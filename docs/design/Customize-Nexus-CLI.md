# Overview

This page lists ways one can customize the Nexus CLI behaviour

### Nexus CLI verbose output
You can enable verbose output by adding the `--debug` flag to any nexus command
```shell
nexus app init --name <app-name> --debug
```

### Nexus version and Upgrading your Nexus CLI
Find out what Nexus version you're running
```shell
nexus version
```

Besides printing the Nexus CLI version, it also provides versions of the other Nexus components the CLI uses. For instance:

```shell
~$ nexus version
A new version of Nexus CLI (0.0.24) is available
NexusCli: v0.0.23
NexusCompiler: v0.0.7
NexusAppTemplates: v0.0.4
NexusDatamodelTemplates: v0.0.4
```
In addition, if there's a new version of the CLI, you'll be prompted to upgrade. You can also upgrade your CLI to a specific version.

```shell
nexus upgrade cli --version 0.0.24 # if no 'version' is provided, the CLI will be upgraded to the latest available version 
```

Finally, if you do not want to be always prompted to upgrade Nexus CLI, you can disable prompts:
```shell
nexus config set --disable-upgrade-prompt
```
This will stop the Nexus CLI from prompting you to upgrade, but it will still notify you if a newer version is available.

### Nexus config
Nexus provides a mechanism to set and view user preferences. Nexus stores these in the user's home directory at `$HOME/.nexus/config`

You can view the config using:
```shell
nexus config view
```

You can set preferences using the `nexus config set <flag>` command. We currently have a limited set of flags available. We'll add more as we need them and based on the feedback we receive.

Finally, there's a flag that will cause the Nexus CLI to print verbose output for every command without explicitly having to specify the `--debug` option
```shell
nexus config set --debug-always
```

### Check Nexus pre-requisites
You can check if your system satisfies all the Nexus CLI prerequisites by:
```shell
nexus prereq list # only lists the prerequisites
nexus prereq verify # verifies prerequisites
```

#### Prerequisites for each Nexus command
It is also possible to check if the prerequisites for that specific command are satisfied and if not, which ones are not. For instance, you can list the prerequisites for `nexus datamodel build` are met by doing: 
```shell
nexus datamodel build --list-prereq
```
The `--list-prereq` flag is enabled for all functional Nexus commands. The command will exit after listing down the prerequisites.

Before every command runs, that specific command's prerequisites are validated. For instance, `nexus datamodel build` requires docker to be running. If docker weren't running, you would encounter an error:
```shell
$ nexus datamodel build --name vmware
Error: verify if docker is running failed with error exit status 1
Usage:
  nexus datamodel build [flags]

Flags:
  -h, --help          help for build
  -n, --name string   name of the datamodel to be build

Global Flags:
      --debug               Enables extra logging
      --list-prereq         List prerequisites
      --skip-prereq-check   Skip prerequisites check
```

If, however, you still want to proceed with running the command, you can use the `--skip-prereq-check` flag. 
