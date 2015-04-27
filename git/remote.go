// two types of URL we are concerned with
//
// git@github.com:wolfeidau/buildkite-cli.git
// https://github.com/wolfeidau/buildkite-cli.git

package git

import (
	"fmt"
	"regexp"
	"strings"
)

var ( //
	remoteRegex = regexp.MustCompile(`\S+\s(.+)\s\S+`)
	urlRegex    = regexp.MustCompile(`github.com\/(?P<owner>.+)\/(?P<project>.+)`)
	gitRegex    = regexp.MustCompile(`github.com:(?P<owner>.+)\/(?P<project>.+)`)
)

// GitRepo git repo details parsed from the remote url.
type GitRepo struct {
	URL   string
	Owner string
	Name  string
}

func (gr *GitRepo) String() string {
	return fmt.Sprintf("%s/%s", gr.Owner, gr.Name)
}

func GitRemoteMatch(remote string) (bool, *GitRepo) {

	ok, url := extractGitURL(remote)

	if !ok {
		return false, nil
	}

	if ok, keys := decodeGitURLRemote(url); ok {
		r := strings.TrimSuffix(keys["project"], ".git")
		return true, &GitRepo{url, keys["owner"], r}
	}
	if ok, keys := decodeGitRemote(url); ok {
		r := strings.TrimSuffix(keys["project"], ".git")
		return true, &GitRepo{url, keys["owner"], r}
	}
	return false, nil
}

func extractGitURL(remote string) (bool, string) {
	m := remoteRegex.FindStringSubmatch(remote)

	if len(m) == 0 {
		return false, ""
	}

	return true, m[1]
}

func decodeGitURLRemote(remote string) (bool, map[string]string) {

	m := urlRegex.FindAllStringSubmatch(remote, -1)

	if len(m) == 0 {
		return false, nil
	}

	matches := m[0]
	keys := urlRegex.SubexpNames()

	md := make(map[string]string)

	for i, m := range matches {
		md[keys[i]] = m
	}

	return true, md
}

func decodeGitRemote(remote string) (bool, map[string]string) {

	m := gitRegex.FindAllStringSubmatch(remote, -1)

	if len(m) == 0 {
		return false, nil
	}

	matches := m[0]
	keys := gitRegex.SubexpNames()

	md := make(map[string]string)

	for i, m := range matches {
		md[keys[i]] = m
	}

	return true, md
}
