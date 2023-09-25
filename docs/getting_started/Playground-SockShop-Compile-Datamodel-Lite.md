# Compile the data model

[[Prev]](Playground-SockShop-Complete-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Install-Datamodel-Lite.md)

![SockShop](../images/Playground-8-Compile-Datamodel.png)

Let's build our data model and generate all the artifacts needed for the runtime.

## Compile data model

Nexus compiler can be invoked to build the datamodel with the following command:

```
COMPILER_TAG=letsplay make datamodel_build
```

This will generate all the artifacts needed at install and runtime.

The generated artifacts are available in the $PWD/build directory.

[[Prev]](Playground-SockShop-Complete-Datamodel-Lite.md) [[Exit]](../../README.md) [[Next]](Playground-SockShop-Install-Datamodel-Lite.md)