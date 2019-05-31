package transport

import (
	"fmt"
	"strings"
)

// SSHTransport abstracts SSH communication
type SSHTransport struct {
	HostAddress string
	Login string
}

func CreateSSHTransport(params map[string]string) Transport {
	return &SSHTransport{HostAddress: params["ansible_host"], Login: params["ansible_user"]}
}

func (ssh *SSHTransport) Exec(command string, args... string) (resultCode int, stdout string, stderr string) {
	fmt.Printf(">>> %s [%s %s]\n", ssh.HostAddress, command, strings.Join(args, " "))
	return 0, "default", "default"
}
