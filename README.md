# Lia CLI
```bash
lia is a CLI tool for easier development of Lia bots.

Usage:
  lia [flags]
  lia [command]

Available Commands:
  bot         Create new bot
  compile     Compiles/prepares bot in specified dir
  fetch       Fetches a bot from url and sets a new name
  generate    Generates a game
  help        Help about any command
  play        Compiles and generates a game between bot1 and bot2
  replay      Runs a replay viewer
  settings    Views the user's settings
  tutorial    Runs tutorial specified by number with chosen bot
  update      Updates Lia development tools
  verify      Verifies if the content in bot-dir is valid
  zip         Verifies, compiles and zips the bot in botDir

Flags:
  -h, --help        help for lia
  -l, --languages   show all supported languages
  -v, --version     show tools version

Use "lia [command] --help" for more information about a command.
```

## Build
```go
go build -o build/lia
```
[How to build go programs.](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04)

### Dependencies ###
```bash
go get github.com/spf13/cobra
go get github.com/palantir/stacktrace
go get github.com/mholt/archiver
go get github.com/fatih/color
```