package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	VaultPath []string
	ArchivePath string
	ArchiveType uint8
}


func main() {
	start := time.Now()
	// config_path, err := filepath.Abs("~/.config/go-archive-it/config.yaml")
	config_path, _ := filepath.Abs("./test-conf/config.yaml") // test path

	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		log.Print("creating config directory")

		err := os.MkdirAll(filepath.Dir(config_path), os.ModePerm)
		if err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}

		config := Config{
			VaultPath: []string{"~/notes", "~/dev"},
			ArchivePath: "./archive",
			ArchiveType: 2,
		}

		c, err := yaml.Marshal(config)

		if err != nil {
			log.Fatalf("Failed to serialize data: %v", err)
		}

		err = ioutil.WriteFile(config_path, c, os.ModeAppend | 0664)
		if err != nil {
			log.Fatalf("Unable to write file: %v", err)
		}

		log.Printf("Config Created at %s, make any neccesary changes and run the program again", config_path)
		os.Exit(0)
	}

	elapsed := time.Since(start)
	log.Printf("Archive(s) created in [[ %s ]]", elapsed)
}
