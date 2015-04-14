package git

import (
	"reflect"
	"testing"

	bk "github.com/wolfeidau/go-buildkite/buildkite"
)

type mockGitCmd struct {
	result []string
}

func (gc *mockGitCmd) Remotes() ([]string, error) {
	return gc.result, nil
}

func setupMockCmd(result []string) {
	git = &mockGitCmd{result}
}

func TestLocateProject(t *testing.T) {

	data := []struct {
		input    []string
		projects []bk.Project
		want     *bk.Project
	}{
		{
			[]string{"origin	git@github.com:wolfeidau/go-buildkite.git (fetch)"},
			[]bk.Project{
				bk.Project{ID: bk.String("123"), Repository: bk.String("git@github.com:wolfeidau/go-buildkite.git")},
				bk.Project{ID: bk.String("345"), Repository: bk.String("git@github.com:wolfeidau/someother.git")},
			},
			&bk.Project{ID: bk.String("123"), Repository: bk.String("git@github.com:wolfeidau/go-buildkite.git")},
		},
	}

	for _, d := range data {
		setupMockCmd(d.input)

		p := LocateProject(d.projects)

		if !reflect.DeepEqual(p, d.want) {
			t.Errorf("LocateProject returned %+v, want %+v", p, d.want)
		}
	}

}
