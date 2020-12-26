// build integration

package integration

import (
	"ansiblego/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

// TODO: run tests inside docker container. It might be easier to build container with out hardcoded ssh key in it,
//   quickly start and quickly kill it and run all these tests in parallel.
//

func TestBasicPlaybook(t *testing.T) {
	setup(t)
	wd, _ := os.Getwd()
	//logging.Global = logging.NewGosibleVerboseLogger(5)
	r := internal.NewRunner(path.Join(wd, "hosts"), path.Join(wd, "site.yaml"))
	err := r.Run()
	assert.NoError(t, err)
}

func setup(t *testing.T) {
	err := os.Chdir("files")
	assert.Nil(t, err)
	dir, err := os.Getwd()
	assert.Nil(t, err)
	t.Logf("current dir %s", dir)
}
