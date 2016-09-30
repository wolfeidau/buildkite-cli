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

func TestLocatePipeline(t *testing.T) {

	data := []struct {
		input    []string
		pipelines []bk.Pipeline
		want     *bk.Pipeline
	}{
		{
			[]string{"origin	git@github.com:wolfeidau/go-buildkite.git (fetch)"},
			[]bk.Pipeline{
				bk.Pipeline{ID: bk.String("123"), Repository: bk.String("git@github.com:wolfeidau/go-buildkite.git")},
				bk.Pipeline{ID: bk.String("345"), Repository: bk.String("git@github.com:wolfeidau/someother.git")},
			},
			&bk.Pipeline{ID: bk.String("123"), Repository: bk.String("git@github.com:wolfeidau/go-buildkite.git")},
		},
	}

	for _, d := range data {
		setupMockCmd(d.input)

		p := LocatePipeline(d.pipelines)

		if !reflect.DeepEqual(p, d.want) {
			t.Errorf("LocatePipeline returned %+v, want %+v", p, d.want)
		}
	}

}
