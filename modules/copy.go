package modules

import (
	"ansiblego/templating"
	"ansiblego/transport"
)

type Copy struct {
	// Source file on local machine
	Src string
	// Destination file on remote machine
	Dest string
	// Permission mode to apply on destination, example: 0644
	Mode string
	// File owner to be set with chmod
	Owner string
}

func LoadCopy(args map[string]string) Module {
	module := &Copy{Src: args["src"], Dest:  args["dest"]}
	// Optional attributes
	module.Mode = args["mode"]
	module.Owner, _ = args["owner"]
	// Default values
	if module.Mode == "" {
		module.Mode = "0600"
	}
	return module
}

func(m *Copy) Run(transport transport.Transport, vars map[string]interface{}) *ModuleExecResult {
	// Render source file path
	sourcePath, err := templating.TemplateExec(m.Src, vars)
	if err != nil {
		return ErrorModuleConfig("failed to determine source path: %v", err)
	}
	// Render destination path
	destinationPath, err := templating.TemplateExec(m.Dest, vars)
	if err != nil {
		return ErrorModuleConfig("failed to determine destination path: %v", err)
	}

	err = transport.SendFileToRemote(sourcePath, destinationPath, m.Mode)
	if err != nil {
		return &ModuleExecResult{ Result: false, StdOut: "", StdErr: err.Error()}
	}
	return &ModuleExecResult{ Result: true, StdOut: "", StdErr: ""}
}
