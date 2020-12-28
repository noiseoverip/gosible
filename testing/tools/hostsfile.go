package tools

import (
	"os"
	"text/template"
)

func RenderHostsFile(templatePath, outputPath, hostIPAddress string) error {
	hostsFileTemplate, err := template.New("hosts_template").ParseFiles(templatePath)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	err = hostsFileTemplate.Execute(f, map[string]string{"targetHost": hostIPAddress})
	if err != nil {
		return err
	}

	return f.Close()
}
