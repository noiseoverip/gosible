package transport

import (
	"bbgithub.dev.bloomberg.com/babka/cli/pkg/util"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

// SSHTransport abstracts SSH communication
type SSHTransport struct {
	HostAddress string
	Login string
}

func CreateSSHTransport(params map[string]string) Transport {
	return &SSHTransport{HostAddress: params["ansible_host"], Login: params["ansible_user"]}
}

// Exec executes command
func (t *SSHTransport) Exec(command string, args... string) (resultCode int, stdout string, stderr string) {
	fmt.Printf(">>> host:%s [%s %s]\n", t.HostAddress, command, strings.Join(args, " "))
	// TODO: add ability to define key at the host level or use default one
	session, err := t.openSession(t.Login, t.HostAddress, "/Users/salisauskas/.ssh/id_rsa")
	defer session.Close() // TODO: reuse session
	if err != nil {
		return -1, "", fmt.Sprintf("Failed to create ssh session: %v", err)
	}

	var berr bytes.Buffer
	session.Stderr = &berr
	var bout bytes.Buffer
	session.Stdout = &bout
	cmd := fmt.Sprintf("/bin/sh -c \"%s %s\"", command, strings.Join(args, " "))

	err = session.Run(cmd)
	if err != nil {
		return -1, string(bout.Bytes()), string(berr.Bytes())
	}

	return 0, string(bout.Bytes()), ""
}

func (t *SSHTransport) openSession(loginName string, hostIpAddress string, privateKeyPath string) (*ssh.Session, error) {
	sshAuth, err := SSHAuthWithKey(privateKeyPath)
	if err != nil {
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: loginName,
		Auth: []ssh.AuthMethod{
			sshAuth, //TODO: we should load this at the startup
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Ignore host fingerprints
			return nil
		},
		Timeout: time.Second * 5,
	}
	nodeAddr := fmt.Sprintf("%s:22", hostIpAddress)
	util.Success("Opening SSH connection to %s@%s", loginName, nodeAddr)
	connection, err := ssh.Dial("tcp", nodeAddr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	//modes := ssh.TerminalModes{
	//	ssh.ECHO:          0,     // disable echoing
	//	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	//	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	//}

	//if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
	//	session.Close()
	//	return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	//}

	//stdin, err := session.StdinPipe()
	//if err != nil {
	//	return fmt.Errorf("unable to setup stdin for session: %v", err)
	//}
	//go io.Copy(stdin, os.Stdin)
	//
	//stdout, err := session.StdoutPipe()
	//if err != nil {
	//	return fmt.Errorf("unable to setup stdout for session: %v", err)
	//}
	//go io.Copy(os.Stdout, stdout)
	//
	//stderr, err := session.StderrPipe()
	//if err != nil {
	//	return fmt.Errorf("unable to setup stderr for session: %v", err)
	//}
	//go io.Copy(os.Stderr, stderr)

	//if err := session.Start("/bin/bash"); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := session.Wait(); err != nil {
	//	log.Fatal(err)
	//}
	return session, nil
}

// SSHAuthWithKey creates ssh.AuthMethod based on certificate
func SSHAuthWithKey(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
