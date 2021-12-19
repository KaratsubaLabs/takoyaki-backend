package main

import (
	"time"
	"bytes"
	"os"
	"os/exec"
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

// RUN COMMAND ON HOST
// TODO: need to also pass output back inside container so we know when
// the command has finished executing
/*
func RunCommandOnHost(args []string) error {

	cmd := exec.Command("echo", args[0:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil { return err }

	return nil
}
*/

func RunCommand(args []string) error {

	cmd := exec.Command(args[0], args[1:]...)

	// log these instead
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil { return err }

	return nil
}

func RunCommandWithOutput(args []string) (string, error) {

	cmd := exec.Command(args[0], args[1:]...)

	var output bytes.Buffer

	cmd.Stdout = &output
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil { return "", err }

	return output.String(), nil
}

