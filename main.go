package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	app = kingpin.New("bk", "A command-line interface for buildkite.com.")

	credsFile = app.Flag("creds", "Credentials file.").Default("~/.buildkite").String()
	list      = app.Command("list", "List projects under an account.")
	account   = list.Arg("account", "Buildkite account.").Required().String()
)

func main() {
	kingpin.Version(Version)
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// List projects
	case list.FullCommand():
		kingpin.FatalIfError(applyList(), "List failed")

	}

}

func applyRequest(req *http.Request, payload interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&payload)

	return err
}

func applyList() error {

	apiKey, err := getAPIKey()

	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.buildkite.com/v1/accounts/%s/projects?api_key=%s", *account, apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	var payload interface{}

	err = applyRequest(req, &payload)

	if err != nil {
		return err
	}

	w := &tabwriter.Writer{}

	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, "a\tb\tc\td\t.")

	//spew.Dump(payload)
	fmt.Printf("%-36s\t%-30s\t%4s\t%s\n", "id", "name", "no", "state")

	switch p := payload.(type) {
	case []interface{}:
		for _, value := range p {
			if info, ok := value.(map[string]interface{}); ok {
				if build, ok := info["featured_build"].(map[string]interface{}); ok {
					fmt.Printf("%-36s\t%-30s\t%4.0f\t%s\n", info["id"], info["name"], build["number"], build["state"])
				}
			}
		}
	}

	return nil
}

func getAPIKey() (string, error) {

	path := *credsFile

	// Check in case of paths like "~/something/"
	if path[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// Get List of Projects "https://api.buildkite.com/v1/accounts/:account/projects"
