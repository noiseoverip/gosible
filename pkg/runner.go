package pkg

import (
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
	Execute(playbook *Playbook, inventory *Inventory, vars GroupVariables) error
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

type HostRunSummary struct {
	Executed int
	Failed   int // failure happened on host
	Errors   int // failure happened while attempting to run taks
}

type HostRunSummaries map[string]*HostRunSummary

func (h HostRunSummaries) hasFails() bool {
	for _, sumarry := range h {
		if sumarry.Failed > 0 {
			return true
		}
	}
	return false
}

type TaskResultWrap struct {
	*modules.ModuleExecResult
	Host *Host
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
	inventory := &Inventory{}
	inventory.Dir = path.Dir(r.Context.InventoryFilePath)
	err = ReadInventory(inventoryFile, inventory)
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

	groupVars, err := LoadGroupVars(path.Dir(r.Context.InventoryFilePath))
	if err != nil {
		return fmt.Errorf("failed to load host group variables")
	}

	playbookFile, err := os.Open(r.Context.PlaybookFilePath)
	if err != nil {
		return fmt.Errorf("failed to read playbook from path %s: %v", r.Context.PlaybookFilePath, err)
	}

	playbook := &Playbook{}
	playbook.Dir = path.Dir(r.Context.PlaybookFilePath)
	err = ReadPlaybook(playbookFile, playbook)
	if err != nil {
		return fmt.Errorf("failed to load playbook from path %s: %v", r.Context.PlaybookFilePath, err)
	}

	//TODO: this should receive pointers to stdout and stdr se we could control them from higher level
	return r.Strategy.Execute(playbook, inventory, groupVars)
}

func NewParalelExecutor() *ParalelExecutor {
	return &ParalelExecutor{}
}

type ParalelExecutor struct{}

func (p *ParalelExecutor) Execute(playbook *Playbook, inventory *Inventory, vars GroupVariables) error {
	context := modules.Context{
		PlaybookDir:  playbook.Dir,
		InventoryDir: inventory.Dir,
	}

	hostRunSummary := &HostRunSummaries{}
	for _, play := range playbook.Plays {
		logging.Header("PLAY [%s]", play.HostSelector)
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
			(*hostRunSummary)[host.Name] = &HostRunSummary{}
		}
		for _, task := range play.Tasks {
			if task.Name == "" {
				task.Name = task.ModuleName
			}
			logging.Header("TASK [%s]", task.Name)

			results := make(chan *TaskResultWrap, len(hosts))
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
				(*hostRunSummary)[host.Name].Executed += 1
				go func(host *Host) {
					defer wg.Done()
					// TODO: here we are passing a copy of host variables yet set_fact seems to work so i am confused.
					//   ideally modules should not have ability to modify host variables directly.
					moduleResult := task.Module.Run(context, &modules.Host{Transport: host.Transport, Vars: host.Vars})

					logging.Debug("Module exec: %s", moduleResult)
					// register module output as variable
					if task.Register != "" {
						host.Vars[task.Register] = moduleResult
					}
					results <- &TaskResultWrap{ModuleExecResult: moduleResult, Host: host}
				}(host)
			}

			wgHostsDone := make(chan struct{})
			go func() {
				defer close(wgHostsDone)
				wg.Wait()
			}()

		loop:
			for {
				select {
				case result := <-results:
					if result.Result {
						logging.Info("ok: [%s]", result.Host.Name)
					} else {
						(*hostRunSummary)[result.Host.Name].Failed += 1
						logging.Info("failed: [%s] %s", result.Host.Name, result.String())
					}
				case <-wgHostsDone:
					// do nothing
					logging.Debug("All host finished")
					break loop

				case <-time.After(1 * time.Minute): //TODO: this should be configurable option
					logging.Error("timed out waiting for ansible-inventory runs")
					break loop
				}
			}
			if hostRunSummary.hasFails() {
				break
			}
		}
	}

	logging.Header("PLAY RECAP")
	for hostName, summary := range *hostRunSummary {
		logging.Info(fmt.Sprintf("%-12s : ok=%d failed=%d", hostName, summary.Executed-summary.Failed, summary.Failed))
	}
	if hostRunSummary.hasFails() {
		return fmt.Errorf("there were failures")
	}

	return nil
}
