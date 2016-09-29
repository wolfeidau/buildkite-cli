package commands

import (
	"bytes"
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

var (
	pipelineColumns = []string{"ID", "NAME", "BUILD", "BRANCH", "MESSAGE", "STATE", "FINISHED"}
	jobColumns     = []string{"NAME", "STARTED", "FINISHED", "STATE"}
	buildColumns   = []string{"PROJECT", "NUMBER", "BRANCH", "MESSAGE", "STATE", "COMMIT"}

	pipelineOrgRegex = regexp.MustCompile(`\/organizations\/([\w_-]+)\/`)
)

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

// Get List of Pipelines for all the orginizations.
func (cli *bkCli) pipelineList(quietList bool) error {

	t := time.Now()

	pipelines, err := cli.listPipelines()

	if err != nil {
		return err
	}

	if quietList {
		for _, proj := range pipelines {
			fmt.Printf("%-36s\n", *proj.ID)
		}
		return nil // we are done
	}

	tb := table.New(pipelineColumns)
	vals := make(map[string]interface{})

	for _, proj := range pipelines {
		vals = utils.ToMap(pipelineColumns, []interface{}{*proj.ID, *proj.Name, 0, "", "", "", ""})
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

	pipelines, err := cli.listPipelines()

	if err != nil {
		return err
	}

	// did we locate a pipeline
	pipeline := git.LocatePipeline(pipelines)

	if pipeline != nil {
		fmt.Printf("Listing for pipeline = %s\n\n", *pipeline.Name)

		org := extractOrg(*pipeline.URL)

		builds, _, err = cli.client.Builds.ListByPipeline(org, *pipeline.Slug, nil)

	} else {
		utils.Check(fmt.Errorf("Failed to locate the buildkite pipeline using git.")) // TODO tidy this up
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
		vals := utils.ToMap(buildColumns, []interface{}{*build.Pipeline.Name, *build.Number, *build.Branch, *build.Message, *build.State, *build.Commit})
		tb.AddRow(vals)
	}

	tb.Markdown = true
	tb.Print()

	fmt.Printf("\nTime taken: %s\n", time.Now().Sub(t))

	return nil
}

func (cli *bkCli) openPipelineBuilds() error {

	pipelines, err := cli.listPipelines()

	if err != nil {
		return err
	}

	// did we locate a pipeline
	pipeline := git.LocatePipeline(pipelines)

	if pipeline != nil {
		fmt.Printf("Opening pipeline = %s\n\n", *pipeline.Name)

	} else {
		utils.Check(fmt.Errorf("Failed to locate the buildkite pipeline using git.")) // TODO tidy this up
		return nil
	}

	org := extractOrg(*pipeline.URL)

	pipelineURL := fmt.Sprintf("https://buildkite.com/%s/%s/builds/last", org, *pipeline.Slug) // TODO URL should come from REST interface

	args, err := utils.BrowserLauncher()

	utils.Check(err) // TODO tidy this up

	cmd := cmd.New(args[0])

	args = append(args, pipelineURL)

	cmd.WithArgs(args[1:]...)

	_, err = cmd.CombinedOutput()

	return err
}

func (cli *bkCli) tailLogs(number string) error {

	pipelines, err := cli.listPipelines()

	if err != nil {
		return err
	}

	// did we locate a pipeline
	pipeline := git.LocatePipeline(pipelines)

	if pipeline != nil {
		fmt.Printf("Opening pipeline = %s\n\n", *pipeline.Name)

	} else {
		utils.Check(fmt.Errorf("Failed to locate the buildkite pipeline using git.")) // TODO tidy this up
		return nil
	}

	if number == "" {

	}

	ok, j := cli.getLastJob(pipeline, number)
	if ok {

		tb := table.New(jobColumns)

		vals := utils.ToMap(jobColumns, []interface{}{*j.Name, *j.StartedAt, *j.FinishedAt, *j.State})
		tb.AddRow(vals)
		tb.Markdown = true
		tb.Print()

		fmt.Println()

		req, err := cli.client.NewRequest("GET", *j.RawLogsURL, nil)

		if err != nil {
			return err
		}
		buffer := new(bytes.Buffer)

		_, err = cli.client.Do(req, buffer)

		if err != nil {
			return err
		}

		fmt.Printf("%s\n", string(buffer.Bytes()))
	}

	return nil
}

func (cli *bkCli) getLastJob(pipeline *bk.Pipeline, number string) (bool, *bk.Job) {
	org := extractOrg(*pipeline.URL)

	build, _, err := cli.client.Builds.Get(org, *pipeline.Slug, number)

	if err != nil {
		return false, nil
	}

	jobs := build.Jobs

	if len(jobs) == 0 {
		return false, nil
	}

	j := jobs[len(jobs)-1]

	return true, j
}

func (cli *bkCli) setup() error {
	return cli.config.PromptForConfig()
}

func (cli *bkCli) listPipelines() ([]bk.Pipeline, error) {
	var pipelines []bk.Pipeline

	orgs, _, err := cli.client.Organizations.List(nil)

	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		projs, _, err := cli.client.Pipelines.List(*org.Slug, nil)

		if err != nil {
			return nil, err
		}

		pipelines = append(pipelines, projs...)
	}

	return pipelines, nil
}

func newClient(config *config.Config) (*bk.Client, error) {

	if config.OAuthToken == "" {
		err := config.PromptForConfig()
		if err != nil {
			return nil, err
		}
	}

	tconf, err := bk.NewTokenConfig(config.OAuthToken, config.Debug)

	if err != nil {
		return nil, err
	}

	return bk.NewClient(tconf.Client()), nil
}

// PipelineList just get a list of pipelines
func PipelineList(quietList bool) error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.pipelineList(quietList)
}

// BuildsList retrieve a list of builds for the current pipeline using the git remote to locate it.
func BuildsList(quietList bool) error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.buildList(quietList)
}

// LogsList retrieve the logs for the last build using the supplied build number
func LogsList(number string) error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.tailLogs(number)
}

// Open buildkite pipeline for the current pipeline using the git remote to locate it.
func Open() error {
	cli, err := newBkCli()
	if err != nil {
		return err
	}

	return cli.openPipelineBuilds()
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
	m := pipelineOrgRegex.FindStringSubmatch(url)

	if len(m) == 2 {
		return m[1]
	}

	return ""
}

func toString(str *string) string {
	return *str
}

func valString(thing interface{}) string {
	if thing == nil {
		return ""
	}
	return fmt.Sprintf("%s", thing)
}
