package modules

import (
	"ansiblego/pkg/templating"
	"path"
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
	module := &Copy{Src: args["src"], Dest: args["dest"]}
	// Optional attributes
	module.Mode = args["mode"]
	module.Owner = args["owner"]
	// Default values
	if module.Mode == "" {
		module.Mode = "0600"
	}
	return module
}

func (m *Copy) Run(ctx Context, host *Host) *ModuleExecResult {
	// Render source file path
	sourcePath, err := templating.TemplateExec(m.Src, host.Vars)
	if err != nil {
		return ErrorModuleConfig("failed to determine source path: %v", err)
	}
	sourcePath = path.Join(ctx.PlaybookDir, sourcePath)
	// Render destination path
	destinationPath, err := templating.TemplateExec(m.Dest, host.Vars)
	if err != nil {
		return ErrorModuleConfig("failed to determine destination path: %v", err)
	}

	err = host.Transport.SendFileToRemote(sourcePath, destinationPath, m.Mode)
	if err != nil {
		return &ModuleExecResult{Result: false, StdOut: "", StdErr: err.Error()}
	}
	return &ModuleExecResult{Result: true, StdOut: "", StdErr: ""}
}
