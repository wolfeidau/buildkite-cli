package git

import (
	"strings"

	"github.com/github/hub/cmd"
	bk "github.com/wolfeidau/go-buildkite/buildkite"
)

var git GitCmd

func init() {
	git = &gitCmd{}
}

// LocateProject the project which represents the current director.
func LocateProject(projects []bk.Project) *bk.Project {

	// git dem remotes
	remotes, err := git.Remotes()

	if err != nil {
		return nil
	}

	for _, p := range projects {
		for _, r := range remotes {
			if strings.Contains(r, *p.Repository) {
				return &p
			}
		}
	}

	return nil
}

type GitCmd interface {
	Remotes() ([]string, error)
}

type gitCmd struct{}

// Remotes locate the remotes for the current project
func (*gitCmd) Remotes() ([]string, error) {
	return gitOutput("remote", "-v")
}

func gitOutput(input ...string) (outputs []string, err error) {
	cmd := cmd.New("git")

	for _, i := range input {
		cmd.WithArg(i)
	}

	out, err := cmd.CombinedOutput()
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			outputs = append(outputs, string(line))
		}
	}

	return outputs, err
}
