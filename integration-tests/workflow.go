package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Workflow []string

var backgroundProcesses []*exec.Cmd // store handles to any background processes

func RunWorkflow(workflow Workflow) {
	for _, file := range workflow {
		err := ProcessWorkflow(filepath.Join(pwd, file))
		CheckIfError(err)
	}
}

func ProcessWorkflow(path string) error {
	fmt.Printf("Opening file %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "```" {
			for scanner.Scan() {
				line = scanner.Text()
				if strings.TrimSpace(line) == "```" {
					break
				}
				err := RunCmd(line)
				if err != nil {
					return err
				}
			}
		} else if strings.Contains(line, "```") {
			// skip annotated code blocks like ```golang
			for scanner.Scan() {
				line = scanner.Text()
				if strings.Contains(line, "```") {
					break
				}
			}
		}
	}
	// kill any background processes created
	for _, backgroundProcess := range backgroundProcesses {
		backgroundProcess.Process.Kill()
	}
	return nil
}

func RunCmd(line string) error {
	if strings.HasPrefix(line, "#") {
		return nil
	}
	fmt.Println(line)
	// TODO support ; as a separator of commands
	for _, command := range strings.Split(line, "&&") {
		if strings.TrimSpace(line) == "" { // ignore empty lines
			continue
		}
		// fill in any env vars
		parts := strings.Fields(command)
		if len(parts) <= 1 {
			return fmt.Errorf("invalid format: too few parameters")
		}

		var err error
		if parts[0] == "cd" || parts[0] == "/usr/bin/cd" {
			fmt.Printf("Executing `%v`\n", append([]string{parts[0]}, parts[1]))
			// TODO - this is far from perfect but OK for now
			err = os.Chdir(os.ExpandEnv(parts[1]))
			continue
		}

		if parts[0] == "export" {
			fmt.Printf("Executing `%v`\n", append([]string{parts[0]}, parts[1]))
			tmp := strings.Split(parts[1], "=")
			if len(tmp) == 2 {
				var value string
				if strings.HasPrefix(tmp[1], "$") {
					startOfVar := strings.Split(tmp[1], ":")
					if len(startOfVar) > 1 {
						varName := strings.TrimPrefix(startOfVar[0], "${")
						if os.Getenv(varName) == "" {
							value = strings.TrimSuffix(strings.TrimPrefix(startOfVar[1], "-"), "}")
						} else {
							value = os.Getenv(varName)
						}

					} else {
						value = os.ExpandEnv(tmp[1])
					}
					fmt.Printf("%s:%s\n", tmp[0], value)
					if err = os.Setenv(tmp[0], value); err != nil {
						return err
					}
				} else {
					if err = os.Setenv(tmp[0], os.ExpandEnv(tmp[1])); err != nil {
						return err
					}
				}
			} else if len(tmp) == 1 {
				if err = os.Setenv(tmp[0], ""); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Invalid export found: %s\n", line)
			}
			continue
		}

		if parts[0] == "kubectl" && len(parts) > 1 && parts[1] == "port-forward" {
			var args []string
			for _, arg := range parts[1:] {
				args = append(args, os.ExpandEnv(arg))
			}
			fmt.Printf("Executing `%v`\n", append([]string{parts[0]}, args...))
			cmd := exec.Command(parts[0], args...)
			err = cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			backgroundProcesses = append(backgroundProcesses, cmd)
			time.Sleep(5 * time.Second)
			continue
		}

		var args []string
		for _, arg := range parts[1:] {
			args = append(args, os.ExpandEnv(arg))
		}

		if parts[0] == "nexus" {
			args = append(args, "--debug")
		}

		fmt.Printf("Executing `%s`\n", append([]string{parts[0]}, args...))
		execCmd := exec.Command(parts[0], args...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stdout
		err = execCmd.Run()
		if err != nil {
			return fmt.Errorf("Failed to run `%s` due to error: %s\n", command, err.Error())
		}

		// add a small delay to offset any timing issues for kubectl commands
		if parts[0] == "kubectl" {
			time.Sleep(5 * time.Second)
		}
	}
	return nil
}

func cleanupCluster() error {
	err := RunCmd("nexus runtime uninstall -n $NAMESPACE")
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Second) // give uninstall some time...

	err = RunCmd("kubectl delete deployment --all -n $NAMESPACE")
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Second)

	return nil
}
