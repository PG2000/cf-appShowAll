package main

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
	"fmt"
	"os"
	"strings"
)

const OPEN_PROTOCOL = "https"

type AppShowAll struct {
	ui terminal.UI
}

func (c *AppShowAll) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), nil)

	s, e := cliConnection.GetApps()

	if e != nil {
		c.ui.Failed(e.Error())
	}

	table := c.ui.Table([]string{"Name", "Routes", "Bound Services", "Bound Routings"})
	for _, v := range s {

		model, e := cliConnection.GetApp(v.Name)

		if e != nil {
			c.ui.Failed(e.Error())
		}

		table.Add(v.Name, c.join(getRoutes(v.Routes)), c.join(getServices(model.Services)), string(v.TotalInstances))
	}

	table.Print()

}

func (c *AppShowAll) join(toJoin []string) string {
	return strings.Join(toJoin, "\n")
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
		route = append(route, fmt.Sprintf("%s://%s.%s", OPEN_PROTOCOL, v.Host, v.Domain.Name))
	}

	return route
}

func (c *AppShowAll) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "app-show-all",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "app-show-all",
				HelpText: "show all service instance bindings, route bindings",
				UsageDetails: plugin.Usage{
					Usage:   "app-show-all",
					Options: map[string]string{},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(AppShowAll))
}
