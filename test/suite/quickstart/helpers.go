package quickstart

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"

	"github.com/jenkins-x/bdd-jx/test/helpers"

	"github.com/jenkins-x/bdd-jx/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	IncludedQuickstarts = []string{"node-http", "spring-boot-rest-prometheus-java11", "spring-boot-http-gradle", "golang-http"}
	_                   = AllQuickstartsTest()
)

// AllQuickstartsTest is responsible for running `jx get quickstarts, and creating a test for each quickstart currnetly
// available
// Individual tests can be run with `go test test/quickstart -ginkgo.focus <quickstart name>`
func AllQuickstartsTest() []bool {
	cmd := exec.Command("jx", "get", "quickstarts")
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("running jx get quickstarts, output was %s: %w", string(bytes), err))
	}
	tests := make([]bool, 0)
	for _, testQuickstartName := range IncludedQuickstarts {
		tests = append(tests, CreateQuickstartsTests(testQuickstartName))
	}
	return tests
}

// CreateQuickstartsTests creates a batch quickstart test for the given quickstart
func CreateQuickstartsTests(quickstartName string) bool {
	return createQuickstartTests(quickstartName)
}

// CreateQuickstartTest Creates quickstart tests.
func createQuickstartTests(quickstartName string) bool {
	return Describe("quickstart "+quickstartName+"\n", func() {
		var T helpers.TestOptions

		BeforeEach(func() {
			qsNameParts := strings.Split(quickstartName, "-")
			qsAbbr := ""
			for s := range qsNameParts {
				qsAbbr = qsAbbr + qsNameParts[s][:1]

			}
			applicationName := helpers.TempDirPrefix + qsAbbr + "-" + strconv.FormatInt(GinkgoRandomSeed(), 10)
			T = helpers.TestOptions{
				ApplicationName: applicationName,
				WorkDir:         helpers.WorkDir,
			}
			T.GitProviderURL()

			utils.LogInfof("Creating application %s in dir %s\n", termcolor.ColorInfo(applicationName), termcolor.ColorInfo(helpers.WorkDir))
		})

		Describe("Create a quickstart", func() {
			Context(fmt.Sprintf("by running jx create quickstart %s", quickstartName), func() {
				It("creates a new source repository and promotes it to staging", func() {
					args := []string{"create", "quickstart", "-b", "--org", T.GetGitOrganisation(), "-p", T.ApplicationName, "-f", quickstartName}

					gitProviderUrl, err := T.GitProviderURL()
					Expect(err).NotTo(HaveOccurred())
					if gitProviderUrl != "" {
						utils.LogInfof("Using Git provider URL %s\n", gitProviderUrl)
						args = append(args, "--git-provider-url", gitProviderUrl)
					}
					gitKind := os.Getenv("GIT_KIND")
					if gitKind != "" {
						args = append(args, "--git-kind", gitKind)
					}

					argsStr := strings.Join(args, " ")
					By(fmt.Sprintf("calling jx %s", argsStr), func() {
						T.ExpectJxExecution(T.WorkDir, helpers.TimeoutSessionWait, 0, args...)
					})

					applicationName := T.GetApplicationName()
					owner := T.GetGitOrganisation()
					branch := T.GetDefaultBranch()
					jobName := owner + "/" + applicationName + "/" + branch

					if T.WaitForFirstRelease() {
						//FIXME Need to wait a little here to ensure that the build has started before asking for the log as the jx create quickstart command returns slightly before the build log is available
						time.Sleep(30 * time.Second)
						By(fmt.Sprintf("waiting for the first release of %s", applicationName), func() {
							T.ThereShouldBeAJobThatCompletesSuccessfully(jobName, helpers.TimeoutBuildCompletes)

							if T.ViewPromotePRPipelines() {
								T.ViewPromotePRPipelineLog(helpers.TimeoutBuildCompletes)
							}

							T.TheApplicationIsRunningInStaging(200)
						})

					} else {
						By(fmt.Sprintf("waiting for the first successful build of %s of %s", branch, applicationName), func() {
							T.ThereShouldBeAJobThatCompletesSuccessfully(jobName, helpers.TimeoutBuildCompletes)
						})
					}

					if T.DeleteApplications() {
						utils.LogInfof("deleting the app %s", T.ApplicationName)

						args = []string{"application", "delete", "--no-source", "--repo", T.ApplicationName}
						argsStr := strings.Join(args, " ")
						By(fmt.Sprintf("calling %s to delete the application", argsStr), func() {
							T.ExpectJxExecution(T.WorkDir, helpers.TimeoutSessionWait, 0, args...)
						})

						if T.ViewPromotePRPipelines() {
							T.ViewPromotePRPipelineLog(helpers.TimeoutBuildCompletes)

							T.ViewBootJob(helpers.TimeoutBuildCompletes)
						}
					} else {
						utils.LogInfof("not the app %s", T.ApplicationName)
					}

					if T.TestPullRequest() {
						utils.LogInfof("now performing a PR to test a preview")
						By("performing a pull request on the source and asserting that a preview environment is created", func() {
							T.CreatePullRequestAndGetPreviewEnvironment(200)
						})
					}

					if T.DeleteRepos() {
						args = []string{"delete", "repo", "-b", "--github", "-o", T.GetGitOrganisation(), "-n", T.ApplicationName}
						argsStr = strings.Join(args, " ")

						By(fmt.Sprintf("calling %s to delete the repository", os.Args), func() {
							T.ExpectJxExecution(T.WorkDir, helpers.TimeoutSessionWait, 0, args...)
						})
					}
				})
			})
		})
		Describe("Create a quickstart with invalid parameters", func() {
			Context("when -p param (project name) is missing", func() {
				It("exits with signal 1", func() {
					args := []string{"create", "quickstart", "-b", "--org", T.GetGitOrganisation(), "-f", quickstartName}
					argsStr := strings.Join(args, " ")
					By(fmt.Sprintf("calling jx %s", argsStr), func() {
						T.ExpectJxExecution(T.WorkDir, helpers.TimeoutSessionWait, 1, args...)
					})
				})
			})
			Context("when -f param (filter) does not match any quickstart", func() {
				It("exits with signal 1", func() {
					args := []string{"create", "quickstart", "-b", "--org", T.GetGitOrganisation(), "-p", T.ApplicationName, "-f", "the_derek_zoolander_app_for_being_really_really_good_looking"}
					argsStr := strings.Join(args, " ")
					By(fmt.Sprintf("calling jx %s", argsStr), func() {
						T.ExpectJxExecution(T.WorkDir, helpers.TimeoutSessionWait, 1, args...)
					})
				})
			})
		})
	})
}
