package ansible

import (
	"ansiblego/templating"
	"ansiblego/transport"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
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
		// Build initial host variables by looping though groups it belongs to add variables of that group
		// Precedence:
		// -- group vars (ordered alphabetically )
		// -- host params (from inventory)
		// -- TODO: cli
		for _, host := range hosts {
			for _, group := range host.Groups {
				if vars, found := groupVars[group]; found {
					host.Vars.add(vars)
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
					conditional := fmt.Sprintf("{%% if %s %%} True {%% else %%} False {%% endif %%}", task.When)
					if outRaw, err := templating.TemplateExec(conditional, host.Vars); err == nil {
						out := strings.TrimSpace(outRaw)
						if out == "False" {
							log("Skipping task [%s] on host %s\n", task.Name, host)
							continue
						}
					} else {
						return fmt.Errorf("%v", err)
					}
				}

				log("Running task [%s] on host %s\n", task.Name, host)
				if host.Transport == nil {
					host.Transport = transport.CreateSSHTransport(host.Params)
				}
				r := task.Module.Run(host.Transport, host.Vars)
				log("Module exec: %s", r)
			}
		}

	}
	return nil
}

func log(msg string, args... interface{}) {
	fmt.Printf(msg, args...)
}