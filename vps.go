package main

import (
    "fmt"
    "os"
)

const (
    RAM_LOW    = 512
    RAM_MEDIUM = 1024
    RAM_HIGH   = 2048
)

type osInfo struct {
    ImageFile     string
}

const osOptions = map[string]osInfo{
    'ubuntu':    osInfo{ImageFile: "ubuntu-server.img"}
}

type VPSConfig struct {
    Hostname      string
    Username      string
    Password      string
    SSHKey        string
    RAM           int // make this 'enum' or sm
    OS            string
}

func VPSCreate(config VPSConfig) {

    // TODO do some validation on config to make sure there is no
    // script injection or insane settings going on

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
        "cidata", "-r", "-J", userdataLocation, metadataLocation
    }

    // create disk image

    // create the vm
    cmd = []string{
        "virt-install"
    }

}

func VPSDelete() {

    cmd := []string{

    }

}

