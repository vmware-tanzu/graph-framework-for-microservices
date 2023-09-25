# Install Sock Shop data model on Nexus Runtime

[[Prev]](Playground-SockShop-Compile-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Access-Datamodel-API-Lite.md)


## Export KUBECONFIG to Nexus Runtime
```
export HOST_KUBECONFIG=$NEXUS_REPO_DIR/nexus-runtime-manifests/k0s/.kubeconfig
```

## Install data model
```
make dm.install
```

[[Prev]](Playground-SockShop-Compile-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Access-Datamodel-API-Lite.md)
