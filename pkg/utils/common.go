package utils

import (
	"api-gw/pkg/config"
	"embed"
	"io/fs"
	"net/http"
	"os"
)

func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsServerConfigValid(conf *config.Config) bool {
	if conf != nil {
		if conf.Server.Address != "" && conf.Server.CertPath != "" && conf.Server.KeyPath != "" {
			return true
		}
	}
	return false
}

func GetHttpFS(embedFs embed.FS, name string) (http.FileSystem, error) {
	fsys, err := fs.Sub(embedFs, name)
	if err != nil {
		return nil, err
	}

	return http.FS(fsys), nil
}
