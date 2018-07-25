# Lia CLI
```bash
lia-cli is a CLI tool for easier development of Lia bots.

Usage:
  lia-cli [flags]
  lia-cli [command]

Available Commands:
  bot         Create new bot
  compile     Compiles/prepares bot in specified dir
  fetch       Fetches a bot from url and sets a new name
  generate    Generates a game
  help        Help about any command
  play        Compiles and generates a game between bot1 and bot2
  replay      Runs a replay viewer
  tutorial    Runs tutorial specified by number with chosen bot
  update      Updates Lia development tools
  verify      Verifies if the content in bot-dir is valid
  zip         Verifies, compiles and zips the bot in botDir

Flags:
  -h, --help        help for lia-cli
  -l, --languages   Show all supported languages
  -v, --version     Show tools version

Use "lia-cli [command] --help" for more information about a command.
```

## Build
```go
go build -o build/lia
```
