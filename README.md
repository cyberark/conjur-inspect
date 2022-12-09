# Conjur Preflight

Conjur Preflight is a CLI qualification tool for environments running Conjur
Enterprise.

For the historic hackathon preflight project, see the report as the prior
commit [here](https://github.com/conjurinc/conjur-preflight/tree/451c9378f7df89659c2e9d05da1ea0e2da3c5269).

## Getting Started

To run the Conjur Preflight tool, download the latest version for your system
architecture from Releases to your target host machine. This is the machine
or VM where the Conjur Enterprise container will run.

Extract the tool, enable it to execute, and run it with:
```sh-session
$ tar -xvf conjur-preflight-1.0.0_amd64.tgz
$ chmod +x conjur-preflight-1.0.0_amd65
$ ./conjur-preflight-1.0.0_amd65
...
PASS / WARN / FAIL / ERROR
```

## Development

### Running

```
make run
```

### Building

```
make build
```

### Running unit tests

```
make test
```

### Adding new checks

==TODO==

## Contributing

For information on how to contribute to this project, see [CONTRIBUTING.md](./CONTRIBUTING.md).
