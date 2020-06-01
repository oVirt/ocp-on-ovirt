# Openshift Release Finder

Under the openshift organization we have various github repos and each merged PR to those repos will trigger an openshift release build. Finding which release was triggered by a PR can be a pain, we will need to manually look at when the PR was merged and then go to [Openshift release page](openshift-release.svc.ci.openshift.org) and find a release job which was triggered after the time the PR was merged, check its page and see if it contains a PR if not move on to the next release page. My goal in this small tool is to automate this task.

## Usage

1. Clone the repo or download the executable.
2. Give it execute permissions: chmod +x ./ocp-release-finder
3. Trigger the tool:
   1. Simple usage:

   ```base
   ./ocp-release-finder https://github.com/openshift/console/pull/5628
   ```

   2. With github AccessToken:
   ```base
   export AccessToken="PASTE YOUR ACCESS TOKEN"
   ./ocp-release-finder https://github.com/openshift/console/pull/5628
   ```
   ** This is needed in case you are using this tool many times in a short time frame. We use github API to retrieve information about the PR, when we use the API to frequently as an unauthenticated user we will get blocked, see information about how to set up an access token in [link](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)

   3. With Debug:

   ```base
   ./ocp-release-finder https://github.com/openshift/console/pull/5628 --debug
   ```

## Important note

This tool will work only one PRs that are newlly merged PRs and have releases on Openshift release page that don't contain them. The reason is that this tool is a web crawler that automate the manualy process so if the PR is old then all the release will contain it, and it will not find it in any release page. In the future I will add logic to inform the user in that case.

