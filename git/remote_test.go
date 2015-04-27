package git

import (
	"reflect"
	"testing"
)

func TestParseRemote(t *testing.T) {

	data := []struct {
		input  string
		result *GitRepo
		match  bool
	}{
		{
			"origin	https://github.com/wolfeidau/buildkite-cli.git (fetch)",
			&GitRepo{URL: "https://github.com/wolfeidau/buildkite-cli.git", Owner: "wolfeidau", Name: "buildkite-cli"},
			true,
		},
		{
			"https://github.com/wolfeidau/",
			nil,
			false,
		},
		{
			"origin	git@github.com:wolfeidau/buildkite-cli.git (fetch)",
			&GitRepo{URL: "git@github.com:wolfeidau/buildkite-cli.git", Owner: "wolfeidau", Name: "buildkite-cli"},
			true,
		},
	}

	for _, datum := range data {
		ok, repo := GitRemoteMatch(datum.input)

		if datum.match != ok {
			t.Errorf("expected %v got %v", datum.match, ok)
		}

		if !reflect.DeepEqual(datum.result, repo) {
			t.Errorf("expected %v got %v", datum.result, repo)
		}

	}

}
