package main

import (
	"os"

	"github.com/wolfeidau/buildkite-cli/commands"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	app   = kingpin.New("bk", "A command-line interface for buildkite.com.")
	quiet = kingpin.Flag("quiet", "Only display numeric IDs").Bool()

	projects = app.Command("projects", "List projects under an orginization.")
	builds   = app.Command("builds", "List latest builds for the current project.")
	open     = app.Command("open", "Open builds list in your browser for the current project.")
	setup    = app.Command("setup", "Configure the buildkite cli with a new token.")
)

func main() {

	kingpin.Version(Version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case projects.FullCommand():
		kingpin.FatalIfError(commands.ProjectList(*quiet), "List projects failed")
	case builds.FullCommand():
		kingpin.FatalIfError(commands.BuildsList(*quiet), "List builds failed")
	case open.FullCommand():
		kingpin.FatalIfError(commands.Open(), "Open failed")
	case setup.FullCommand():
		kingpin.FatalIfError(commands.Setup(), "Setup failed")
	default:
		kingpin.UsageErrorf("missing sub command.")
	}

}
