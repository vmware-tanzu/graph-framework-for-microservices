package utils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	net_http "net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func GitClone() (string, error) {
	fmt.Println("git clone https://gitlab.eng.vmware.com/nsx-allspark_users/m7/policymodel.git")

	manifestDir := "policymodel" // Should remove this directory after installation is complete.
	_, err := os.Stat(manifestDir)
	if err == nil {
		os.RemoveAll(manifestDir)
	}
	err = os.Mkdir(manifestDir, os.ModePerm)
	if err != nil {
		fmt.Printf("tenant-manifest directory creation failed with an error: %v", err)
		return "", err
	}
	var branch plumbing.ReferenceName = "refs/heads/nexus-sdk-dev"
	var Auth transport.AuthMethod
	var URL string = ""
	if os.Getenv("GIT_USER") != "" && os.Getenv("GIT_PASS") != "" {
		URL = "https://gitlab.eng.vmware.com/nsx-allspark_users/m7/policymodel"
		Auth = &http.BasicAuth{Username: os.Getenv("GIT_USER"), Password: os.Getenv("GIT_PASS")}
	} else {
		var path string
		URL = "git@gitlab.eng.vmware.com:nsx-allspark_users/m7/policymodel.git"
		if os.Getenv("GIT_SSH_KEY") == "" {
			homedir, _ := os.UserHomeDir()
			path = filepath.Join(homedir, ".ssh/id_rsa")
		} else {
			path = os.Getenv("GIT_SSH_KEY")
		}
		pvk, _ := ioutil.ReadFile(path)
		signer, _ := ssh.ParsePrivateKey(pvk)
		Auth = &ssh2.PublicKeys{User: "git", Signer: signer}
	}

	r, err := git.PlainClone(manifestDir, false, &git.CloneOptions{
		URL:           URL,
		Progress:      os.Stdout,
		Auth:          Auth,
		Depth:         1,
		ReferenceName: branch,
		SingleBranch:  true,
	})
	if err != nil {
		fmt.Printf("policymodel git-clone failed with an error: %v", err)
		return "", err
	}
	_ = r.DeleteRemote("origin")

	//remove remote origin...

	return manifestDir, nil
}

func GoModInit(path string) error {
	fmt.Printf("Intializing gomodule\n")
	os.Chdir(path)
	cmd := exec.Command("go", "mod", "init", path)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	os.Chdir("..")
	return nil
}

func GetModuleName(path string) (string, error) {
	fmt.Printf("getting the current modulename")
	os.Chdir(path)
	cmd := exec.Command("go", "list", "-m")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	os.Chdir("..")
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

func SystemCommand(envList []string, name string, args ...string) error {
	fmt.Printf("envList: %v\n", envList)
	fmt.Printf("command: %v\n", name)
	fmt.Printf("args: %v\n", args)

	command := exec.Command(name, args...)
	command.Env = os.Environ()

	if len(envList) > 0 {
		command.Env = append(command.Env, envList...)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf(err.Error())

	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	err = command.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return err
	}

	err = command.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return err
	}

	return nil
}

func DownloadFile(url string, filename string) error {
	url = fmt.Sprintf("%s", url)
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
func Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			if _, err = os.Stat(target); os.IsNotExist(err) {
				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))

				if err != nil {
					return err
				}

				// copy over contents
				if _, err := io.Copy(f, tr); err != nil {
					return err
				}

				// manually close here after each file operation; defering would cause each file close
				// to wait until all operations have completed.
				f.Close()
			}
		}
	}
}

func CreateNexusDirectory(NEXUS_DIR string, NEXUS_TEMPLATE_URL string) error {
	fmt.Print("run this command outside of nexus home directory\n")
	if _, err := os.Stat(NEXUS_DIR); os.IsNotExist(err) {
		fmt.Printf("creating nexus home directory\n")
		os.Mkdir(NEXUS_DIR, 0755)
	}
	os.Chdir(NEXUS_DIR)
	err := DownloadFile(NEXUS_TEMPLATE_URL, "nexus.tar")
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
	os.Chdir("..")
	return nil
}
