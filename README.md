![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/korbexmachina/go-archive-it?style=for-the-badge)
![GitHub commits since latest release (by SemVer including pre-releases)](https://img.shields.io/github/commits-since/korbexmachina/go-archive-it/latest?style=for-the-badge)
![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/korbexmachina/go-archive-it?style=for-the-badge)

# Go Archive It!

A lightweight archive management utility with YAML based configuration.

In a world where your files and data are increasingly stored and controlled by large corporations, there are those who wish to take that power back. However, with great power comes great responsibility, and it is important to practice good data stewardship. A good place to start is by keeping rolling backups, a process which _Go Archive It_ aims to make as straightforward and painless as possible.

## Documentation

Full documentation available [here](https://korbexmachina.github.io/go-archive-it/)

## Usage

- I reccommend using a cronjob to automate your archives
  - If you are using the homebrew installation, that might look something like this:
    ```
    # Run go-archive-it at 9am every day, output logs to ~/logs/go-archive-it.txt
    0 9 * * * /usr/local/Cellar/go-archive-it/0.0.7/bin/go-archive-it >> ~/logs/go-archive-it.txt 2>&1
    ```
  - If you built the project from source, I will assume you know what you are doing, but the only real change will be the path to the binary
- The program looks for a config file at `~/.config/go-archive-it/config.yaml`
  - If the config does not exist, the program will create it with the following default contents:
    ```yaml
    vaultpath:
      - ~/example
      - ~/directories
    archivepath: ~/archive
    archivetype: 1
    retention: 10
    ```
  - Let's break this down
    - `vaultpath` is a list of directories to be archived (add as many as you need, they will be archived asynchronously)
    - `archivepath` is the path to the directory where your archives will be stored
    - `archivetype` can be set to `0` for uncompressed `.tar` archives or `1` for `.tar.gz`
    - `retention` is the number of archives you want to keep at any given time for each of the directories in vaultpath (it is stored as an 8 bit integer, so it must be less than 256)
   
### Arguments

- `-h, help`
  - Display the help message
- `-e, ext`
  - Runs the program based on a file at `~/.config/go-archive-it/ext.yaml`
    - Configuration is the same as running the program with the default path
  - Intended for archiving onto an external drive with preconfigured options
- `-i, init [NAME]`
  - Initialize a named config file at `~/.config/go-archive-it/[NAME].yaml`
- `-p, path [NAME]`
  - Run the program withthe named config file at `~/.config/go-archive-it/[NAME].yaml`
 
### Help
```
Usage: go-archive-it [OPTION] ...
---------------------------------
-h, help                Display this help message
-e, ext                 Use external config file (~/.config/go-archive-it/ext.yaml)
-i, init [NAME]         Initialize named config file (~/.config/go-archive-it/[NAME].yaml)
-p, path [NAME]         Use named config file (~/.config/go-archive-it/[NAME].yaml)
---------------------------------
Running with no arguments will use the default config file (~/.config/go-archive-it/config.yaml)
```

## Installation

### Homebrew (MacOS)

```
brew tap korbexmachina/tap
brew install go-archive-it
```

## Roadmap

- ~~`-e` as an alternative to the `ext` argument~~
- ~~`-p` argument for passing an arbitrary filename for configuration, to allow for as many user configurations as needed~~
- ~~`-h` argument for help~~
- ~~`init`/`-i` argument for initializing config files with arbitrary names~~
- Default config initialization option (templates)
- Archive the path beggining at the directory being archived rather than including the directories above it
  - The program does currently ignore their other contents, but the nesting is still mildly annoying when accessing the archives

## Dependencies
[yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) - see `notice.md` for details
