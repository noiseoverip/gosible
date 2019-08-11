package runner

import (
	"ansiblego/ansible"
	"fmt"
	"os"
	"path"
)

// Runner is responsible for loading all required files and executing a playbook
type Runner struct {
	InventoryFilePath string
	PlaybookFilePath string
}

func (r *Runner) Run() error {
	inventoryFile, err := os.Open(r.InventoryFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %v", err)
	}
	inventory := &ansible.Inventory{}
	err = ansible.ReadInventory(inventoryFile, inventory)
	if err != nil {
		return fmt.Errorf("failed to load inventory from path %s: %v", r.InventoryFilePath, err)
	}

	fmt.Printf("\n# INVENTORY:\n")
	for _, g := range inventory.Groups {
		fmt.Printf("\tGroup: %s\n", g.Name)
		for _, h := range g.Hosts {
			fmt.Printf("\t\tHost: %s %s\n", h.Name, h.IpAddr)
		}
	}
	fmt.Printf("\n")

	groupVars, err := ansible.LoadGroupVars(path.Dir(r.InventoryFilePath))
	if err != nil {
		return fmt.Errorf("failed to load host group variables")
	}

	playbookFile, err := os.Open(r.PlaybookFilePath)
	if err != nil {
		return fmt.Errorf("failed to read playbook from path %s: %v", r.PlaybookFilePath, err)
	}

	playbook := &ansible.Playbook{}
	err = ansible.ReadPlaybook(playbookFile, playbook)
	if err != nil {
		return fmt.Errorf("failed to load playbook from path %s: %v", r.PlaybookFilePath, err)
	}

	//TODO: this should receive pointers to stdout and stdr se we could control them from higher level
	return playbook.Run(inventory, groupVars)
}