# Contributing

For general contribution and community guidelines, please see the [community repo](https://github.com/cyberark/community).

1. [Fork the project](https://help.github.com/en/github/getting-started-with-github/fork-a-repo)
1. [Clone your fork](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository)
1. Make local changes to your fork by editing files
1. [Commit your changes](https://help.github.com/en/github/managing-files-in-a-repository/adding-a-file-to-a-repository-using-the-command-line)
1. [Push your local changes to the remote server](https://help.github.com/en/github/using-git/pushing-commits-to-a-remote-repository)
1. [Create new Pull Request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request-from-a-fork)

From here your pull request will be reviewed and once you've responded to all
feedback it will be merged into the project. Congratulations, you're a
contributor!

## Running development version locally

To run a development version of Conjur Inspect, start the development
environment with:

```sh-session
$ cd dev
$ ./start
...
[root@a0a4483ca6c9 conjur-inspect]#
```

This starts a development container and begins a terminal within. Then run
the CLI from source with the command:

```sh
make run
```

### Building the Conjur Inspect CLI

To build the CLI and its release artifacts, run the command:

```sh
bin/build-release
```

When this command finishes, the CLI binary and installers are available under
the 'dist/` directory.

### Running unit tests

Run the Conjur Inspect unit tests with the command:

```sh
bin/test-unit
```

Alternatively, within the `dev/` containerized development environment, the unit
tests may be run with the command:

```sh
make test
```

### Running integration tests

To run a set of integration tests that exercise `conjur-inspect` in a couple of
representative environments, use the command:

```sh
bin/test-integration
```

## Releasing

To create a new release, follow the instructions in our general release
guidelines [here](https://github.com/cyberark/community/blob/master/Conjur/CONTRIBUTING.md#release-process).
