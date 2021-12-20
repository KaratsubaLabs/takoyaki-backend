package main

import (
	"os"
	"os/exec"
	"bytes"
	"time"
	"math/rand"
)

func ContainsString(slice []string, elem string) bool {
	for _, a := range slice {
		if a == elem {
			return true
		}
	}
	return false
}

const randomStringLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const randomStringLen = 10
func RandomString() string {

	rand.Seed(time.Now().UnixNano())

	out := make([]byte, randomStringLen)
	for i := range out {
		out[i] = randomStringLetters[rand.Intn(len(randomStringLetters))]
	}
	return string(out)

}

func WriteToFile(filepath string, content string) error {

	f, err := os.Create(filepath)
	if err != nil { return err }

	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil { return err }

	return nil
}

func RunCommand(args []string) (string, error) {

	cmd := exec.Command(args[0], args[1:]...)

	var output bytes.Buffer

	cmd.Stdout = &output
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil { return "", err }

	return output.String(), nil
}

// RUN COMMAND ON HOST
// TODO: need to also pass output back inside container so we know when
// the command has finished executing
const HOST_PIPE = "/var/run/takoyaki"
func RunCommandOnHost(args []string) error {

	cmd := exec.Command("echo", args[0:]...)

	pf, err := os.OpenFile(HOST_PIPE, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeNamedPipe)
	if err != nil { return err }

	defer pf.Close()

	cmd.Stdout = pf
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil { return err }

	return nil
}

