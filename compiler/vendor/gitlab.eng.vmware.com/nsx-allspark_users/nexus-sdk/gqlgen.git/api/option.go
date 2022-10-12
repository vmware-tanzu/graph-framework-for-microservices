package api

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/codegen/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/plugin"
)

type Option func(cfg *config.Config, plugins *[]plugin.Plugin)

func NoPlugins() Option {
	return func(cfg *config.Config, plugins *[]plugin.Plugin) {
		*plugins = nil
	}
}

func AddPlugin(p plugin.Plugin) Option {
	return func(cfg *config.Config, plugins *[]plugin.Plugin) {
		*plugins = append(*plugins, p)
	}
}

// PrependPlugin prepends plugin any existing plugins
func PrependPlugin(p plugin.Plugin) Option {
	return func(cfg *config.Config, plugins *[]plugin.Plugin) {
		*plugins = append([]plugin.Plugin{p}, *plugins...)
	}
}

// ReplacePlugin replaces any existing plugin with a matching plugin name
func ReplacePlugin(p plugin.Plugin) Option {
	return func(cfg *config.Config, plugins *[]plugin.Plugin) {
		if plugins != nil {
			found := false
			ps := *plugins
			for i, o := range ps {
				if p.Name() == o.Name() {
					ps[i] = p
					found = true
				}
			}
			if !found {
				ps = append(ps, p)
			}
			*plugins = ps
		}
	}
}
