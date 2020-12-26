package pkg

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Inventory struct {
	// Inventory from host group perspective. Here, group containers pointers to Host objects
	Groups []*HostGroup
	// Inventory from host perspective
	Hosts []*Host

	Dir string
}

// TODO: Here we should interpret host selector and build a list of hosts we should run tasks play on but for now we
//   limit our selves to only supporting running on groups
func (i *Inventory) GetHosts(selector string) ([]*Host, error) {
	if group, found := i.groupByName(selector); found {
		return group.Hosts, nil
	}
	return nil, fmt.Errorf("currently only support groups as host selector and groupByName %s was not found", selector)
}

func (i *Inventory) groupByName(name string) (group *HostGroup, ok bool) {
	for _, g := range i.Groups {
		if g.Name == name {
			return g, true
		}
	}
	return nil, false
}

func ReadInventory(in io.Reader, inventory *Inventory) error {
	sc := bufio.NewScanner(in)
	inventory.Groups = []*HostGroup{}
	currentGroup := &HostGroup{Name: "all", Hosts: []*Host{}}
	groupAll := currentGroup
	inventory.Groups = append(inventory.Groups, currentGroup)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "[") {
			groupName := strings.TrimSpace(line[1 : len(line)-1])
			if gtemp, ok := inventory.groupByName(groupName); ok {
				// Found existing groupByName
				currentGroup = gtemp
			} else {
				// Found new groupByName
				currentGroup = &HostGroup{Name: groupName, Hosts: []*Host{}}
				inventory.Groups = append(inventory.Groups, currentGroup)
			}
		} else if isValidHostStart(line) {
			// Found host entry
			host := new(Host)
			err := UnmarshalHost(line, host)
			if err != nil {
				panic(err)
			}
			if existingHost, found := findHost(host.Name, inventory.Hosts); found {
				// Existing host means we already took care of group all
				host = existingHost
			} else if currentGroup.Name != "all" {
				// In case this is first time we are seeing the host and we are already in context of some other group
				host.Groups = append(host.Groups, "all")
				groupAll.Hosts = append(groupAll.Hosts, host)
			}
			host.Groups = append(host.Groups, currentGroup.Name)
			currentGroup.Hosts = append(currentGroup.Hosts, host)
			inventory.Hosts = append(inventory.Hosts, host)
		}
	}
	// Sort host group alphabetically
	for _, host := range inventory.Hosts {
		sort.Strings(host.Groups)
	}
	return nil
}

// Find host by name in a given slice
func findHost(name string, hosts []*Host) (host *Host, found bool) {
	for _, h := range hosts {
		if h.Name == name {
			return h, true
		}
	}
	return nil, false
}

func isValidHostStart(hostLine string) bool {
	return len(hostLine) > 0 && !strings.HasPrefix(strings.TrimLeft(hostLine, " "), "#")
}
