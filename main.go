package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/serve"

	"github.com/people-data-labs/steampipe-plugin-pdl/pdl"
)

func main() {
	serve.Serve(&serve.Options{PluginFunc: pdl.Plugin})
}
