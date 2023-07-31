package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Template struct {
	vault_path     []string
	archive_path   string
	archive_type   int
}


func main() {
	start := time.Now()
	// config_path, err := filepath.Abs("~/.config/go-archive-it/config.yaml")
	config_path, err := filepath.Abs("./go-archive-it/config.yaml") // test path

	if os.IsNotExist(err) {
		log.Print("creating config directory")

		os.Mkdir(filepath.Dir(config_path), os.ModePerm)
		if err != nil {
			log.Fatal("Unable to create directory: %v", err)
		}

		template := Template{
			vault_path: []string{"~/notes, ~/dev"},
			archive_path: "~/go-archive-it",
			archive_type: 2,
		}

		err := os.WriteFile("filename.txt", template, 0755)
		if err != nil {
			log.Fatal("Unable to write file: %v", err)
		}

		fmt.Printf("Config Created at %s, make any neccesary changes and run the program again", config_path)
		os.Exit(0)
	}
	elapsed := time.Since(start)
	log.Printf("Archive(s) created in [[ %s ]]", elapsed)
}
