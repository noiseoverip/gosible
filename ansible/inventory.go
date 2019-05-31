package ansible

import (
	"bufio"
	"io"
	"strings"
)

type Inventory struct {
	Groups []*HostGroup
}

func (i *Inventory) Group(name string) (group *HostGroup, ok bool ){
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
	inventory.Groups = append(inventory.Groups, currentGroup)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "[") {
			groupName := strings.TrimSpace(line[1:len(line)-1])
			if gtemp, ok := inventory.Group(groupName); ok {
				// Found existing Group
				currentGroup = gtemp
			} else {
				// Found new Group
				currentGroup = &HostGroup{Name: groupName, Hosts: []*Host{}}
				inventory.Groups = append(inventory.Groups, currentGroup)
			}
		} else if isValidHostStart(line) {
			// Found host entry
			// TODO: make sure host object is unique
			host := new(Host)
			err := UnmarshalHost(line, host)
			if err != nil {
				panic(err)
			}
			currentGroup.Hosts = append(currentGroup.Hosts, host)
		}
	}
	return nil
}

func isValidHostStart(hostLine string) bool {
	return len(hostLine) > 0 && !strings.HasPrefix(strings.TrimLeft(hostLine, " "), "#")
}