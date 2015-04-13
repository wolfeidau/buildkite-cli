package commands

import (
	"fmt"
	"regexp"
	"time"

	table "github.com/crackcomm/go-clitable"
	"github.com/github/hub/cmd"
	"github.com/wolfeidau/buildkite-cli/config"
	"github.com/wolfeidau/buildkite-cli/git"
	"github.com/wolfeidau/buildkite-cli/utils"
	bk "github.com/wolfeidau/go-buildkite/buildkite"
)

var projectColumns = []string{"ID", "NAME", "BUILD", "BRANCH", "MESSAGE", "STATE", "FINISHED"}
var buildColumns = []string{"PROJECT", "NUMBER", "BRANCH", "MESSAGE", "STATE", "COMMIT"}

var projectOrgRegex = regexp.MustCompile(`\/organizations\/([\w_-]+)\/`)

// BkCli manages the config and state for the buildkite cli
type bkCli struct {
	config *config.Config
	client *bk.Client
}

// NewBkCli configure the buildkite cli using the supplied config
func newBkCli() (*bkCli, error) {
	config := config.CurrentConfig()

	client, err := newClient(config)

	if err != nil {
		return nil, err
	}

	return &bkCli{config, client}, nil
}

// Get List of Projects for all the orginizations.
func (cli *bkCli) projectList(quietList bool) error {

	t := time.Now()

	projects, err := cli.listProjects()

	if err != nil {
		return err
	}

	if quietList {
		for _, proj := range projects {
			fmt.Printf("%-36s\n", *proj.ID)
		}
		return nil // we are done
	}

	tb := table.New(projectColumns)
	vals := make(map[string]interface{})

	for _, proj := range projects {
		if proj.FeaturedBuild != nil {
			vals = utils.ToMap(projectColumns, []interface{}{*proj.ID, *proj.Name, *proj.FeaturedBuild.Number, *proj.FeaturedBuild.Branch, *proj.FeaturedBuild.Message, *proj.FeaturedBuild.State, *proj.FeaturedBuild.FinishedAt})
		} else {
			vals = utils.ToMap(projectColumns, []interface{}{*proj.ID, *proj.Name, 0, "", "", "", ""})
		}
		tb.AddRow(vals)
	}
	tb.Markdown = true
	tb.Print()

	fmt.Printf("\nTime taken: %s\n", time.Now().Sub(t))

	return err
}

// List Get List of Builds
func (cli *bkCli) buildList(quietList bool) error {

	var (
		builds []bk.Build
		err    error
	)

	t := time.Now()

	projects, err := cli.listProjects()

	if err != nil {
		return err
	}

	// did we locate a project
	project := git.LocateProject(projects)

	if project != nil {
		fmt.Printf("Listing for project = %s\n\n", *project.Name)

		org := extractOrg(*project.URL)

		builds, _, err = cli.client.Builds.ListByProject(org, *project.Slug, nil)

	} else {
		utils.Check(fmt.Errorf("Failed to locate the buildkite project using git.")) // TODO tidy this up
		return nil
	}

	if err != nil {
		return err
	}

	if quietList {
		for _, build := range builds {
			fmt.Printf("%-36s\n", *build.ID)
		}
		return nil // we are done
	}

	tb := table.New(buildColumns)

	for _, build := range builds {
		vals := utils.ToMap(buildColumns, []interface{}{*build.Project.Name, *build.Number, *build.Branch, *build.Message, *build.State, *build.Commit})
		tb.AddRow(vals)
	}

	tb.Markdown = true
	tb.Print()

	fmt.Printf("\nTime taken: %s\n", time.Now().Sub(t))

	return nil
}

func (cli *bkCli) openProjectBuilds() error {

	projects, err := cli.listProjects()

	if err != nil {
		return err
	}

	// did we locate a project
	project := git.LocateProject(projects)

	if project != nil {
		fmt.Printf("Opening project = %s\n\n", *project.Name)

	} else {
		utils.Check(fmt.Errorf("Failed to locate the buildkite project using git.")) // TODO tidy this up
		return nil
	}

	org := extractOrg(*project.URL)

	projectURL := fmt.Sprintf("https://buildkite.com/%s/%s/builds", org, *project.Slug) // TODO URL should come from REST interface

	args, err := utils.BrowserLauncher()

	utils.Check(err) // TODO tidy this up

	cmd := cmd.New(args[0])

	args = append(args, projectURL)

	cmd.WithArgs(args[1:]...)

	_, err = cmd.CombinedOutput()

	return err
}

func (cli *bkCli) setup() error {
	return cli.config.PromptForConfig()
}

func (cli *bkCli) listProjects() ([]bk.Project, error) {
	var projects []bk.Project

	orgs, _, err := cli.client.Organizations.List(nil)

	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		projs, _, err := cli.client.Projects.List(*org.Slug, nil)

		if err != nil {
			return nil, err
		}

		projects = append(projects, projs...)
	}

	return projects, nil
}

func newClient(config *config.Config) (*bk.Client, error) {

	if config.OAuthToken == "" {
		err := config.PromptForConfig()
		if err != nil {
			return nil, err
		}
	}

	tconf, err := bk.NewTokenConfig(config.OAuthToken)

	if err != nil {
		return nil, err
	}

	return bk.NewClient(tconf.Client()), nil
}

// ProjectList just get a list of projects
func ProjectList(quietList bool) error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.projectList(quietList)
}

// BuildsList retrieve a list of builds for the current project using the git remote to locate it.
func BuildsList(quietList bool) error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.buildList(quietList)
}

// Open buildkite project for the current project using the git remote to locate it.
func Open() error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.openProjectBuilds()
}

// Setup configure the buildkite cli with a new token.
func Setup() error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.setup()
}

func extractOrg(url string) string {
	m := projectOrgRegex.FindStringSubmatch(url)

	if len(m) == 2 {
		return m[1]
	}

	return ""
}
