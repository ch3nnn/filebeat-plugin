package main

import (
	"os"

	_ "filebeat-plugin/pkg/processors/actions"
	"github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
)

func main() {
	if err := cmd.Filebeat(inputs.Init, cmd.FilebeatSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
