package main

import (
	"github.com/turbot/steampipe-plugin-kubernetes/kubernetes"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: kubernetes.Plugin})
}
