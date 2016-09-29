package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/wolfeidau/buildkite-cli/commands"
	bk "github.com/wolfeidau/go-buildkite/buildkite"
)

var (
	app       = kingpin.New("bk", "A command-line interface for buildkite.com.")
	quiet     = app.Flag("quiet", "Only display numeric IDs").Bool()
	debugHTTP = app.Flag("debug-http", "Display detailed HTTP debugging").Bool()

	pipelines    = app.Command("pipelines", "List pipelines under an orginization.")
	builds      = app.Command("builds", "List latest builds for the current pipeline.")
	logs        = app.Command("logs", "Retrieve the logs for the current pipelines last build.")
	buildNumber = logs.Arg("number", "supply a build number to retrieve the logs for.").Default("").String()
	open        = app.Command("open", "Open builds list in your browser for the current pipeline.")
	setup       = app.Command("setup", "Configure the buildkite cli with a new token.")
)

func main() {

	kingpin.Version(Version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case pipelines.FullCommand():
		bk.SetHttpDebug(*debugHTTP)
		kingpin.FatalIfError(commands.PipelineList(*quiet), "List pipelines failed")
	case builds.FullCommand():
		kingpin.FatalIfError(commands.BuildsList(*quiet), "List builds failed")
	case logs.FullCommand():
		kingpin.FatalIfError(commands.LogsList(*buildNumber), "List builds failed")
	case open.FullCommand():
		kingpin.FatalIfError(commands.Open(), "Open failed")
	case setup.FullCommand():
		kingpin.FatalIfError(commands.Setup(), "Setup failed")
	default:
		kingpin.Errorf("missing sub command.")
	}

}
