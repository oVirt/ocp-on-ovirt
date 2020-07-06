package ovirtci

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	ssh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type ChannelLease struct {
	FileName string
	Content  string
	address  string
}

type Proxyvm struct {
	Address    string
	Sshconfig  *ssh.ClientConfig
	Sshconfig2 *ssh.ClientConfig
	session    *ssh.Session
	stdoutBuf  bytes.Buffer
	SshUser    string
	SshKey     string
	VmSshKey   string
}

func (p *Proxyvm) SetSshconfig(config *ssh.ClientConfig) {
	p.Sshconfig = config
}

func (p *Proxyvm) SetSshconfig2(config *ssh.ClientConfig) {
	p.Sshconfig2 = config
}

func (p *Proxyvm) ConnectSsh() {

	bClient, err := ssh.Dial("tcp", p.Address, p.Sshconfig)
	if err != nil {
		log.Fatal(err)
	}
	session, err := bClient.NewSession()
	session.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	session.Stderr = os.Stderr

	p.stdoutBuf = bytes.Buffer{}
	session.Stdout = &p.stdoutBuf
	p.session = session
	//p.session
}

//run command on ocp VM
func (p *Proxyvm) RunSshVM(cmd string, addr string, c chan ChannelLease) {

	//connect to the proxy VM
	connected := true

	bClient, err := ssh.Dial("tcp", p.Address, p.Sshconfig)
	if err != nil {
		log.Error(err)
		connected = false
	}

	//connect to the internal VM
	conn, err := bClient.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
		connected = false
	}

	ch := ChannelLease{}
	if connected {
		ncc, chans, reqs, err := ssh.NewClientConn(conn, addr, p.Sshconfig2)
		if err != nil {
			log.Error(err)
		}
		sClient := ssh.NewClient(ncc, chans, reqs)
		buff := bytes.Buffer{}

		session, err := sClient.NewSession()
		session.Stdin = os.Stdin
		session.Stdout = &buff
		session.Stderr = os.Stderr

		//p.stdoutBuf = bytes.Buffer{}

		if err := session.Run(cmd); err != nil {
			switch v := err.(type) {
			case *ssh.ExitError:
				os.Exit(v.Waitmsg.ExitStatus())
			default:
				log.Fatalln(err)
			}
		}
		ch = ChannelLease{FileName: cmd, Content: buff.String(), address: addr}
	}

	c <- ch
}

func (p *Proxyvm) openTTY(addr string, ctx context.Context) error {

	//connect to the proxy VM
	connected := true

	bClient, err := ssh.Dial("tcp", p.Address, p.Sshconfig)
	if err != nil {
		log.Error(err)
		connected = false
	}

	//connect to the internal VM
	conn, err := bClient.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
		connected = false
	}

	if connected {
		ncc, chans, reqs, err := ssh.NewClientConn(conn, addr, p.Sshconfig2)
		if err != nil {
			log.Error(err)
		}
		sClient := ssh.NewClient(ncc, chans, reqs)

		session, err := sClient.NewSession()
		//session.Stdin = os.Stdin
		//session.Stdout = &p.stdoutBuf
		//session.Stderr = os.Stderr
		defer session.Close()

		go func() {
			<-ctx.Done()
			conn.Close()
		}()

		fd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(fd)
		if err != nil {
			return fmt.Errorf("terminal make raw: %s", err)
		}
		defer terminal.Restore(fd, state)

		w, h, err := terminal.GetSize(fd)
		if err != nil {
			return fmt.Errorf("terminal get size: %s", err)
		}

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 36000,
			ssh.TTY_OP_OSPEED: 36000,
		}

		term := os.Getenv("TERM")
		if term == "" {
			term = "xterm-256color"
		}
		if err := session.RequestPty(term, h, w, modes); err != nil {
			return fmt.Errorf("session xterm: %s", err)
		}

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin

		if err := session.Shell(); err != nil {
			return fmt.Errorf("session shell: %s", err)
		}

		if err := session.Wait(); err != nil {
			if e, ok := err.(*ssh.ExitError); ok {
				switch e.ExitStatus() {
				case 130:
					return nil
				}
			}
			return fmt.Errorf("ssh: %s", err)
		}

	}
	return nil
}

func (p *Proxyvm) RunSsh(cmd string, c chan ChannelLease) {
	if p.session == nil {
		p.ConnectSsh()
	}
	if err := p.session.Run(cmd); err != nil {
		switch v := err.(type) {
		case *ssh.ExitError:
			os.Exit(v.Waitmsg.ExitStatus())
		default:
			log.Fatalln(err)
		}
	}

	ch := ChannelLease{FileName: cmd, Content: p.stdoutBuf.String()}
	//os.Stdin.Sync()
	p.session.Stdout = os.Stdout

	c <- ch
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

//TODO: RunCommands - run multiply commands in bastion VM

//RunProxyVM - run given command in bastion VM (proxy VM)
func RunProxyVM(cmd string, proxy Proxyvm) ChannelLease {
	//proxy := Proxyvm{Address: vmaddresss}
	sshConfig := &ssh.ClientConfig{
		User:            proxy.SshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(proxy.SshKey),
		},
	}
	proxy.SetSshconfig(sshConfig)
	c := make(chan ChannelLease)
	go proxy.RunSsh(cmd, c)
	res := <-c
	return res
}

//RunVM - run given command in given address VM throught proxy VM
func RunVM(address string, proxy Proxyvm, cmd string) (string, error) {
	sshConfig := &ssh.ClientConfig{
		User:            proxy.SshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(proxy.SshKey),
		},
	}

	sshConfig2 := &ssh.ClientConfig{
		User:            "core",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(proxy.VmSshKey),
		},
	}

	proxy.SetSshconfig(sshConfig)
	proxy.SetSshconfig2(sshConfig2)
	c := make(chan ChannelLease)
	go proxy.RunSshVM(cmd, fmt.Sprintf("%s:22", address), c)
	res := <-c
	if res.Content != "" {
		return res.Content, nil
	} else {
		return res.Content, fmt.Errorf("error connecting")
	}
}

//RunVMMany - connect to multiply VMs and run single command
func RunVMMany(addressList []string, proxy Proxyvm, cmd string) (map[string]string, error) {
	sshConfig := &ssh.ClientConfig{
		User:            proxy.SshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(proxy.SshKey),
		},
	}

	sshConfig2 := &ssh.ClientConfig{
		User:            "core",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(proxy.VmSshKey),
		},
	}

	proxy.SetSshconfig(sshConfig)
	proxy.SetSshconfig2(sshConfig2)
	c := make(chan ChannelLease)
	for _, address := range addressList {
		log.Debugln("adding ", address)
		go proxy.RunSshVM(cmd, fmt.Sprintf("%s:22", address), c)
	}

	retmap := map[string]string{}
	for range addressList {
		res := <-c
		retmap[res.address] = res.Content
	}

	return retmap, nil
}

func OpenTTY(address string, proxy Proxyvm) error {
	sshConfig := &ssh.ClientConfig{
		User:            proxy.SshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(proxy.SshKey),
		},
		Timeout: 2,
	}

	sshConfig2 := &ssh.ClientConfig{
		User:            "core",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile(proxy.VmSshKey),
		},
		Timeout: 2,
	}

	proxy.SetSshconfig(sshConfig)
	proxy.SetSshconfig2(sshConfig2)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := proxy.openTTY(fmt.Sprintf("%s:22", address), ctx); err != nil {
			log.Print(err)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}

	return nil
}
