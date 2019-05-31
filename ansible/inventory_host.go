package ansible

import (
	"ansiblego/transport"
	"strings"
)

type Host struct {
	Name string
	IpAddr string
	Login string
	Transport transport.Transport
	Params map[string]string
}

// 	UnmarshalHost parses inventory line representing host and sets appropriate host fields
// 	Format: [alias key=value key=value]
//		Required keys:
//			ansible_host (default: 127.0.0.1)
//			ansible_user (default: root)
func UnmarshalHost(input string, h *Host) error {
	els := strings.Split(input, " ")

	// For now assume that first element is always an alias, so we force to use aliases !
	h.Name = els[0]
	for _, el := range els[1:] {
		// Put all key-value pairs into a map so we can have any attributes attached and pass them to transport...
		keyVal:= strings.Split(el, "=")
		h.Params = make(map[string]string, 2)
		// Set defaults
		h.Params["ansible_user"] = "root"
		h.Params["ansible_host"] = "127.0.0.1"

		h.Params[keyVal[0]] = keyVal[1]
		// Parse default keys-value pairs
		switch keyVal[0] {
		case "ansible_host":
			h.IpAddr = keyVal[1]
		case "ansible_user":
			h.Login = keyVal[1]
		}
	}
	return nil
}

func (h *Host) String() string {
	return h.IpAddr
}
