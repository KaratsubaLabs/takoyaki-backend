package util

import (
	"bytes"
	"os"
	"os/exec"
)

func WriteToFile(filepath string, content string) error {

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func RunCommand(args []string) (string, error) {

	cmd := exec.Command(args[0], args[1:]...)

	var output bytes.Buffer

	cmd.Stdout = &output
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

const HOST_PIPE = "/var/run/takoyaki/pipe"

func RunCommandOnHost(args []string) error {

	cmd := exec.Command("echo", args[0:]...)

	pf, err := os.OpenFile(HOST_PIPE, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeNamedPipe)
	if err != nil {
		return err
	}

	defer pf.Close()

	cmd.Stdout = pf
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
