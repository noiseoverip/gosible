package pkg

import (
	"ansiblego/pkg/ansible"
	"ansiblego/pkg/logging"
	"ansiblego/pkg/templating"
	"ansiblego/pkg/transport"
	"fmt"
	"os"
	"path"
)

type Executer interface {
	Execute(playbook *ansible.Playbook, inventory *ansible.Inventory, vars ansible.GroupVariables) error
}

// Context holds shared objects
type Context struct {
	InventoryFilePath string
	PlaybookFilePath  string
}

// Runner is responsible for loading all required files and executing a playbook
type Runner struct {
	Context  *Context
	Strategy Executer
}

func NewRunner(inventory string, playbook string) *Runner {
	return &Runner{
		Context:  &Context{InventoryFilePath: inventory, PlaybookFilePath: playbook},
		Strategy: NewSequentialExecuter(),
	}
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

	logging.Info("\n# INVENTORY:\n")
	for _, g := range inventory.Groups {
		logging.Info("\tGroup: %s\n", g.Name)
		for _, h := range g.Hosts {
			logging.Info("\t\tHost: %s %s\n", h.Name, h.IpAddr)
		}
	}
	logging.Info("\n")

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
	return r.Strategy.Execute(playbook, inventory, groupVars)
}

func NewSequentialExecuter() Executer {
	return &SequentialExecuter{}
}

type SequentialExecuter struct {
}

func (s SequentialExecuter) Execute(playbook *ansible.Playbook, inventory *ansible.Inventory, vars ansible.GroupVariables) error {
	for _, play := range playbook.Plays {
		logging.L.InfoLogger.Printf("\n### Running play [%s] ###\n\n", play.HostSelector)

		hosts, err := inventory.GetHosts(play.HostSelector)
		if err != nil {
			return err
		}

		// TODO: the way it is done now, host variable will not persist across plays
		// Build initial host variables by looping though groups it belongs to add variables of that group
		// Precedence:
		// -- group vars (ordered alphabetically )
		// -- host params (from inventory)
		// -- TODO: cli
		for _, host := range hosts {
			for _, group := range host.Groups {
				if vars, found := vars[group]; found {
					host.Vars.Add(vars)
				}
			}
			// Override group variables with host params from inventory
			for k, v := range host.Params {
				host.Vars[k] = v
			}
		}

		// Execute tasks on each host sequentially (for now)
		for _, host := range hosts {
			// Since role is just a list of tasks, we could simply expand it to tasks and attach role information to
			// task (for tracking and logging mostly). Then task execution would not change.
			for _, task := range play.Tasks {
				// Handle conditional task execution 'when'
				if task.When != "" {
					result, err := templating.Assert(task.When, host.Vars)
					if err != nil {
						return err
					}
					if !result {
						logging.Debug("Skipping task [%s] on host %s\n", task.Name, host)
						continue
					}
				}

				logging.Debug("\n### Running task [%s] on host %s\n", task.Name, host)
				if host.Transport == nil {
					host.Transport = transport.CreateSSHTransport(host.Params)
				}
				r := task.Module.Run(host.Transport, host.Vars)
				logging.Debug("Module exec: %s", r)
				if !r.Result {
					return fmt.Errorf("module execution failed")
				}
				// register module output as variable
				if task.Register != "" {
					// TODO: module execusion result should probably be wrapped to add extra information such as
					// execution time, module name...
					host.Vars[task.Register] = r
				}
			}
		}

	}
	return nil
}
