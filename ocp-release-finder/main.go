package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// Filter represents a function which accept a string as input and returns a filtered string
type Filter func(string) string

// PullRequestData holds all the data needed for a PR
type PullRequestData struct {
	org   string
	repo  string
	id    int
	mDate time.Time
}

var prd PullRequestData
var releaseStatusCrawler *colly.Collector
var releasePageCrawler *colly.Collector

const OPENSHIFT_RELEASE_DOMAIN string = "openshift-release.svc.ci.openshift.org"
const PR_DATE_FORMAT string = "2006-01-02-150405"

// Gets a PR URL and returns the organization, repo and ID of that PR.
// If the URL is not a valid PR url then returns an error
func parsePrURL(prURL string) (string, string, int, error) {
	githubReg := regexp.MustCompile(`(https*:\/\/)?github\.com\/(.+\/){2}pull\/\d+`)
	if !githubReg.MatchString(prURL) {
		return "", "", 0, fmt.Errorf("ERROR: prURL is not a github PR URL, got url %v", prURL)
	}
	u, err := url.Parse(prURL)
	if err != nil {
		return "", "", 0, err
	}
	pArr := strings.Split(u.Path, "/")
	if len(pArr) != 5 {
		return "", "", 0, fmt.Errorf("Expected PR URL with the form of github.com/ORG/REPO/pull/ID, but instead got %v", prURL)
	}
	id, err := strconv.Atoi(pArr[4])
	if err != nil {
		return "", "", 0, fmt.Errorf("ERROR: Expected PR URL with the form of github.com/ORG/REPO/pull/ID, but instead got %v", prURL)
	}
	return pArr[1], pArr[2], id, nil
}

// Creates a Github client using AccessToken if it exists or an un authenticated client
// if no AccessToken is available and retrieves the PR details from github
func getGithubPrData(org string, repo string, id int) (*github.PullRequest, error) {
	var client *github.Client
	ctx := context.Background()
	accessToken := os.Getenv("AccessToken")
	if accessToken == "" {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	pr, _, err := client.PullRequests.Get(ctx, org, repo, id)
	return pr, err
}

func createPullRequestData(prURL string) (PullRequestData, error) {
	org, repo, id, err := parsePrURL(prURL)
	if err != nil {
		return PullRequestData{}, err
	}
	pr, err := getGithubPrData(org, repo, id)
	if err != nil {
		return PullRequestData{}, err
	}
	return PullRequestData{
		org:   org,
		repo:  repo,
		id:    id,
		mDate: pr.GetMergedAt(),
	}, nil
}

func filterReleasesLinks(e *colly.HTMLElement) {
	reRelease := regexp.MustCompile(`\d\.\d\.\d-\d\.(nightly|ci)-\d{4}-\d{2}-\d{2}-\d{6}`)
	if !reRelease.MatchString(e.Text) {
		return
	}
	d := strings.SplitN(e.Text, "-", 3)[2]
	d = strings.Split(d, " ")[0]
	tr, err := time.Parse(PR_DATE_FORMAT, d)
	if err != nil {
		fmt.Printf("Error time.Parse: ", err)
	}
	// Check if PR merge time t is after the release creation time tr
	if prd.mDate.After(tr) {
		return
	}
	link := e.Attr("href")
	releasePageCrawler.Visit(e.Request.AbsoluteURL(link))
}

func findPR(e *colly.HTMLElement) {
	link := e.Attr("href")
	_, _, id, err := parsePrURL(link)
	if err != nil {
		return
	}
	if prd.id == id {
		fmt.Println("found release, link ", e.Request.URL.String())
	}
}

func setupReleaseStatusCrawler(debug bool) {
	releaseStatusCrawler = colly.NewCollector(
		colly.AllowedDomains(OPENSHIFT_RELEASE_DOMAIN),
		colly.MaxDepth(1),
	)
	if debug == true {
		releaseStatusCrawler.OnRequest(func(r *colly.Request) {
			fmt.Println("Debug: release status crawler is visiting: ", r.URL.String())
		})
	}
	releaseStatusCrawler.OnHTML("a[href]", filterReleasesLinks)
}

func setupReleasePageCrawler(debug bool) {
	releasePageCrawler = colly.NewCollector(
		colly.AllowedDomains(OPENSHIFT_RELEASE_DOMAIN),
		colly.MaxDepth(1),
		colly.Async(true),
	)
	//Set max Parallelism and introduce a Random Delay
	releasePageCrawler.Limit(&colly.LimitRule{
		Parallelism: 6,
		RandomDelay: 1 * time.Second,
	})
	if debug == true {
		releasePageCrawler.OnRequest(func(r *colly.Request) {
			fmt.Println("Debug: release page crawler is visiting: ", r.URL.String())
		})
	}
	releasePageCrawler.OnHTML("a[href]", findPR)
}

func main() {
	var err error
	debug := flag.Bool("debug", false, "run with debug")
	flag.Parse()
	url := flag.Arg(0)
	prd, err = createPullRequestData(url)
	if err != nil {
		panic(err)
	}
	setupReleasePageCrawler(*debug)
	setupReleaseStatusCrawler(*debug)
	releaseStatusCrawler.Visit("https://" + OPENSHIFT_RELEASE_DOMAIN)
	releasePageCrawler.Wait()
	fmt.Println("Done!")
}
