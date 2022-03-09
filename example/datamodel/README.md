# Example graph to parse

```
                  Root
                   |
                   |
                  Config
                 /     \
                /       \
               /         \
            Gns         AccessControlPolicy
             |                  |
             |                  |
         SvcGroup <----------> ACPConfig
                   softlink

```

DSL example based on https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=NSBU&title=Nexus+Platform#NexusPlatform-TL;DR;

TODO:
Extend examples with:
1. Link types
- Link1 CustomLinkType `nexus:"link"`
- NamedLink1 syncMap[string]CustomNamedLinkType `nexus:"link"`
2. Status type
- status CustomStatusType `nexus:"status"`
3. Custom comments
- rest-api
- version
- validation-endpoints
