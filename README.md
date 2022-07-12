# DUP Control

Execution helper for various Data Usage Projects (DUP) that are implemented as run-once container images and (may be) 
split into retrieval and analysis parts.

## Installation

### Prerequisites

#### Docker

* Docker must be installed on the system. Please follow the [official installation instructions][docker-install].
* Also consider allowing your non-root Linux user to use docker (and dupctl) by adding it to the docker group 
(see [docker docs][docker-ugroup]) otherwise only a root user will be able to execute dupctl. 

### Download the Latest Version

| Operating System | AMD64 Processor Architecture ([?][wiki-amd64]) | ARM64 Processor Architecture ([?][wiki-arm64])  |
|------------------|------------------------------------------------|-------------------------------------------------|
| Windows          | [Download][windows-amd64]                      | -                                               |
| Linux            | [Download][linux-amd64]                        | [Download][linux-arm64]                         |
| MacOS            | [Download][darwin-amd64]                       | [Download][darwin-arm64]                        |

### Windows

Move the downloaded `dupctl.exe` into a directory in your [PATH][wiki-path] (`echo $env:PATH` (PowerShell)). 

* using `C:\Windows\System32\` will enable dupctl for all users. *Note: Access to C:\Windows\System32\ may require administrator access privileges.*
* using `C:\Users\%USERNAME%\AppData\Local\Microsoft\WindowsApps` will enable dupctl for the current user. *Note: Make sure the directory is in your [PATH][wiki-path] (see `echo $env:PATH` (PowerShell))!*

dupctl can then be executed via cmd or powershell. 

**Download using Command Line (Windows)**

You can also download the executable using a command line:

```shell
curl https://dupctl.s3.amazonaws.com/dupctl-windows-amd64.exe -O dupctl.exe
```

### Linux / macOS

Move the downloaded `dupctl` binary into a directory in your [PATH][wiki-path] (`echo $PATH`).

```shell
sudo mv dupctl /usr/local/bin/dupctl
sudo chmod +x /usr/local/bin/dupctl
```

**Download using Command Line (Linux/macOS)**

You can also download the executable using a command line, fill in the appropriate link from the [download table](#download-the-latest-version).

```shell
curl [link] -O dupctl
```

## Usage

### Working Directory

Choose a working directory from where you will execute dupctl commands. This is very important, as dupctl uses the current 
working directory (cwd) to store results from the executed dups and to find the dupctl config.  

### Create dupctl Config

Create a file called `config.toml` within your chosen working directory. The minimal configuration contains the projects 
container registry, the registry credentials (formerly used with `docker login` command) and a project name in the following form:
```toml
project = "some-project"
registry = "example.com/container-registry"
registryUser = "<username>"
registryPass = "<password>"
```

*Note: The project name will distinguish docker containers belonging to different projects.*

*Note: Container Registry Credentials are provided per DIC by the UC PheP Development Team. Please contact
[Jonas Wagner](mailto:jwagner@life.uni-leipzig.de) or [Frank Meineke](mailto:Frank.Meineke@imise.uni-leipzig.de).*

### Global Settings

Some settings can be set via CLI flag or config file. The table below lists the flags and corresponding keys for
the config file. *CLI opts will override config settings.*

| CLI Flag               | Config Key         | Description                                                      | Optional? | Default     |
|------------------------|--------------------|------------------------------------------------------------------|-----------|-------------|
| --config               |                    | Specify a config file rather than using the default config path  | Yes       | config.toml |
| --disable-update-check | disableUpdateCheck | Disable upgrade check on startup                                 | Yes       | false       |
| --offline              | offline            | Assumes an air-gapped environment (No Update Check / Image Pull) | Yes       | false       |

### Retrieval

```shell
dupctl retrieve --dup <dup-name> --fhir-server-endpoint "https://some-fhir-server" [flags] 
```

#### Settings

Some settings can be set via CLI flag or config file. The table below lists the flags and corresponding keys for 
the config file. *CLI opts will override config settings.*

| CLI Flag               | Config Key                  | Description                                                                                                                                                    | Optional? | Default |
|------------------------|-----------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|---------|
| --dup                  |                             | DUP to execute, e.g. 'vhf'                                                                                                                                     | No        |         | 
| --version              | retrieve.version            | Determines which image to use, as images can be versioned or hand-tailored for different dic sites. (e.g. '0.1', dic-giessen', 'dic-leipzig', 'dic-muenchen'). | Yes       | latest  |
| --fhir-server-endpoint | retrieve.fhirServerEndpoint | URL including base path of the FHIR Server to be queried, e.g.: 'https://example.com/r4/'                                                                      | No        |         |
| --fhir-server-user     | retrieve.fhirServerUser     | Username for basic auth protected communication with FHIR Server                                                                                               | Yes       |         |
| --fhir-server-pass     | retrieve.fhirServerPass     | Password for basic auth protected communication with FHIR Server                                                                                               | Yes       |         |
| --fhir-server-cacert   | retrieve.fhirServerCACert   | CA Certificate file[^cafile] for https connection to FHIR Server                                                                                               | Yes       |         |
| --fhir-server-token    | retrieve.fhirServerToken    | Token for token based auth protected communication with FHIR Server                                                                                            | Yes       |         |
| --env / -e             | retrieve.env                | Passes environment variables to the dup scripts, e.g.: -e "MAX_BUNDLES=5"                                                                                      | Yes       |         |

#### Example

```shell
dupctl retrieve --dup vhf --fhir-server-endpoint "https://mii-agiop-3p.life.uni-leipzig.de/fhir/"
```

### Analysis

```shell
dupctl analyze --dup <dup-name> [flags] 
```

#### Settings

Some settings can be set via CLI flag or config file. The table below lists the flags and corresponding keys for
the config file. *CLI opts will override config settings.*

| CLI Flag   | Config Key      | Description                                                               | Optional? | Default |
|------------|-----------------|---------------------------------------------------------------------------|-----------|---------|
| --dup      |                 | DUP to execute, e.g. 'vhf'                                                | No        |         | 
| --version  | analyze.version | Determines which version of the analysis algorithm to use                 | Yes       | latest  |
| --env / -e | analyze.env     | Passes environment variables to the dup scripts, e.g.: -e "MAX_BUNDLES=5" | Yes       |         |

#### Example

```shell
dupctl analyze --dup vhf --version "1.0"
```

### Example Configuration

```toml
registryUser = "some-dic"
registryPass = "some-individual-password"

[retrieve]
fhirServerEndpoint = "https://example.com/fhir"
fhirServerUser = "some-fhir-server-user"
fhirServerPass = "some-fhir-server-pass"
env = {"MAX_BUNDLES" = "100", "COUNT" = 200}
```

## Troubleshooting

### Permission denied

Getting a `permission denied` error when using `dupctl upgrade` usually means you require access rights on the 
currently installed dupctl file. On linux using `sudo dupctl upgrade` should suffice.

[docker-install]: https://docs.docker.com/get-docker/
[docker-ugroup]: https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user

[wiki-amd64]: https://en.wikipedia.org/wiki/X86-64#AMD64
[wiki-arm64]: https://de.wikipedia.org/wiki/Arm-Architektur#Armv8-A_(2011)
[wiki-path]: https://en.wikipedia.org/wiki/PATH_(variable)

[windows-amd64]: https://dupctl.s3.amazonaws.com/dupctl-windows-amd64.exe
[linux-amd64]: https://dupctl.s3.amazonaws.com/dupctl-linux-amd64
[linux-arm64]: https://dupctl.s3.amazonaws.com/dupctl-linux-arm64
[darwin-amd64]: https://dupctl.s3.amazonaws.com/dupctl-darwin-amd64
[darwin-arm64]: https://dupctl.s3.amazonaws.com/dupctl-darwin-arm64
