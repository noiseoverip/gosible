package ansible

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

var hostsExample = []byte(`
---

- hosts: all
  tasks:
    - name: taskPlay1
      command: "hola 1"

- hosts: all
  roles:
    - role1

# TODO: Add support for role as a map
#- hosts: all
#  roles:
#    - { name: role1 }

`)

func TestLoadPlaybook(t *testing.T) {

	var r io.Reader
	r = bytes.NewReader(hostsExample)

	playbook := new(Playbook)
	err := ReadPlaybook(r, playbook)

	failIfError(t, err)
	if len(playbook.Plays) < 1 {
		t.Fail()
	}

	for _, play := range playbook.Plays {
		assert.NotNil(t, play.HostSelector)
	}

	// Make sure play index 0 has tasks
	assert.NotNil(t, playbook.Plays[0].Tasks)
	assert.True(t, len(playbook.Plays[0].Tasks) > 0)

	// Make sure play index 1 has roles
	assert.NotNil(t, playbook.Plays[1].Roles)
	assert.True(t, len(playbook.Plays[1].Roles) > 0)

}

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
