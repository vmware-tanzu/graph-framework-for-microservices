package parser

import (
	"io/fs"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// ParseDSLNodes walks recursively through given path and looks for structs types definitions to add them to graph
func ParseGraphQLFiles(startPath string) map[string]string {
	graphqlFiles := make(map[string]string, 0)
	err := filepath.Walk(startPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == "build" {
				log.Infof("Ignoring build directory...")
				return filepath.SkipDir
			}
			if info.Name() == "vendor" {
				log.Infof("Ignoring vendor directory...")
				return filepath.SkipDir
			}
		} else {
			if filepath.Ext(path) == ".graphql" {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				graphqlFiles[path] = string(data)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to ParseDSLNodes %v", err)
	}

	return graphqlFiles
}
