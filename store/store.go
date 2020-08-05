package store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const configSubDirName string = ".grpctl"
const configFileName string = "grpctl.json"
const defaultMode = 0755

type config struct {
	Target string `json:"target,omitempty"`
}

func configFileLocation() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home dir: %v", err)
	}

	return path.Join(dir, configSubDirName, configFileName)
}

func maybeCreateConfigFile() (*os.File, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(configFileLocation())
	if err != nil {
		os.Mkdir(path.Join(dir, configSubDirName), defaultMode)
	}

	return os.OpenFile(
		configFileLocation(),
		os.O_RDWR|os.O_CREATE,
		defaultMode,
	)
}

// SaveTarget saves the target info in a local config file
func SaveTarget(target string) error {
	f, err := maybeCreateConfigFile()
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(config{target}, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	_, err = f.WriteString("\n")
	if err != nil {
		return err
	}

	Debugf("target saved into config file at %v: %v", configFileLocation(), target)
	return nil
}

// LoadTarget loads the target address from a local config file
func LoadTarget() string {
	b, err := ioutil.ReadFile(configFileLocation())
	if err != nil {
		Debugf("failed to read config file: %v", err)
		return ""
	}

	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		Debugf("failed to parse config: %v", err)
		return ""
	}
	Debugf("loaded target from config file at %v: %v", configFileLocation(), c.Target)
	return c.Target
}
