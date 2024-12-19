![Version](https://img.shields.io/github/v/release/captainhook-go/captainhook?style=flat-square&label=version)
[![License](https://poser.pugx.org/captainhook/captainhook/license.svg?v=1)](https://packagist.org/packages/captainhook/captainhook)
[![Go Report Card](https://goreportcard.com/badge/github.com/captainhook-go/captainhook?style=flat-square)](https://goreportcard.com/report/github.com/captainhook-go/captainhook)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.21-61CFDD.svg?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/captainhook-go/captainhook.svg?style=flat-square)](https://pkg.go.dev/github.com/captainhook-go/captainhook)
[![Mastodon](https://img.shields.io/badge/Mastodon-%40captainhook-purple.svg)](https://phpc.social/@captainhook)


# CaptainHook

<img src="https://captainhook-go.github.io/captainhook/gfx/ch.png" alt="CaptainHook logo" align="right" width="200"/>

*CaptainHook* is an easy to use and very flexible git hook manager for software developers.
It enables you to configure your git hook actions in a simple json file and easily share them within your team.

You can use *CaptainHook* to validate or prepare your commit messages, ensure code quality
or run unit tests before you commit or push changes to git. You can automatically clear
local caches or install the correct dependencies after pulling the latest changes.

You can run your own commands or use loads of built-in functionality.
For more information have a look at the [documentation](https://captainhook-go.github.io/captainhook/ "CaptainHook Documentation").

## Installation

You can download the application binary for your platform from the [release page](https://github.com/captainhook-go/captainhook/releases/latest "Latest CaptainHook Release").
Or use one of the following options.

Use `Homebrew` to install *CaptainHook*.
```bash
brew tap captainhook-go/captainhook
brew install captainhook
```
Use `go install` to install *CaptainHook*.
```bash
go install github.com/captainhook-go/captainhook/cmd/captainhook@latest
```


## Setup

After installing CaptainHook, navigate to your project directory and use the *captainhook* init command to create a configuration file.
```bash
cd my-project-repo
captainhook init
```

As soon as you have a configuration file the only thing left is to activate the hooks by installing them to
your local `.git/hooks` directory. To do so just run the following *captainhook* command.
```bash
captainhook install
```

## Configuration

Here's an example *captainhook.json* configuration file.
```json
{
  "hooks": {
    "commit-msg": {
      "actions": [
        {
          "run": "CaptainHook::Message.MustFollowBeamsRules"
        }
      ]
    },
    "pre-commit": {
      "actions": [
        {
          "run": "unittest"
        }
       ]
    },
    "pre-push": {
      "actions": [
        {
          "run": "CaptainHook::Branch.PreventPushOfFixupAndSquashCommits",
          "options": {
            "branches-to-protect": ["main", "integration"]
          }
        }
      ]
    }
  }
}
```

## Contributing

So you'd like to contribute to `CaptainHook`? Excellent! Thank you very much.
I can absolutely use your help.

Have a look at the [contribution guidelines](CONTRIBUTING.md).