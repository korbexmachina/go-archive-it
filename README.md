# Go Archive It!

A lightweight archive management utility that can be managed with a YAML config file. I reccomend running it with a cron job, the script will not generate more than one archive per directory per day.

## What you need to know

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
    - `retention` is the number* of archives you want to keep at any givven time for each of the directories in vaultpath (it is stored as an 8 bit integer, so it must be less than 256)
