package main

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
	"fmt"
	"os"
	"strings"
)

type AppShowAll struct {
	ui terminal.UI
}

func (c *AppShowAll) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), nil)

	apps, appsError := cliConnection.GetApps()

	if appsError != nil {
		c.ui.Failed(appsError.Error())
	}

	table := c.ui.Table([]string{"Name", "Routes", "Bound Services", "Bound Routings"})

	for _, app := range apps {
		model, e := cliConnection.GetApp(app.Name)

		if e != nil {
			c.ui.Failed(e.Error())
		}

		table.Add(app.Name, join(getRoutes(app.Routes)), join(getServices(model.Services)), string(app.TotalInstances))
	}

	if len(apps) > 0 {
		table.Print()
	} else {
		c.ui.Say("No apps in space")
	}

}

func join(toJoin []string) string {
	return strings.Join(toJoin, ", ")
}

func getServices(summaries []plugin_models.GetApp_ServiceSummary) []string {
	var service []string

	for _, v := range summaries {
		service = append(service, v.Name)
	}

	return service
}

func getRoutes(summaries []plugin_models.GetAppsRouteSummary) []string {
	var route []string

	for _, v := range summaries {
		route = append(route, fmt.Sprintf("%s://%s.%s", "https", v.Host, v.Domain.Name))
	}

	return route
}

func (c *AppShowAll) GetMetadata() plugin.PluginMetadata {
	const pluginCommand = "apps-show-all"
	return plugin.PluginMetadata{
		Name: pluginCommand,
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     pluginCommand,
				HelpText: "show all service instance bindings, route bindings for all apps in current space",
				UsageDetails: plugin.Usage{
					Usage:   pluginCommand,
					Options: map[string]string{},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(AppShowAll))
}
