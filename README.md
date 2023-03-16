# Conjur Inspect

Conjur Inspect is a command-line tool that validates a system for running
Conjur Enterprise successfully and gathers system information for
troubleshooting when Conjur Enterprise isn't running correctly.

Conjur Inspect is part of the CyberArk Conjur
[Open Source Suite](https://cyberark.github.io/conjur/) of tools.

## Certification level

![](https://img.shields.io/badge/Certification%20Level-Trusted-28A745?link=https://github.com/cyberark/community/blob/master/Conjur/conventions/certification-levels.md)

This repo is a **Trusted** level project. It is supported by CyberArk and has
been verified to work with Conjur Open Source. For more detailed information on
our certification levels, see
[our community guidelines](https://github.com/cyberark/community/blob/master/Conjur/conventions/certification-levels.md#trusted).

## Getting Started

To run the Conjur Inspect tool, download the latest version for your system
architecture from [Releases](https://github.com/cyberark/conjur-inspect/releases)
to your target host machine. This is the machine or VM where the Conjur
Enterprise container will run.

Extract the tool, enable it to execute, and run it with:
```sh-session
$ tar -xvf conjur-inspect-1.0.0_amd64.tgz
$ chmod +x conjur-inspect-1.0.0_amd64
$ ./conjur-inspect-1.0.0_amd64
...
PASS / WARN / FAIL / ERROR
```

## Community Support

Our primary channel for support is through our CyberArk Commons community
[here](https://discuss.cyberarkcommons.org/c/conjur/5).

## Code Maintainers

CyberArk Conjur Team

## Contributing

For information on how to contribute to this project, see [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

Copyright (c) 2023 CyberArk Software Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this software except in compliance with the License. You may obtain a copy of
the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.

For the full license text see [LICENSE](./LICENSE).
