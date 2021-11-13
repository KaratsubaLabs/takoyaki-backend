package main

import (
    "fmt"
    _ "os"
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

func VPSCreate(config VPSConfig) {

    // TODO do some validation on config to make sure there is no
    // script injection or insane settings going on

    // check if vm already exists

	// generate random name for vm
	vmName := RandomString()
	tempDir := "/tmp/takoyaki/" + vmName

	// create temp dir

	// some vars
    cidataLocation := tempDir + "cidata.iso"
    userdataLocation := tempDir + "user-data"
    metadataLocation := tempDir + "meta-data"
	volumeName := vmName + "-vol"
	// make sure config.OS is valid
	cloudImg := OSOptions[config.OS].ImageFile // also concat image location
	osVariant := OSOptions[config.OS].OSVariant // determine based on image (full list from osinfo-query os)

    // generate meta-data and user-data files
    fmt.Sprintf(`
        instance-id: %s
        local-hostname: %s
    `, config.Hostname, config.Hostname)

    fmt.Sprintf(`
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

    // create cidata image (maybe do it in /temp ?)

    cmd := []string{
        "genisoimage", "-output", cidataLocation, "-V",
        "cidata", "-r", "-J", userdataLocation, metadataLocation,
    }

    // create volume
	cmd = []string {
		"virsh", "-c", "qemu:///system", "vol-create-as",
		VolumePoolName, volumeName, fmt.Sprintf("%d", config.Disk), "--format", "qcow2",
	}

	// load cloud image into volume
	cmd = []string {
		"virsh", "-c", "qemu:///system", "vol-upload",
		"--pool", VolumePoolName, volumeName, cloudImg,
	}

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

    _ = cmd

	// clean up temp dir

}

// possibly keep user data for recovery for a set amout of time
func VPSDestroy() {

}

// when a user requests for vps specs to be upgraded
func VPSModify() {

}

