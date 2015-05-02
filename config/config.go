package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/wolfeidau/buildkite-cli/utils"
	"gopkg.in/yaml.v2"
)

// TODO support different configurations for different "instances" of buildbox
const host = "api.buildkite.com"
const clientID = "buildkitecli"

var defaultConfigsFile string

func init() {
	homeDir := os.Getenv("HOME")

	if homeDir == "" {
		if u, err := user.Current(); err == nil {
			homeDir = u.HomeDir
		}
	}

	if homeDir == "" {
		utils.Check(fmt.Errorf("Can't get current user's home dir"))
	}

	defaultConfigsFile = filepath.Join(homeDir, ".buildkite", "bk")
}

// Config for the buildkite cli
type Config struct {
	OAuthToken   string `yaml:"oauth_token"`
	Orginization string `yaml:"orginization"`
	Debug        bool   `yaml:"debug"`
}

func (c *Config) PromptForConfig() (err error) {

	token := c.PromptForToken(host)

	c.OAuthToken = token

	utils.Printf("Saving configuration for %s\n", host)
	err = newConfigService().Save(configsFile(), c)

	return
}

func (c *Config) PromptForToken(host string) (token string) {
	token = os.Getenv("BUILDKITE_TOKEN")
	if token != "" {
		return
	}

	utils.Printf("Grab your token from https://buildkite.com/user/api-access-tokens\n")
	utils.Printf("%s token: ", host)
	token = c.scanLine()

	return
}

func (c *Config) PromptForOrginization(host string) (user string) {
	user = os.Getenv("BUILDKITE_ORG")
	if user != "" {
		return
	}

	utils.Printf("%s orginization: ", host)
	user = c.scanLine()

	return
}

func (c *Config) scanLine() string {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	utils.Check(scanner.Err())

	return line
}

func configsFile() string {
	configsFile := os.Getenv("BUILDKITE_CONFIG")
	if configsFile == "" {
		configsFile = defaultConfigsFile
	}

	return configsFile
}

func CurrentConfig() *Config {
	c := &Config{}
	newConfigService().Load(configsFile(), c)

	return c
}

func newConfigService() *configService {
	return &configService{}
}

type configService struct {
}

func (s *configService) Save(filename string, c *Config) error {

	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&c)

	return ioutil.WriteFile(filename, data, 0600)
}

func (s *configService) Load(filename string, c *Config) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
