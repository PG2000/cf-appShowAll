package main_test

import (
	"code.cloudfoundry.org/cli/cf/util/testhelpers/rpcserver"
	"code.cloudfoundry.org/cli/cf/util/testhelpers/rpcserver/rpcserverfakes"
	"code.cloudfoundry.org/cli/plugin/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os/exec"
)

var _ = Describe("AppShowAll", func() {
	const validPluginPath = "./appShowAll"

	var (
		rpcHandlers *rpcserverfakes.FakeHandlers
		ts          *rpcserver.TestServer
		err         error
	)

	BeforeEach(func() {
		rpcHandlers = new(rpcserverfakes.FakeHandlers)
		ts, err = rpcserver.NewTestRPCServer(rpcHandlers)
		Expect(err).NotTo(HaveOccurred())

		rpcHandlers.CallCoreCommandStub = func(_ []string, retVal *bool) error {
			*retVal = true
			return nil
		}

		//set rpc.GetOutputAndReset to return empty string; this is used by CliCommand()/CliWithoutTerminalOutput()
		rpcHandlers.GetOutputAndResetStub = func(_ bool, retVal *[]string) error {
			*retVal = []string{"{}"}
			return nil
		}
	})

	JustBeforeEach(func() {
		err = ts.Start()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ts.Stop()
	})

	Describe("list-apps", func() {
		const tableHeader = "Name   Routes   Bound Services   Bound Routings"
		const command = "apps-show-all"

		Context("no apps given", func() {
			It("app-show-all says no apps in space when no apps given", func() {
				args := []string{ts.Port(), command}
				session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
				session.Wait()
				Expect(err).NotTo(HaveOccurred())
				Expect(session).ToNot(gbytes.Say(tableHeader))
				Expect(session).To(gbytes.Say("No apps in space"))
			})
		})

		When("show table with services for apps ", func() {
			Context("apps available with bound services", func() {
				BeforeEach(func() {
					rpcHandlers.GetAppsStub = func(args string, retVal *[]plugin_models.GetAppsModel) error {
						*retVal = exampleAppsModel()
						return nil
					}

					rpcHandlers.GetAppStub = func(appName string, retVal *plugin_models.GetAppModel) error {
						*retVal = exampleAppModel()
						return nil
					}
				})

				It("lists all apps with bounded services", func() {
					args := []string{ts.Port(), command}
					session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
					session.Wait()
					Expect(err).NotTo(HaveOccurred())
					Expect(session).To(gbytes.Say("app1   https://app1.cf.local   autoscaler, example-service"))
				})
			})
		})
	})

})

func exampleAppsModel() []plugin_models.GetAppsModel {
	const appName = "app1"
	const domain = "cf.local"
	return []plugin_models.GetAppsModel{
		{
			Name:             appName,
			Guid:             "",
			State:            "",
			TotalInstances:   1,
			RunningInstances: 1,
			Memory:           0,
			DiskQuota:        0,
			Routes:           exampleAppRouteSummary(appName, domain),
		},
	}
}

func exampleAppRouteSummary(hostname string, domain string) []plugin_models.GetAppsRouteSummary {
	return []plugin_models.GetAppsRouteSummary{
		struct {
			Guid   string
			Host   string
			Domain plugin_models.GetAppsDomainFields
		}{
			"",
			hostname,
			struct {
				Guid                   string
				Name                   string
				OwningOrganizationGuid string
				Shared                 bool
			}{"", domain, "", false},
		},
	}
}

func exampleAppModel() plugin_models.GetAppModel {
	return plugin_models.GetAppModel{
		Guid:                 "",
		Name:                 "app1",
		BuildpackUrl:         "",
		Command:              "",
		DetectedStartCommand: "",
		DiskQuota:            0,
		EnvironmentVars:      nil,
		InstanceCount:        0,
		Memory:               0,
		RunningInstances:     0,
		HealthCheckTimeout:   0,
		State:                "",
		SpaceGuid:            "",
		PackageUpdatedAt:     nil,
		PackageState:         "",
		StagingFailedReason:  "",
		Stack:                nil,
		Instances:            nil,
		Routes:               []plugin_models.GetApp_RouteSummary{},
		Services: []plugin_models.GetApp_ServiceSummary{struct {
			Guid string
			Name string
		}{Guid: "dd", Name: "autoscaler"}, struct {
			Guid string
			Name string
		}{Guid: "ee", Name: "example-service"}},
	}

}
