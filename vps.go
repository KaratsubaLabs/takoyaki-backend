package main

import (
    "os"
)

type VPSConfig struct {
    Hostname      string
    Username      string
    Password      string
    SSHKey        string
}

func VPSCreate() {

    // generate user-data and meta-data files

    // create cidata image

    // create disk image

    // create the vm
    cmd := []string{
        "virt-install"
    }

}

func VPSDelete() {

    cmd := []string{

    }

}

