package runner

import (
	"ansiblego/ansible"
	"ansiblego/utils"
	"fmt"
	"os"
	"path"
)

// Context holds shared objects
type Context struct {
	Logger *utils.GosibleLogger
	InventoryFilePath string
	PlaybookFilePath string
}

// Runner is responsible for loading all required files and executing a playbook
type Runner struct {
	Context *Context
}

func (r *Runner) Run() error {
	inventoryFile, err := os.Open(r.Context.InventoryFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %v", err)
	}
	inventory := &ansible.Inventory{}
	err = ansible.ReadInventory(inventoryFile, inventory)
	if err != nil {
		return fmt.Errorf("failed to load inventory from path %s: %v", r.Context.InventoryFilePath, err)
	}

	fmt.Printf("\n# INVENTORY:\n")
	for _, g := range inventory.Groups {
		fmt.Printf("\tGroup: %s\n", g.Name)
		for _, h := range g.Hosts {
			fmt.Printf("\t\tHost: %s %s\n", h.Name, h.IpAddr)
		}
	}
	fmt.Printf("\n")

	groupVars, err := ansible.LoadGroupVars(path.Dir(r.Context.InventoryFilePath))
	if err != nil {
		return fmt.Errorf("failed to load host group variables")
	}

	playbookFile, err := os.Open(r.Context.PlaybookFilePath)
	if err != nil {
		return fmt.Errorf("failed to read playbook from path %s: %v", r.Context.PlaybookFilePath, err)
	}

	playbook := &ansible.Playbook{}
	err = ansible.ReadPlaybook(playbookFile, playbook)
	if err != nil {
		return fmt.Errorf("failed to load playbook from path %s: %v", r.Context.PlaybookFilePath, err)
	}

	//TODO: this should receive pointers to stdout and stdr se we could control them from higher level
	return playbook.Run(inventory, groupVars)
}