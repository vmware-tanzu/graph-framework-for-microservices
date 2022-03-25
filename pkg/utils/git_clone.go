package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
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
	fmt.Printf("Intializing gomodule")
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

func RenderTemplateFiles(data interface{}, directory string) error {
	fmt.Printf("running render template for %s\n", directory)
	fi, err := os.Stat(directory)
	if err != nil {
		return err
	}
	if fi.IsDir() {
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
				err = RenderTemplateFiles(data, filename)
				if err != nil {
					return err
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
