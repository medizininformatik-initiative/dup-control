# Polar Control

Utility to manage [POLAR][polar] workpackage and -pilot algorithm
execution using docker.

## Installation

### Download the Latest Version

* [Windows][windows-amd64]
* [Linux][linux-amd64] ([arm][linux-arm64])
* [MacOS][darwin-amd64] ([arm][darwin-arm64])

### Windows

Move the downloaded `polarctl.exe` into a directory in your PATH (`echo $env:PATH` (PowerShell)). 

* using `C:\Windows\System32\` will enable polarctl for all users. *Note: Access to C:\Windows\System32\ may require administrator access privileges.*
* using `C:\Users\%USERNAME%\AppData\Local\Microsoft\WindowsApps` will enable polarctl for the current user. *Note: Make sure the directory is in your PATH!*

polarctl can then be executed via cmd or powershell. 

### Linux / macOS

Move the downloaded `polarctl` binary into a directory in your PATH (`echo $PATH`).

```shell
sudo mv polarctl /usr/local/bin/polarctl
```

## Usage

### Working Directory

Choose a working directory from where you will execute polarctl commands. This is very important, as polarctl uses the current 
working directory (cwd) to store results from the executed workpackages and to find the polarctl config.  

### Create polarctl Config

Create a file called `config.toml` within your chosen polar working directory. The minimal configuration contains the polar 
container registry credentials (formerly used with `docker login` command) in the following form:
```
registryUser = "polar-dic-<site>"
registryPass = "<password>"
```

### Global Settings

Some settings can be set via CLI flag or config file. The table below lists the flags and corresponding keys for
the config file. *CLI opts will override config settings.*

| CLI Flag               | Config Key          | Description                                                          | Optional? | Default |
|------------------------|---------------------|----------------------------------------------------------------------|-----------|---------|
| --config               |                     | Specify a config file rather than using the default config path      | Yes       | config.toml |
| --disable-update-check | disableUpdateCheck  | Disable upgrade check on startup                                     | Yes       | false  |

### Retrieval

```shell
polarctl retrieve --wp <workpackage> --fhir-server-endpoint "https://some-fhir-server" [flags] 
```

#### Settings

Some settings can be set via CLI flag or config file. The table below lists the flags and corresponding keys for 
the config file. *CLI opts will override config settings.*

| CLI Flag               | Config Key                   | Description                                                          | Optional? | Default |
|------------------------|------------------------------|----------------------------------------------------------------------|-----------|---------|
| --wp                   |                              | Workpackage algorithm to execute, e.g. 'wp-1-1-pilot'                | No        |     | 
| --site                 | retrieve.site                | Determines which image to use, as images are (not necessarily) hand-tailored for different dic sites. (e.g. 'dic-giessen', 'dic-leipzig', 'dic-muenchen'). | Yes        | latest |
| --fhir-server-endpoint | retrieve.fhirServerEndpoint  | URL including base path of the FHIR Server to be queried, e.g.: 'https://example.com/r4/' | No        |     |
| --fhir-server-user     | retrieve.fhirServerUser      | Username for basic auth protected communication with FHIR Server     | Yes       |         |
| --fhir-server-pass     | retrieve.fhirServerPass      | Password for basic auth protected communication with FHIR Server     | Yes       |         |
| --fhir-server-cacert   | retrieve.fhirServerCACert    | CA Certificate file[^cafile] for https connection to FHIR Server     | Yes       |         |
| --dev                  |                              | Enables settings for local development                               | Yes       | false   |

#### Example

```shell
polarctl retrieve --wp wp-1-1-pilot --fhir-server-endpoint "https://mii-agiop-3p.life.uni-leipzig.de/fhir/" --site "dic-giessen"
```

### Analysis

```shell
polarctl analyze --wp <workpackage> [flags] 
```

#### Settings

Some settings can be set via CLI flag or config file. The table below lists the flags and corresponding keys for
the config file. *CLI opts will override config settings.*

| CLI Flag               | Config Key          | Description                                                          | Optional? | Default |
|------------------------|---------------------|----------------------------------------------------------------------|-----------|---------|
| --wp                   |                     | Workpackage algorithm to execute, e.g. 'wp-1-1-pilot'                | No        |        | 
| --version              | analyze.version     | Determines which version of the analysis algorithm to use            | Yes       | latest |
| --dev                  |                     | Enables settings for local development                               | Yes       | false  |

#### Example

```shell
polarctl analyze --wp wp-1-1-pilot --version "1.0"
```


[polar]: https://www.medizininformatik-initiative.de/de/POLAR

[windows-amd64]: https://git.smith.care/smith/uc-phep/polar/polar-control-2/-/jobs/artifacts/main/raw/builds/polarctl-windows-amd64.exe?job=build-branch
[linux-amd64]: https://git.smith.care/smith/uc-phep/polar/polar-control-2/-/jobs/artifacts/main/raw/builds/polarctl-linux-amd64?job=build-branch
[linux-arm64]: https://git.smith.care/smith/uc-phep/polar/polar-control-2/-/jobs/artifacts/main/raw/builds/polarctl-linux-arm64?job=build-branch
[darwin-amd64]: https://git.smith.care/smith/uc-phep/polar/polar-control-2/-/jobs/artifacts/main/raw/builds/polarctl-darwin-amd64?job=build-branch
[darwin-arm64]: https://git.smith.care/smith/uc-phep/polar/polar-control-2/-/jobs/artifacts/main/raw/builds/polarctl-darwin-arm64?job=build-branch
