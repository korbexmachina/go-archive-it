package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	VaultPath []string
	ArchivePath string
	ArchiveType uint8
}


func ConfigExists(configPath string) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Print("creating config directory")

		err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		if err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}

		config := Config{
			VaultPath: []string{"~/notes", "~/dev"},
			ArchivePath: "./archive",
			ArchiveType: 2,
		}

		c, err := yaml.Marshal(config) // Serialize the struct

		if err != nil {
			log.Fatalf("Failed to serialize data: %v", err)
		}

		err = ioutil.WriteFile(configPath, c, os.ModeAppend | 0664)
		if err != nil {
			log.Fatalf("Unable to write file: %v", err)
		}

		log.Printf("Config Created at %s, make any neccesary changes and run the program again", configPath)
		os.Exit(0)
	}
}

func LoadConfig(configPath string) Config {
	var config Config

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Unable to parse file: %v", err)
	}

	return config
}
