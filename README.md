# bk

A simple command line interface to the [buildkite](http://buildkite.com) service.

![ScreenShot](/docs/buildkite-cli-builds.gif)

# usage

```
usage: bk <command> [<flags>] [<args> ...]

A command-line interface for buildkite.com.

Flags:
  --help  Show help.

Commands:
  help [<command>]
    Show help for a command.

  projects
    List projects under an orginization.

  builds
    List latest builds for the current project.

  open
    Open buildkite project in your browser for the current project.

  setup
    Configure the buildkite cli with a new token.

```

# license

Released under MIT license.