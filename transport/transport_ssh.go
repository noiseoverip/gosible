package transport

import (
	"ansiblego/logging"
	"bytes"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// SSHTransport abstracts SSH communication
type SSHTransport struct {
	HostAddress string
	Login       string
	SSHClient   *ssh.Client
	SCPSession  *scp.Client
}

func CreateSSHTransport(params map[string]string) Transport {
	return &SSHTransport{HostAddress: params["ansible_host"], Login: params["ansible_user"]}
}

// Exec executes command
func (t *SSHTransport) Exec(command string, args ...string) (resultCode int, stdout string, stderr string, err error) {
	logging.Info(">>> host:%s [%s %s]\n", t.HostAddress, command, strings.Join(args, " "))
	// TODO: re-use ssh sessions. Host should keep "connection" object, each command should run in its own session
	if t.SSHClient == nil {
		client, err := t.createSSHClient(t.Login, t.HostAddress, "/Users/salisauskas/.ssh/id_rsa")
		if err != nil {
			return -1, "", "", fmt.Errorf("failed to create client: %s", err)
		}
		t.SSHClient = client
	}

	session, err := t.SSHClient.NewSession()
	if err != nil {
		return -1, "", "", fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	var berr bytes.Buffer
	session.Stderr = &berr
	var bout bytes.Buffer
	session.Stdout = &bout
	cmd := fmt.Sprintf("/bin/sh -c \"%s %s\"", command, strings.Join(args, " "))

	err = session.Run(cmd)
	if err != nil {
		logging.Info("ERROR %v", err)
		return -1, string(bout.Bytes()), string(berr.Bytes()), err
	}

	return 0, string(bout.Bytes()), "", nil
}

func (t *SSHTransport) openFileTransferSession(loginName string, hostIpAddress string, privateKeyPath string) (session *scp.Client, err error) {
	clientConfig, _ := auth.PrivateKey(loginName, privateKeyPath, ssh.InsecureIgnoreHostKey())
	client := scp.NewClient(hostIpAddress, &clientConfig)
	err = client.Connect()
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (t *SSHTransport) SendFileToRemote(srcFilePath string, destFilePath string, mode string) (err error) {
	// TODO: re-use scp sessions.
	t.SCPSession, err = t.openFileTransferSession(t.Login, fmt.Sprintf("%s:22", t.HostAddress), "/Users/salisauskas/.ssh/id_rsa")
	if err != nil {
		return fmt.Errorf("failed to create session %v:", err)
	}
	defer t.SCPSession.Close()
	f, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("failed opening file %v:", err)
	}
	defer f.Close()

	err = t.SCPSession.CopyFile(f, destFilePath, mode)
	if err != nil {
		return fmt.Errorf("error while copying file: %v", err)
	}
	return nil
}

func (t *SSHTransport) createSSHClient(loginName string, hostIpAddress string, privateKeyPath string) (*ssh.Client, error) {
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

	log.Printf("Opening SSH connection to %s@%s", loginName, nodeAddr)
	return ssh.Dial("tcp", nodeAddr, sshConfig)

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
