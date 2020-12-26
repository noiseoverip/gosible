package pkg

import (
	"ansiblego/pkg/transport"
	"strings"
)

type Host struct {
	Name      string
	IpAddr    string
	Transport transport.Transport
	Params    map[string]string // These are host params set in inventory as part of host declaration
	Groups    []string
	// These are variables set through group_vars, set_facts, cli...
	Vars HostVariables
}

// 	UnmarshalHost parses inventory line representing host and sets appropriate host fields
// 	Format: [alias key=value key=value]
func UnmarshalHost(input string, h *Host) (err error) {
	els := strings.Split(input, " ")

	// Set defaults
	h.Params = make(map[string]string, 2)
	h.Params["ansible_user"] = "root"
	h.Params["ansible_host"] = "127.0.0.1"
	h.Params["ansible_port"] = "22"

	// For now assume that first element is always an alias, so we force to use aliases !
	h.Name = els[0]
	for _, el := range els[1:] {
		// Put all key-value pairs into a map so we can have any attributes attached and pass them to transport...
		keyVal := strings.Split(el, "=")

		h.Params[keyVal[0]] = keyVal[1]

		// Parse default keys-value pairs
		switch keyVal[0] {
		case "ansible_host":
			h.IpAddr = keyVal[1]
		}
	}
	h.Vars = make(map[string]interface{})
	return nil
}

func (h *Host) String() string {
	return h.IpAddr
}
