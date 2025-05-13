package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Configuration représente la structure du fichier de configuration YAML.
type Configuration struct {
	KeyPath    string            `yaml:"keypath"`
	Port       string            `yaml:"port"`
	SSHVersion string            `yaml:"version"`
	Users      map[string]string `yaml:"users"`
	NoPassDump []string          `yaml:"nopassdump"`
	Hosts      map[string]string `yaml:"hosts"`
}

func readconfig(configfile string) (config Configuration) {

	// Lire le contenu du fichier.
	data, err := os.ReadFile(configfile)
	if err != nil {
		fmt.Printf("Error reading config file %s: %v\n", configfile, err)
		return
	}

	// Désérialiser le contenu YAML dans la structure.
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return
	}

	return
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ipcut(ipport string) string {
	index := strings.Index(ipport, ":")
	return ipport[:index]
}
