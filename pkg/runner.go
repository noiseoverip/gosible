package pkg

import (
	"ansiblego/pkg/ansible"
	"ansiblego/pkg/logging"
	"ansiblego/pkg/modules"
	"ansiblego/pkg/templating"
	"ansiblego/pkg/transport"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
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
		Strategy: NewParalelExecutor(),
	}
}

func (r *Runner) Run() error {
	inventoryFile, err := os.Open(r.Context.InventoryFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %v", err)
	}
	inventory := &ansible.Inventory{}
	inventory.Dir = path.Dir(r.Context.InventoryFilePath)
	err = ansible.ReadInventory(inventoryFile, inventory)
	if err != nil {
		return fmt.Errorf("failed to load inventory from path %s: %v", r.Context.InventoryFilePath, err)
	}

	logging.Debug("\n# INVENTORY:\n")
	for _, g := range inventory.Groups {
		logging.Debug("\tGroup: %s\n", g.Name)
		for _, h := range g.Hosts {
			logging.Debug("\t\tHost: %s %s\n", h.Name, h.IpAddr)
		}
	}

	groupVars, err := ansible.LoadGroupVars(path.Dir(r.Context.InventoryFilePath))
	if err != nil {
		return fmt.Errorf("failed to load host group variables")
	}

	playbookFile, err := os.Open(r.Context.PlaybookFilePath)
	if err != nil {
		return fmt.Errorf("failed to read playbook from path %s: %v", r.Context.PlaybookFilePath, err)
	}

	playbook := &ansible.Playbook{}
	playbook.Dir = path.Dir(r.Context.PlaybookFilePath)
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
	context := modules.Context{
		PlaybookDir:  playbook.Dir,
		InventoryDir: inventory.Dir,
	}

	for _, play := range playbook.Plays {
		logging.Display("PLAY [%s]", play.HostSelector)
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

		for _, task := range play.Tasks {
			logging.Display("TASK [%s]", task.Name)
			for _, host := range hosts {
				// Handle conditional task execution 'when'
				if task.When != "" {
					result, err := templating.Assert(task.When, host.Vars)
					if err != nil {
						return err
					}
					if !result {
						logging.Info("skipped: [%s]", host.Name)
						continue
					}
				}
				if host.Transport == nil {
					host.Transport = transport.CreateSSHTransport(host.Params)
				}
				r := task.Module.Run(context, host.Transport, host.Vars)
				logging.Debug("Module exec: %s", r)

				if r.Result {
					logging.Info("ok: [%s]", host.Name)
				} else {
					logging.Info("failed: [%s]", host.Name)
					return fmt.Errorf("module execution failed: %s", r.String())
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


type ParalelExecutor struct {

}

func NewParalelExecutor() *ParalelExecutor {
	return &ParalelExecutor{}
}

func(p *ParalelExecutor) Execute(playbook *ansible.Playbook, inventory *ansible.Inventory, vars ansible.GroupVariables) error {
	context := modules.Context{
		PlaybookDir:  playbook.Dir,
		InventoryDir: inventory.Dir,
	}

	for _, play := range playbook.Plays {
		logging.Display("PLAY [%s]", play.HostSelector)
		hosts, err := inventory.GetHosts(play.HostSelector)
		if err != nil {
			return err
		}

		// Build initial host variables by looping though groups it belongs to add variables of that group
		// Precedence:
		// -- group vars (ordered alphabetically )
		// -- host params (from inventory)
		for _, host := range hosts {
			// If there we no variable, it means this is the first play for this host so we load for it variables
			// from group_vars, inventory
			if len(host.Vars) < 1 {
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
		}

		for _, task := range play.Tasks {
			if task.Name == "" {
				task.Name = task.ModuleName
			}
			logging.Display("TASK [%s]", task.Name)

			errors := make(chan error, len(hosts))
			var wg sync.WaitGroup
			for _, host := range hosts {
				// Handle conditional task execution 'when'
				if task.When != "" {
					result, err := templating.Assert(task.When, host.Vars)
					if err != nil {
						return err
					}
					if !result {
						logging.Info("skipped: [%s]", host.Name)
						continue
					}
				}
				if host.Transport == nil {
					host.Transport = transport.CreateSSHTransport(host.Params)
				}
				wg.Add(1)
				go func(host *ansible.Host) {
					defer wg.Done()
					r := task.Module.Run(context, host.Transport, host.Vars)
					logging.Debug("Module exec: %s", r)

					if r.Result {
						logging.Info("ok: [%s]", host.Name)
					} else {
						logging.Info("failed: [%s]", host.Name)
						errors <- fmt.Errorf("module execution failed: %s", r.String())
					}
					// register module output as variable
					if task.Register != "" {
						// TODO: module execusion result should probably be wrapped to add extra information such as
						// execution time, module name...
						host.Vars[task.Register] = r
					}
				}(host)
			}

			wgHostsDone := make(chan struct{})
			go func() {
				defer close(wgHostsDone)
				wg.Wait()
			}()

			select {
			case <-wgHostsDone:
				// do nothing
				logging.Debug("All host finished")
			case <-time.After(1 * time.Minute):
				return fmt.Errorf("timed out waiting for ansible-inventory runs")
			}

			close(errors)
			for err := range errors {
				if err != nil {
					logging.Error("%s", err)
				}
			}
		}
	}
	return nil
}
