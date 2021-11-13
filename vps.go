package main

import (
    "fmt"
    "os"
	"os/exec"
	"path/filepath"
)

const (
    RAM_LOW    = 512
    RAM_MEDIUM = 1024
    RAM_HIGH   = 2048
)

type VPSConfig struct {
    DisplayName   string
    Hostname      string
    Username      string
    Password      string
    SSHKey        string
    RAM           int // make this 'enum' or sm
    CPU           int
    Disk          int
    OS            string
}

type OSInfo struct {
    ImageFile     string
	OSVariant     string
}

var OSOptions = map[string]OSInfo{
	"ubuntu":    OSInfo{ImageFile: "ubuntu-server.img", OSVariant: "ubuntu20.04"},
}

const VolumePoolName = "vps"

func writeToFile(filepath string, content string) error {

	f, err := os.Create(filepath)
	if err != nil { return err }

	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil { return err }

	return nil
}

func runCommand(args []string) error {

	cmd := exec.Command(args[0], args[1:]...)

	// log these instead
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Wait()
	if err != nil { return err }

	return nil
}

// run these as a go routine since they block
func VPSCreate(config VPSConfig) error {

    // TODO do some validation on config to make sure there is no
    // script injection or insane settings going on

    // check if vm already exists

	// generate random name for vm
	vmName := RandomString()

	// create temp dir
	tempDir, err := os.MkdirTemp("", "takoyaki-*")
	if err != nil { return err }

	defer os.RemoveAll(tempDir)

	// some vars
    cidataLocation := filepath.Join(tempDir, "cidata.iso")
    metadataLocation := filepath.Join(tempDir, "meta-data")
    userdataLocation := filepath.Join(tempDir, "user-data")
	volumeName := vmName + "-vol"
	// make sure config.OS is valid
	cloudImg := OSOptions[config.OS].ImageFile // also concat image location
	osVariant := OSOptions[config.OS].OSVariant // determine based on image (full list from osinfo-query os)

    // generate meta-data and user-data files
	metadataFile := fmt.Sprintf(`
        instance-id: %s
        local-hostname: %s
    `, config.Hostname, config.Hostname)
	err = writeToFile(metadataLocation, metadataFile)
	if err != nil { return err }

	userdataFile := fmt.Sprintf(`
        #cloud-config
        users:
          - name: %s
            ssh_authorized_keys:
              - %s
            groups: sudo
            shell: /bin/bash
            passwd: %s
            lock_passwd: false
    `, config.Username, config.SSHKey, config.Password)
	err = writeToFile(userdataLocation, userdataFile)
	if err != nil { return err }

    cmd := []string{
        "genisoimage", "-output", cidataLocation, "-V",
        "cidata", "-r", "-J", userdataLocation, metadataLocation,
    }
	if err := runCommand(cmd); err != nil { return err }

    // create volume
	cmd = []string {
		"virsh", "-c", "qemu:///system", "vol-create-as",
		VolumePoolName, volumeName, fmt.Sprintf("%d", config.Disk), "--format", "qcow2",
	}
	if err := runCommand(cmd); err != nil { return err }

	// load cloud image into volume
	cmd = []string {
		"virsh", "-c", "qemu:///system", "vol-upload",
		"--pool", VolumePoolName, volumeName, cloudImg,
	}
	if err := runCommand(cmd); err != nil { return err }

    // create the vm
    cmd = []string{
		"virt-install", "-c", "qemu:///system",
		"--name=" + vmName,
		"--boot", "uefi",
		"--os-variant=", osVariant,
		"--memory=" + fmt.Sprintf("%d", config.RAM),
		"--vcpus=" + fmt.Sprintf("%d", config.CPU),
		"--import",
		"--disk", "vol=" + VolumePoolName + "/" + volumeName,
		"--disk", "path=" + cidataLocation + ",device=cdrom",
    }
	if err := runCommand(cmd); err != nil { return err }

	return nil
}

// possibly keep user data for recovery for a set amout of time
func VPSDestroy(vmName string) error {

	volumeName := vmName + "-vol"

	// possibly backup vps

	cmd := []string{
		"virt-install", "-c", "qemu:///system",
		"shutdown", vmName,
	}
	if err := runCommand(cmd); err != nil { return err }

	cmd = []string{
		"virt-install", "-c", "qemu:///system",
		"destroy", vmName,
	}
	if err := runCommand(cmd); err != nil { return err }

	cmd = []string{
		"virt-install", "-c", "qemu:///system",
		"undefine", "--nvram", vmName,
	}
	if err := runCommand(cmd); err != nil { return err }

	cmd = []string{
		"virt-install", "-c", "qemu:///system",
		"vol-delete", "--pool", VolumePoolName, volumeName,
	}
	if err := runCommand(cmd); err != nil { return err }

	return nil
}

// when a user requests for vps specs to be upgraded
func VPSModify() {

}

func VPSBackup() {

}

