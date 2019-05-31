package transport

import (
	"ansiblego/internal/logging"
	"bytes"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// SSHTransport abstracts SSH communication
type SSHTransport struct {
	HostAddress string
	Login       string
	Port        int
	SSHClient   *ssh.Client
	SCPClient   *scp.Client
	KeyPath     string
}

// CreateSSHTransport creates ssh transport based on from provided params. Keeping this as flexible as possible for now:
// simply pass all host variables to here so we can use whatever is needed directly. Only public key authentication is
// supported (password wont' work)
func CreateSSHTransport(params map[string]string) (Transport, error) {
	portString := params["ansible_port"]
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, fmt.Errorf("ansible_port must be integer but got: %s", portString)
	}
	return &SSHTransport{
		HostAddress: params["ansible_host"],
		Login:       params["ansible_user"],
		Port:        port,
		KeyPath:     params["ansible_ssh_private_key"],
	}, nil
}

// Exec executes command on remote host using SSH connection
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
	defer func() {
		_ = session.Close()
	}()

	var berr, bout bytes.Buffer
	session.Stderr = &berr
	session.Stdout = &bout
	cmd := fmt.Sprintf("/bin/sh -c \"%s %s\"", command, strings.Join(args, " "))

	err = session.Run(cmd)
	if err != nil {
		logging.Info("ERROR %v", err)
		return -1, bout.String(), bout.String(), err
	}

	return 0, bout.String(), "", nil
}

// scpClient create SCP client
func (t *SSHTransport) scpClient(loginName string, hostIpAddress string, privateKeyPath string) (session *scp.Client, err error) {
	logging.Debug("initializing scp client")
	clientConfig, _ := auth.PrivateKey(loginName, privateKeyPath, ssh.InsecureIgnoreHostKey())
	client := scp.NewClient(hostIpAddress, &clientConfig)
	return &client, nil
}

// SendFileToRemote sends a given file to remote host.
// TODO: this doesn't seem to be re-using sessions which might be worth improving
func (t *SSHTransport) SendFileToRemote(srcFilePath string, destFilePath string, mode string) (err error) {
	if t.SCPClient == nil {
		t.SCPClient, err = t.scpClient(t.Login, fmt.Sprintf("%s:%d", t.HostAddress, t.Port), t.KeyPath)
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
	logging.Debug("initializing ssh client")
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
			sshAuth, //TODO: should load this at the startup
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Ignore host fingerprints
			return nil
		},
		Timeout: time.Second * 5,
	}
	nodeAddr := fmt.Sprintf("%s:%d", t.HostAddress, t.Port)

	logging.Debug("Opening SSH connection to %s@%s", t.Login, nodeAddr)
	return ssh.Dial("tcp", nodeAddr, sshConfig)
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

func DefaultSSHKeyPath() (keyPath string, err error) {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDirPath, ".ssh", "id_rsa"), nil
}
