package transport

import (
	"ansiblego/pkg/logging"
	"bytes"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

// SSHTransport abstracts SSH communication
type SSHTransport struct {
	HostAddress string
	Login       string
	SSHClient   *ssh.Client
	SCPClient   *scp.Client
	KeyPath 	string
}

func CreateSSHTransport(params map[string]string) Transport {
	return &SSHTransport{
		HostAddress: params["ansible_host"],
		Login: params["ansible_user"],
		KeyPath: params["ansible_ssh_private_key"],
	}
}

// Exec executes command
func (t *SSHTransport) Exec(command string, args ...string) (resultCode int, stdout string, stderr string, err error) {
	logging.Debug(">>> host:%s [%s %s]\n", t.HostAddress, command, strings.Join(args, " "))
	t.SSHClient, err = t.sshClient()
	if err != nil {
		return -1, "", "", fmt.Errorf("failed to create client: %s", err)
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

func (t *SSHTransport) scpClient(loginName string, hostIpAddress string, privateKeyPath string) (session *scp.Client, err error) {
	clientConfig, _ := auth.PrivateKey(loginName, privateKeyPath, ssh.InsecureIgnoreHostKey())
	client := scp.NewClient(hostIpAddress, &clientConfig)
	return &client, nil
}

func (t *SSHTransport) SendFileToRemote(srcFilePath string, destFilePath string, mode string) (err error) {
	// TODO: re-use scp sessions.
	if t.SCPClient == nil {
		t.SCPClient, err = t.scpClient(t.Login, fmt.Sprintf("%s:22", t.HostAddress), t.KeyPath)
		if err != nil {
			return fmt.Errorf("failed to create session :%v ", err)
		}
	}
	t.SSHClient, err = t.sshClient()
	session, err := t.SSHClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()
	t.SCPClient.Session = session
	f, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("failed opening file %v:", err)
	}
	defer f.Close()

	err = t.SCPClient.CopyFile(f, destFilePath, mode)
	if err != nil {
		return fmt.Errorf("error while copying file: %v", err)
	}
	return nil
}

func (t *SSHTransport) sshClient() (*ssh.Client, error) {
	if t.SSHClient != nil {
		return t.SSHClient, nil
	}
	if t.KeyPath == "" {
		defaultKeypath, err := DefaultSSHKeyPath()
		if err != nil {
			return nil, err
		}
		t.KeyPath = defaultKeypath
	}
	logging.Debug("Using ssh key path: %s", t.KeyPath)
	sshAuth, err := SSHAuthWithKey(t.KeyPath)
	if err != nil {
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: t.Login,
		Auth: []ssh.AuthMethod{
			sshAuth, //TODO: we should load this at the startup
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Ignore host fingerprints
			return nil
		},
		Timeout: time.Second * 5,
	}
	nodeAddr := fmt.Sprintf("%s:22", t.HostAddress)

	logging.Debug("Opening SSH connection to %s@%s", t.Login, nodeAddr)
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

func DefaultSSHKeyPath() (keyPath string, err error){
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDirPath, ".ssh", "id_rsa"), nil

}
