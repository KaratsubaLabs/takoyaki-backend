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
    STORAGE       int
    OS            string
}

type OSInfo struct {
    ImageFile     string
}

var OSOptions = map[string]OSInfo{
    "ubuntu":    OSInfo{ImageFile: "ubuntu-server.img"},
}

func VPSCreate(config VPSConfig) {

    // TODO do some validation on config to make sure there is no
    // script injection or insane settings going on

    // check if vm already exists

    // generate user-data and meta-data files
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
    cidataLocation := "cidata.iso"
    userdataLocation := "user-data"
    metadataLocation := "meta-data"

    cmd := []string{
        "genisoimage", "-output", cidataLocation, "-V",
        "cidata", "-r", "-J", userdataLocation, metadataLocation,
    }

    // create disk image

    // create the vm
    cmd = []string{
        "virt-install",
            // "--connect=qemu:///system",
            // "--name=" + config.Hostname,
            // "--memory=" + config.RAM,
            // "--vcpus=" + config.CPU,
            // "--boot uefi",
    }

    _ = cmd

}

// possibly keep user data for recovery for a set amout of time
func VPSDestroy() {

}

// when a user requests for vps specs to be upgraded
func VPSModify() {

}

