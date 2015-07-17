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

	projects    = app.Command("projects", "List projects under an orginization.")
	builds      = app.Command("builds", "List latest builds for the current project.")
	logs        = app.Command("logs", "Retrieve the logs for the current projects last build.")
	buildNumber = logs.Arg("number", "supply a build number to retrieve the logs for.").Default("").String()
	open        = app.Command("open", "Open builds list in your browser for the current project.")
	setup       = app.Command("setup", "Configure the buildkite cli with a new token.")
)

func main() {

	kingpin.Version(Version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case projects.FullCommand():
		bk.SetHttpDebug(*debugHTTP)
		kingpin.FatalIfError(commands.ProjectList(*quiet), "List projects failed")
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
