# Documentation

## main.go

### main()

The main function that runs the program.

It begins by defining the user config directory.

```go
configDir, err := os.UserHomeDir()
if err != nil {
    log.Fatalf("Failed to resolve user config directory: %s", err)
}
configDir = filepath.Join(configDir, ".config")

```

The function then calls the following functions defined in [utils/config.go](#config.go) and sanitizes the input by expanding `~` into the user's home directory

```go
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
```

Finally the function loops though the directories specified in `vaultpath:` and calls [Archive()](#Archive()) asynchronously

```go
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
                err := utils.Cleanup(filepath.Join(archivePath, filepath.Base(path)), config.Retention)
                if err != nil {
                    log.Fatalf("Failed to cleanup: %s", err)
                }
        }(path)
}
```

Finally the function logs the execution time of the program

```go
log.Printf("%d Archive(s) created in [[ %f ]] seconds", count, elapsed.Seconds())
```

## utils

### config.go

#### type Config

A struct for storing config data, defined as follows:

```go
type Config struct {
	VaultPath []string
	ArchivePath string
	ArchiveType uint8
	Retention uint8
}
```
- `VaultPath` contains a slice of the paths to the directories that will be archived, stored as strings
- `ArchivePath` is the path to the directory where your archives will be stored
- `ArchiveType` can be set to 0 for uncompressed .tar archives or 1 for .tar.gz
- `Retention` is the number* of archives that will be kept of each directory in VaultPath at any given time for each of the directories in VaulPpath (it is stored as an 8 bit integer, so it must be less than 256)

#### ConfigExists()

- Checks if a file exists at `~/.config/go-archive-it/config.yaml`
    - if there is no such file, it creates one, containing the following
    ```yaml
    vaultpath:
        - ~/example
        - ~/directories
    archivepath: ~/archive
    archivetype: 1
    retention: 10
    ```
        - The program then exits with error code 0 and prints the following message:
        `Config Created at %s, make any neccesary changes and run the program again`
    - if there is such a file, the function returns nothing

#### LoadConfig()

Takes a path to a `.yaml` file and de-serializes it into a [Config](#struct-Config) and returns the [Config](#struct-Config)

```go
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
```

### archive.go

#### Archive(vaultPath string, archivePath string, archiveType uint8)

- Accepts three arguments a string that corresponds to the path to directory to be archived, another string that corresponds to the directory where the archives will be stored, and a uint 8 which should be either 0 or 1 to indicate what type of archive to create
- Returns nothing on success

#### Cleanup(archivePath string, retention uint8) error

- Checks if the any of the archivePath contains more than the number of files specified by retention
- If the directory exceesd this number, the oldest file in the directory will be deleted
- Returns `nil` if no errors occur, otherwise returns an error

```go
files, err := os.ReadDir(archivePath)
if len(files) < int(retention) {
    return nil
} else {
    log.Printf("Retention cap [[ %d ]] exceeded - Cleaning up %s...", retention, archivePath)
}
if err != nil {
    return err
}

oldestFile, err := files[0].Info()
if err != nil {
    return err
}
oldest := time.Now()
for _, file := range files[1:] {
    info, err := file.Info()
        if err != nil {
            return err
        }
    if info.ModTime().Before(oldest) {
        oldestFile = info
            oldest = info.ModTime()
    }
}

os.Remove(filepath.Join(archivePath, oldestFile.Name()))
    log.Printf("%s succesfully cleaned up!", archivePath)
    return nil
    }
```
