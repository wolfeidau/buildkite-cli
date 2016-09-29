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

// LocatePipeline the pipeline which represents the current director.
func LocatePipeline(pipelines []bk.Pipeline) *bk.Pipeline {

	// git dem remotes
	remotes, err := git.Remotes()

	if err != nil {
		return nil
	}

	for _, p := range pipelines {
		for _, r := range remotes {
			if ok, gitRepo := GitRemoteMatch(r); ok {
				s := gitRepo.String()
				if strings.Contains(*p.Repository, s) {
					return &p
				}
			}

		}
	}

	return nil
}

type GitCmd interface {
	Remotes() ([]string, error)
}

type gitCmd struct{}

// Remotes locate the remotes for the current pipeline
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
