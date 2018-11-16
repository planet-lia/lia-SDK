# Lia-SDK

Download latest release [here](https://github.com/liagame/lia-SDK/releases).


## Lia CLI
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
  playground  Runs playground specified by number with chosen bot
  update      Updates Lia development tools
  verify      Verifies if the content in bot-dir is valid
  zip         Verifies, compiles and zips the bot in botDir

Flags:
  -h, --help        help for lia
  -l, --languages   show all supported languages
  -v, --version     show tools version

Use "lia [command] --help" for more information about a command.
```


## Development
### Dependencies
Install all dependencies and update them by cd-ing into the repository
root and running:
```bash
go get -u ./...
```

### Build
```bash
chmod +x scrpits/build.sh
./scripts/build.sh
```

## Analytics 
To change your status to a tester in google analytics,
find the .lia.json file in your home directory. There you should add a line
so it looks someting like this.

```json
{
  "analyticsallow": true,
  "analyticsallowedversion": "sdk-version",
  "trackingid": "your-tracking-id",
  "testing": true
}
```
