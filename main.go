package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/korbexmachina/go-archive-it/utils"
)

func main() {
	start := time.Now()
	count := 0
	verbose := false

	configDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to resolve user config directory: %s", err)
	}

	configDir = filepath.Join(configDir, ".config")
	// log.Fatalf(configDir)
	configPath := filepath.Join(configDir, "go-archive-it/config.yaml")
	// configPath, _ := filepath.Abs("./test-conf/config.yaml") // test path

	helpMessage := `
Usage: go-archive-it [OPTION] ...
---------------------------------
-h, help		Display this help message
-e, ext			Use external config file (~/.config/go-archive-it/ext.yaml)
-i, init [NAME]		Initialize named config file (~/.config/go-archive-it/[NAME].yaml)
-p, path [NAME]		Use named config file (~/.config/go-archive-it/[NAME].yaml)
-v, verbose		Verbose logging
---------------------------------
Running with no arguments will use the default config file (~/.config/go-archive-it/config.yaml)
`

	if len(os.Args) < 2 {
		log.Print("Running with no arguments\n")
	} else {

		switch os.Args[1] {
		case "-h", "help":
			fmt.Printf(helpMessage)
			os.Exit(0)
		case "-e", "ext":
			configPath = filepath.Join(configDir, "go-archive-it/ext.yaml")
			log.Printf("Running with external config: %s", configPath)
		case "-i", "init":
			name := ""
			if len(os.Args) < 3 || os.Args[2] == "" {
				name = "go-archive-it/config"
			} else {
				name = "go-archive-it/" + os.Args[2]
			}
			utils.ConfigExists(filepath.Join(configDir, name + ".yaml"))
			os.Exit(0)
		case "-p", "path":
			name := "go-archive-it/" + os.Args[2]
			configPath = filepath.Join(configDir, name + ".yaml")
			log.Printf("Running with named config: %s", configPath)
		case "-v", "verbose":
			verbose = true
		default:
			log.Fatalf("Unknown argument: %s", os.Args[1])
		}
	}

	utils.ConfigExists(configPath)
	config := utils.LoadConfig(configPath)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := usr.HomeDir

	archivePath := config.ArchivePath
	if archivePath == "~" {
		archivePath = dir
	} else if strings.HasPrefix(archivePath, "~/") {
		archivePath = filepath.Join(dir, archivePath[2:])
	}

	var wg sync.WaitGroup
	// The loop that actually runs everything
	for _, path := range config.VaultPath {
		count++
		// Tilda expansion
		if path == "~" {
			path = dir
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(dir, path[2:])
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			utils.Archive(path, archivePath, config.ArchiveType)
		}(path)

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			err := utils.Cleanup(filepath.Join(archivePath, filepath.Base(path)), config.Retention, verbose)
			if err != nil {
				log.Fatalf("Failed to cleanup: %s", err)
			}
		}(path)
	}

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("%d Archive(s) created in [[ %f ]] seconds", count, elapsed.Seconds())
}
