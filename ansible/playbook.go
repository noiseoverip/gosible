package ansible

import (
	"ansiblego/transport"
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

func (playbook *Playbook) Run(inventory *Inventory) {
	for _, play := range playbook.Plays {
		log("Running play %s\n", play.Hosts)
		if group, found := inventory.Group(play.Hosts); found {
			playbook.runTasks(play.Tasks, group.Hosts)
		}
	}
}

func (*Playbook) runTasks(tasks []Task,  hosts []*Host) {
	for _, host := range hosts {
		for _, task := range tasks {
			log("Running task [%s] on host %s\n", task.Name, host)
			if host.Transport == nil {
				host.Transport = transport.CreateSSHTransport(host.Params)
			}
			r := task.Module.Run(host.Transport)
			log("Result:%v Stdout:%s", r.Result, r.StdOut)
		}
	}
}

func log(msg string, args... interface{}) {
	fmt.Printf(msg, args...)
}