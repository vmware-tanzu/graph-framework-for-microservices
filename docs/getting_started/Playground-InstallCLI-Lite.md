# Install Nexus CLI
[[Prev]](Playground-Lite.md) [[Exit]](../../README.md)  [[Next]](Playground-InstallRuntime-Lite.md)

![InstallCLI](../images/Playground-2-install-cli.png)

### 1. Clone Nexus Repo
```
git clone -b api-gw-tenancy-plugin git@github.com:xmen4xp/graph-framework-for-microservices.git
cd graph-framework-for-microservices/
export NEXUS_REPO_DIR=${PWD}
```

### 2. Build & Install Nexus CLI
#### For Linux
```
make cli.build.linux
```
For MacOS
```
make cli.build.darwin
```

### 3. Install Nexus CLI
#### For Linux
```
sudo make cli.install.linux
```
For MacOS
```
sudo make cli.install.darwin
```

[[Prev]](Playground-Lite.md) [[Exit]](../../README.md)  [[Next]](Playground-InstallRuntime-Lite.md)
