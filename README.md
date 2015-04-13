# bk

A simple command line interface to the [buildkite](http://buildkite.com) service.

![ScreenShot](/docs/buildkite-cli-builds.gif)

# install

At the moment installation is done just using the go get command.

```
go get github.com/wolfeidau/buildkite-cli/cmd/bk
```

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
    Open builds list in your browser for the current project.

  setup
    Configure the buildkite cli with a new token.

```

# buildkite API notes

* Would be nice to have access to the orginization slug in the project(s) results as it is critical for building URLs relating to associated content.
* Would be nice to have notable website URLs in the project/build content, or some base resource which defines then in the REST API. The rationale for this is if you wanted to use bk for a test or onsite instance of buildkite the sources would currently need to be modified.
* The API is a little slow for an interactive application, can take up to 9 seconds to pull down the builds for a project as I am not currently caching these results. At the moment I do a call to orgs, then projects, then builds to retrieve the details.

# license

Released under MIT license.