package ansible

import (
	"ansiblego/pkg/logging"
	"ansiblego/pkg/templating"
	"ansiblego/pkg/transport"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

type Playbook struct {
	Plays []Play
}

// ReadPlaybook loads playbook yml file into struct
func ReadPlaybook(in io.Reader, playbook *Playbook) error {
	buf := new(bytes.Buffer)
	bytesRead, err := buf.ReadFrom(in)
	if bytesRead < 1 {
		return fmt.Errorf("empty file")
	} else if err != nil  {
		return fmt.Errorf("error reading file: %v", err)
	}
	// Playbook is essentially a slice of unnamed
	var plays []Play
	if err := yaml.Unmarshal(buf.Bytes(), &plays); err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to read yml %v", err))
	}
	playbook.Plays = plays
	return nil
}

// Run playbook:
// Play (- hosts:..) sequentially
// -- each host (sequentially)
// -- -- tasks (sequentially)
func (playbook *Playbook) Run(inventory *Inventory, groupVars GroupVariables) error {
	for _, play := range playbook.Plays {
		log("\n### Running play [%s] ###\n\n", play.HostSelector)

		hosts, err := inventory.GetHosts(play.HostSelector)
		if err != nil {
			return err
		}

		// TODO: the way it is done now, host variable will not persist across plays
		// Build initial host variables by looping though groups it belongs to Add variables of that group
		// Precedence:
		// -- group vars (ordered alphabetically )
		// -- host params (from inventory)
		// -- TODO: cli
		for _, host := range hosts {
			for _, group := range host.Groups {
				if vars, found := groupVars[group]; found {
					host.Vars.Add(vars)
				}
			}
			// Override group variables with host params from inventory
			for k,v := range host.Params {
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
					result, err  := templating.Assert(task.When, host.Vars)
					if err != nil {
						return err
					}
					if !result {
						log("Skipping task [%s] on host %s\n", task.Name, host)
						continue
					}
				}

				log("\n### Running task [%s] on host %s\n", task.Name, host)
				if host.Transport == nil {
					host.Transport = transport.CreateSSHTransport(host.Params)
				}
				r := task.Module.Run(host.Transport, host.Vars)
				log("Module exec: %s", r)
				if !r.Result {
					return fmt.Errorf("module execution failed")
				}
				// register module output as variable
				if task.Register != "" {
					// TODO: module execusion result should probably be wrapped to Add extra information such as
					// execution time, module name...
					host.Vars[task.Register] = r
				}
			}
		}

	}
	return nil
}

func log(msg string, args... interface{}) {
	logging.Info(msg, args...)
}