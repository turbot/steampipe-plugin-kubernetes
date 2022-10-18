package main

import (
	"github.com/turbot/steampipe-plugin-kubernetes/kubernetes"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: kubernetes.Plugin})
}
