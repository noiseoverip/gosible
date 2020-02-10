package modules

import (
	"ansiblego/logging"
	"ansiblego/templating"
	"ansiblego/transport"
	"bytes"
	"io/ioutil"
	"os"
)

// Command implements module interface and executes CLI commands on transport layer
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
func NewTemplate(args map[string]string) Module {
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
		logging.Info("WARN: template %s looks empty", sourcePath)
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
	logging.Info("\nTemplate:\n" +
		"\t\tsource: %s\n" +
		"\t\tdest: %s\n" +
		"\t\traw:\n" +
		"%s\n>>> Raw template end\n" +
		"\t\trendered:\n" +
		">>> Template start\n%s\n>>>> Template end\n\n", sourcePath, destinationPath, string(buf.Bytes()), templateRendered)

	// TODO: we could skip writing it to file and push it directly
	tempFile, err := ioutil.TempFile(os.TempDir(), "ansiblego-*")
	if err != nil {
		return &ModuleExecResult{ Result: false, StdOut: "", StdErr: err.Error()}
	}
	written, err := tempFile.Write([]byte(templateRendered))
	if err != nil {
		return &ModuleExecResult{ Result: false, StdOut: "", StdErr: err.Error()}
	} else if written < 1 {
		logging.Info("WARN: 0 bytes written for template")
		return &ModuleExecResult{ Result: false, StdOut: "", StdErr: ""}
	}
	logging.Info("Saved template to %s\n", tempFile.Name())

	err = transport.SendFileToRemote(tempFile.Name(), destinationPath, "0600")
	if err != nil {
		return &ModuleExecResult{ Result: false, StdOut: "", StdErr: err.Error()}
	}
	return &ModuleExecResult{ Result: true, StdOut: "", StdErr: ""}
}
