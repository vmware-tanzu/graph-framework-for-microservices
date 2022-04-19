package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func SystemCommand(cmd *cobra.Command, customErr ClientErrorCode, envList []string, name string, args ...string) error {
	// Get silent flag status.
	debugEnabled := IsDebug(cmd)

	if debugEnabled {
		if len(envList) != 0 {
			fmt.Printf("envList: %v\n", envList)
		}
		fmt.Printf("command: %v\n", name)
		fmt.Printf("args: %v\n", args)
	}
	command := exec.Command(name, args...)
	command.Env = os.Environ()

	if len(envList) > 0 {
		command.Env = append(command.Env, envList...)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return GetCustomError(INTERNAL_ERROR, err).Print().ExitIfFatalOrReturn()
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return GetCustomError(INTERNAL_ERROR, err).Print().ExitIfFatalOrReturn()
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			if debugEnabled {
				fmt.Printf("\t > %s\n", scanner.Text())
			}
		}
	}()

	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			if debugEnabled {
				fmt.Printf("\t > %s\n", errScanner.Text())
			}
		}
	}()
	err = command.Start()
	if err != nil {
		return GetCustomError(customErr,
			fmt.Errorf("starting cmd %s %v failed with error %v", name, args, err)).
			Print().ExitIfFatalOrReturn()
	}

	err = command.Wait()
	if err != nil {
		return GetCustomError(customErr,
			fmt.Errorf("waiting for cmd %s %v failed with error %v", name, args, err)).
			Print().ExitIfFatalOrReturn()

	}

	return nil
}
