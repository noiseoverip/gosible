package main

import (
	"ansiblego/ansible"
	"flag"
	"fmt"
	"os"
)

var inventoryPath = flag.String("i", "", "Path to inventory")
// Playbook path is defined as cli argument

func main() {
	flag.Parse()
	inventoryFile, err := os.Open(*inventoryPath)
	if err != nil {
		fmt.Printf("Failed to open file %v", err)
		os.Exit(1)
	}
	inventory := new(ansible.Inventory)
	err = ansible.ReadInventory(inventoryFile, inventory)
	if err != nil {
		fmt.Printf("Failed to load inventory from path %s: %v", inventoryPath, err)
		os.Exit(1)
	}
	if len(flag.Args()) < 1 {
		fmt.Printf("Too few arguments provided, please provide path to playbook")
		os.Exit(1)
	}
	playbookPath := flag.Arg(0)
	playbookFile, err := os.Open(playbookPath)
	playbook := new(ansible.Playbook)
	err = ansible.ReadPlaybook(playbookFile, playbook)
	if err != nil {
		fmt.Printf("Failed to load playbook from path %s: %v", playbookPath, err)
		os.Exit(1)
	}

	playbook.Run(inventory)
}