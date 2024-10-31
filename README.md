# Conjur Inspect

Conjur Inspect is a command-line tool that validates a system for running
Conjur Enterprise successfully and gathers system information for
troubleshooting when Conjur Enterprise isn't running correctly.

## Certification level

![](https://img.shields.io/badge/Certification%20Level-Trusted-28A745?link=https://github.com/cyberark/community/blob/master/Conjur/conventions/certification-levels.md)

This repo is a **Trusted** level project. It is supported by CyberArk and has
been verified to work with Conjur Enterprise. For more detailed information on
our certification levels, see
[our community guidelines](https://github.com/cyberark/community/blob/master/Conjur/conventions/certification-levels.md#trusted).

## Getting Started

### Obtaining install package

The latest Conjur Inspect install package may either be downloaded from
[Releases](https://github.com/cyberark/conjur-inspect/releases) or copied
from a Conjur Enterprise Appliance container image (version 13.5 or later)
with the following commands:

1. First, create a container instance using the Conjur Enterprise Appliance
   image. This doesn't need to be a configured Conjur instance, only a running
   container using the image:

   ```sh
   docker run -d --name conjur registry.tld/conjur-appliance:<version>
   ```

2. List the available Conjur Inspect install packages with the command:

   ```sh
   docker exec -it conjur ls -la /opt/conjur/conjur-inspect
   ```

   This will output a file listing similar to:

   ```sh
   -rw-r--r-- 1 root   root   32989 Oct 22 14:19 conjur-inspect_0.4.0_linux_386.tar.gz
   -rw-r--r-- 1 root   root   32991 Oct 22 14:19 conjur-inspect_0.4.0_linux_amd64.tar.gz
   -rw-r--r-- 1 root   root   32991 Oct 22 14:19 conjur-inspect_0.4.0_linux_arm64.tar.gz
   ```

3. Select the correct file based on your host system architecture and copy it
   to the host with the command:

   ```sh
   docker cp conjur:/opt/conjur/conjur-inspect/conjur-inspect_0.4.0_linux_amd64.tar.gz .
   ```

4. Remove the temporary Conjur container with the command:

   ```sh
   docker rm -f conjur
   ```

### Installing

To install Conjur Inspect, copy the install package for your system
architecture to your target host machine. This is the machine or VM where the
Conjur Enterprise container will run.

Extract the tool from the gzip archive:

```sh-session
$ tar -xvf conjur-inspect-0.4.0_amd64.tgz
```

Or install it from one of the system packages:

```sh
yum install conjur-inspect_0.4.0_amd64.rpm
```

```sh
apt install conjur-inspect_0.4.0_386.deb
```

### Running

Run Conjur Inspect with the following command as the same system user that will
run the Conjur container:

```sh-session
$ conjur-inspect
```

This command generates an inspection report of different system indicators
related to Conjur stability performance:

```sh
========================================
Conjur Enterprise Inspection Report
Version: 0.3.0-ef572db (Build 46)
========================================

CPU
---
INFO - CPU Cores: 4
INFO - CPU Architecture: amd64

Disk
----
INFO - Disk Space (tmpfs, /dev): 8.1 GB Total, 0 B Used ( 0%), 8.1 GB Free
...

Memory
------
INFO - Memory Total: 16 GB
INFO - Memory Free: 15 GB
INFO - Memory Used: 813 MB (5.0 %)

Host
----
INFO - Hostname: ip-172-31-17-12.ec2.internal
INFO - Uptime: 8 minutes 35 seconds
INFO - OS: linux, redhat, rhel, 8.6
INFO - Virtualization: None

Follower
--------
ERROR - Leader Hostname: N/A (Leader hostname is not set. Set the 'MASTER_HOSTNAME' environment variable to run this check)

Container Runtime
-----------------
ERROR - Docker: N/A (failed to inspect Docker runtime: exec: "docker": executable file not found in $PATH ())
INFO - Podman Version: 4.2.0
INFO - Podman Driver: overlay
INFO - Podman Graph Root: /home/ec2-user/.local/share/containers/storage
INFO - Podman Run Root: /run/user/1000/containers
INFO - Podman Volume Path: /home/ec2-user/.local/share/containers/storage/volumes

Ulimits
-------
INFO - core file size (blocks, -c): 0
...
```

Available options and flags may be view by running:

```sh
conjur-inspect --help
```

## Gathering container specific inspection data

To gather inspection data for a Conjur Enterprise container that is already
running, include the `--container-id` argument with either the container name
or ID. For example:

```sh
conjur-inspect --container-id conjur
```

## Raw data report

In addition to the output report, `conjur-inspect` records the raw inspection
data in a gzipped archived saved as a filed to the working directory. By default,
this archived is named with the timestamp when `conjur-inspect` was
run. For example, `2023-05-26-15-51.tar.gz`.

The name of the output archive may be customized with the `--report-id` argument.
For example:

```sh
conjur-inspect --report-id standby
```

This results in an output archived named `standby.tar.gz`.

## Inspecting disk performance

The Conjur Inspect disk performance checks require an additional dependency,
`fio`, to be installed as a prerequisite. This may be installed from OS package
management tools. For example:

```sh
yum install fio
```

or

```sh
apt install fio
```

When `fio` is installed, `conjur-inspect` reports additional information about
disk latency and throughput:

```sh-session
./conjur-inspect
========================================
Conjur Enterprise Inspection Report
Version: 0.3.0-ef572db (Build 46)
========================================

...

Disk
----
...
INFO - FIO - Read IOPs (/home/ec2-user/conjur-inspect_0.3.0_linux_amd64): 1624.33 (Min: 874, Max: 3902, StdDev: 221.16)
INFO - FIO - Write IOPs (/home/ec2-user/conjur-inspect_0.3.0_linux_amd64): 1674.32 (Min: 884, Max: 4018, StdDev: 228.43)
INFO - FIO - Read Latency (99%, /home/ec2-user/conjur-inspect_0.3.0_linux_amd64): 0.00 ms
INFO - FIO - Write Latency (99%, /home/ec2-user/conjur-inspect_0.3.0_linux_amd64): 0.00 ms
INFO - FIO - Sync Latency (99%, /home/ec2-user/conjur-inspect_0.3.0_linux_amd64): 1.48 ms

...
```

## Troubleshooting

To troubleshoot issues running `conjur-inspect`, add the `--debug` flag for
detailed execution output.

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
