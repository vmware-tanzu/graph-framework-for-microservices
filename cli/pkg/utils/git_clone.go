package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	net_http "net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func GoModInit(path string, current bool) error {
	if path != "" {
		fmt.Printf("Intializing gomodule\nGo mod init name: %s\n", path)
		cmd := exec.Command("go", "mod", "init", path)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOTOOLCHAIN=go1.18")
		out, err := cmd.Output()
		fmt.Printf("output: %s", out)
		if err != nil {
			return err
		}
		if current == false {
			err := os.Chdir("..")
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("Intializing gomodule\n")
		cmd := exec.Command("go", "mod", "init")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOTOOLCHAIN=go1.18")
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetModuleName(path string) (string, error) {
	fmt.Println("getting the current modulename")
	if path != "" {
		err := os.Chdir(path)
		if err != nil {
			return "", err
		}
	}
	cmd := exec.Command("go", "list", "-m")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if path != "" {
		err := os.Chdir("..")
		if err != nil {
			return "", err
		}
	}
	return string(stdout), nil
}

func RenderFile(filename string, params interface{}) error {
	fmt.Printf("calling renderfile... %s\n", filename)
	var b bytes.Buffer
	if strings.HasSuffix(filename, "tmpl") {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("error in accessing file %s\n", filename)
			return err
		}
		tmpl, err := template.New("template").Parse(string(data))
		if err != nil {
			fmt.Printf("could not get template parser for %s\n", filename)
			return err

		}
		err = tmpl.Execute(&b, params)
		if err != nil {
			fmt.Printf("could not generate template correctly for %s\n", filename)
			return err
		}
		err = ioutil.WriteFile(strings.TrimSuffix(filename, ".tmpl"), b.Bytes(), os.ModePerm)
		if err != nil {
			fmt.Printf("error in writing output file %s for render template %s\n", strings.TrimSuffix(filename, ".tmpl"), filename)
			return err
		}
		_ = os.Remove(filename)
	}
	return nil
}

func RenderTemplateFiles(data interface{}, directory string, skipdirectory string) error {
	fmt.Printf("running render template for %s\n", directory)
	fi, err := os.Stat(directory)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		if directory != skipdirectory {
			files, err := ioutil.ReadDir(directory)
			if err != nil {
				return err
			}
			for _, file := range files {
				filename := filepath.Join(directory, file.Name())
				if !file.IsDir() {
					err = RenderFile(filename, data)
					if err != nil {
						return err
					}
				} else {
					err = RenderTemplateFiles(data, filename, skipdirectory)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		err = RenderFile(directory, data)
		if err != nil {
			return err
		}
	}
	return nil

}

func DownloadFile(url string, filename string) error {
	resp, err := net_http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, _ := os.Create(filename)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
func Untar(targetdir string, reader io.ReadCloser) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		target := path.Join(targetdir, header.Name)
		directory, _ := path.Split(target)
		if directory != "" {
			_, err = os.Stat(directory)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				} else {
					err = os.MkdirAll(directory, os.ModePerm)
					if err != nil {
						return err
					}
				}
			}
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(target, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			err = setAttrs(target, header)
			if err != nil {
				return err
			}

		case tar.TypeReg:
			w, err := os.Create(target)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, tarReader)
			if err != nil {
				return err
			}
			w.Close()

			err = setAttrs(target, header)
			if err != nil {
				return err
			}

		default:
			log.Printf("unsupported type: %v", header.Typeflag)
		}
	}

	return nil
}

func setAttrs(target string, header *tar.Header) error {
	err := os.Chmod(target, os.FileMode(header.Mode))
	if err != nil {
		return err
	}

	return os.Chtimes(target, header.AccessTime, header.ModTime)
}

func CreateNexusDirectory(NEXUS_DIR string, NEXUS_TEMPLATE_URL string) error {
	if _, err := os.Stat(NEXUS_DIR); os.IsNotExist(err) {
		fmt.Printf("creating nexus home directory\n")
		err := os.Mkdir(NEXUS_DIR, 0755)
		if err != nil {
			return err
		}

		err = os.Chdir(NEXUS_DIR)
		if err != nil {
			return err
		}

		err = DownloadFile(NEXUS_TEMPLATE_URL, "nexus.tar")
		if err != nil {
			return fmt.Errorf("could not download template files due to %s\n", err)
		}
		file, err := os.Open("nexus.tar")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer file.Close()
		err = Untar(".", file)
		if err != nil {
			return fmt.Errorf("could not unarchive template files due to %s", err)
		}
		os.Remove("nexus.tar")
		err = os.Chdir("..")
		if err != nil {
			return err
		}
	}
	return nil
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func IsDockerRunning(cmd *cobra.Command) error {
	return SystemCommand(cmd, DOCKER_NOT_RUNNING, []string{}, "docker", "ps")
}
