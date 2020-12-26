package pkg

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

type Playbook struct {
	Plays []Play
	Dir   string
}

// ReadPlaybook loads playbook yml file into struct
func ReadPlaybook(in io.Reader, playbook *Playbook) error {
	buf := new(bytes.Buffer)
	bytesRead, err := buf.ReadFrom(in)
	if bytesRead < 1 {
		return fmt.Errorf("empty file")
	} else if err != nil {
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
