package frontendmockerfile

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ApiVersion string        `yaml:"apiVersion"`
	Images     []ConfigImage `yaml:"images"`
}

type ConfigImage struct {
	From   string `yaml:"from"`
	Parent string `yaml:"parent"`

	External []*ExternalFile `yaml:"external"`

	WorkDir string   `yaml:"workdir"`
	Steps   []string `yaml:"steps"`
	Output  []string `yaml:"output"`

	Package *Package `yaml:"package"`
}

type Package struct {
	Repo    []string `yaml:"repo"`
	Gpg     []string `yaml:"gpg"`
	Install []string `yaml:"install"`
}

type ExternalFile struct {
	Source      string `yaml:"src"`
	Destination string `yaml:"dst"`
	Sha256      string `yaml:"sha256"`

	Install []string `yaml:"install"`
}

func NewFromFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %s", err)
	}

	return NewFromBytes(contents)
}

func NewFromBytes(b []byte) (*Config, error) {
	c := &Config{}
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %s", err)
	}

	return c, nil
}
