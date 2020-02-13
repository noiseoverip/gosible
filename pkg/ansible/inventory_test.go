package ansible

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

var inventoryExample = []byte(`
host1 param1=param1value
host2

[group1]
host1

[group2]
host2

[groupAllHosts]
host1
host2

[groupNoHosts]

`)

func TestReadParams(t *testing.T) {
	var r io.Reader
	r = bytes.NewReader(inventoryExample)

	inv := new(Inventory)
	err := ReadInventory(r, inv)
	failIfError(t, err)
	assert.NotNil(t, inv)
	assert.NotEmpty(t, inv)

	allGroup, ok := inv.groupByName("all")
	assert.True(t, ok)
	assert.Exactly(t, 2, len(allGroup.Hosts))
}

func TestReadInventory(t *testing.T) {
	var r io.Reader
	r = bytes.NewReader(inventoryExample)

	inv := new(Inventory)
	err := ReadInventory(r, inv)
	failIfError(t, err)
	assert.NotNil(t, inv)
	assert.NotEmpty(t, inv)

	allGroup, ok := inv.groupByName("all")
	assert.True(t, ok)
	assert.Exactly(t, 2, len(allGroup.Hosts))
}
