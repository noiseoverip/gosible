package modules

import (
	"ansiblego/templating"
	"ansiblego/transport"
	"bytes"
	"fmt"
	"os"
)

// Command implements module interface and executes CLI commands on transport layer
// Pipes are not supported at this point
type Template struct {
	// Source file on local machine
	Src string
	// Destination file on remote machine
	Dest string
}

//
//	- template:
//		src: $localPath
//		dest: $remotePath
//
func LoadTemplate(args map[string]string) Module {
	return &Template{Src: args["src"], Dest:  args["dest"]}
}

func(t *Template) Run(transport transport.Transport, vars map[string]interface{}) *ModuleExecResult {
	// Render source file path
	sourcePath, err := templating.TemplateExec(t.Src, vars)
	if err != nil {
		return ErrorModuleConfig("failed to determine template source path: %v", err)
	}

	// Search for template in playbook dir. Should be extended once we add support for roles
	templateSrcFile, err := os.Open(sourcePath)
	if err != nil {
		return ErrorModuleConfig("failed to load template file: %v", err)
	}
	buf := new(bytes.Buffer)
	bytesRead, err := buf.ReadFrom(templateSrcFile)
	if bytesRead < 1 {
		fmt.Printf("WARN: template %s looks empty", sourcePath)
	}


	// Render it
	templateRendered, err := templating.TemplateExec(string(buf.Bytes()), vars)
	if err != nil {
		return ErrorModuleConfig("failed to render template: %v", err)
	}

	// Render destination path
	destinationPath, err := templating.TemplateExec(t.Dest, vars)
	if err != nil {
		return ErrorModuleConfig("failed to determine template destination path: %v", err)
	}
	fmt.Printf("\nTemplate:\n" +
		"\t\tsource: %s\n" +
		"\t\tdest: %s\n" +
		"\t\traw:\n" +
		"%s\n>>> Raw template end\n" +
		"\t\trendered:\n" +
		">>> Template start\n%s\n>>>> Template end\n\n", sourcePath, destinationPath, string(buf.Bytes()), templateRendered)

	resultCode, stdout, stdr := transport.Exec("echo", "not sure how to copy it yet")

	return &ModuleExecResult{ Result: resultCode == 0, StdOut: stdout, StdErr: stdr}
}
