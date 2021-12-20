package main

import (
	"time"
    "fmt"
	"errors"
    "os"
	"path/filepath"
)

type OSInfo struct {
    ImageFile     string
	OSVariant     string
}

var OSOptions = map[string]OSInfo{
	"ubuntu":    OSInfo{ImageFile: "ubuntu-server.img", OSVariant: "ubuntu20.04"},
}

const VolumePoolName = "vps"
const CloudImageDir = "/home/pinosaur/Temp/cloud-img"
const SnapshotDir = "/snapshots"

func buildMetadataFile(hostname string) string {
	return fmt.Sprintf(`
instance-id: %s
local-hostname: %s
`, hostname, hostname)
}

func buildUserdataFile(username string, password string, sshKey string) string {

	// TODO: validate ssh key
	if sshKey == "" {
		return fmt.Sprintf(`
#cloud-config
users:
  - name: %s
    groups: sudo
    shell: /bin/bash
    passwd: %s
    lock_passwd: false
`, username, password)
	}

	return fmt.Sprintf(`
#cloud-config
users:
  - name: %s
    groups: sudo
    shell: /bin/bash
    passwd: %s
    lock_passwd: false
    ssh_authorized_keys:
      - %s
`, username, password, sshKey)
}

func makeUserPassword(rawPassword string) (string, error) {

	cmd := []string{"mkpasswd", "--method=SHA-512", "--rounds=4096", rawPassword}
	hashedPass, err := RunCommand(cmd)

	return hashedPass, err
}

// run these as a go routine since they block
func VPSCreate(vmName string, config VPSCreateRequestData) error {

    // TODO do some validation on config to make sure there is no
    // script injection or insane settings going on

    // check if vm already exists

	// create temp dir
	tempDir, err := os.MkdirTemp("", "takoyaki-*")
	if err != nil { return err }

	fmt.Printf("%s\n", tempDir)
	defer os.RemoveAll(tempDir)

	// some vars
    cidataLocation := filepath.Join(tempDir, "cidata.iso")
    metadataLocation := filepath.Join(tempDir, "meta-data")
    userdataLocation := filepath.Join(tempDir, "user-data")
	volumeName := vmName + "-vol"

	osOptions, ok := OSOptions[config.OS]
	if !ok { return errors.New("invalid os") }

	cloudImg := filepath.Join(CloudImageDir, osOptions.ImageFile)
	osVariant := OSOptions[config.OS].OSVariant // determine based on image (full list from osinfo-query os)

    // generate meta-data and user-data files
	metadataFile := buildMetadataFile(config.Hostname)
	err = WriteToFile(metadataLocation, metadataFile)
	if err != nil { return err }

	// encrypt user password
	// TODO make sure container has mkpasswd program installed
	hashedPass, err := makeUserPassword(config.Password)
	if err != nil { return err }

	userdataFile := buildUserdataFile(config.Username, hashedPass, config.SSHKey)
	fmt.Printf("userdata file =-=-=-=-=\n%s", userdataFile)
	err = WriteToFile(userdataLocation, userdataFile)
	if err != nil { return err }

    cmd := []string{
        "genisoimage", "-output", cidataLocation, "-V",
        "cidata", "-r", "-J", userdataLocation, metadataLocation,
    }
	if err := RunCommandOnHost(cmd); err != nil { return err }

    // create volume
	cmd = []string {
		"virsh", "-c", "qemu:///system", "vol-create-as",
		VolumePoolName, volumeName, fmt.Sprintf("%dG", config.Disk), "--format", "qcow2",
	}
	if err := RunCommandOnHost(cmd); err != nil { return err }

	// load cloud image into volume
	cmd = []string {
		"virsh", "-c", "qemu:///system", "vol-upload",
		"--pool", VolumePoolName, volumeName, cloudImg,
	}
	if err := RunCommandOnHost(cmd); err != nil { return err }

    // create the vm
    cmd = []string{
		"virt-install",
		"--connect", "qemu:///system",
		"--name=" + vmName,
		"--boot", "uefi",
		"--os-variant=" + osVariant,
		"--memory=" + fmt.Sprintf("%d", config.RAM),
		"--vcpus=" + fmt.Sprintf("%d", config.CPU),
		"--import",
		"--disk", "vol=" + VolumePoolName + "/" + volumeName,
		"--disk", "path=" + cidataLocation + ",device=cdrom",
		"--graphics", "vnc,port=5911,listen=127.0.0.1", // get rid of this later
		"--noautoconsole",
    }
	if err := RunCommandOnHost(cmd); err != nil { return err }

	return nil
}

// possibly keep user data for recovery for a set amout of time
func VPSDestroy(vmName string) error {

	volumeName := vmName + "-vol"

	// possibly backup vps

	cmd := []string{
		"virsh", "-c", "qemu:///system",
		"shutdown", vmName,
	}
	if err := RunCommandOnHost(cmd); err != nil { return err }

	cmd = []string{
		"virsh", "-c", "qemu:///system",
		"destroy", vmName,
	}
	if err := RunCommandOnHost(cmd); err != nil { return err }

	cmd = []string{
		"virsh", "-c", "qemu:///system",
		"undefine", "--nvram", vmName,
	}
	if err := RunCommandOnHost(cmd); err != nil { return err }

	cmd = []string{
		"virsh", "-c", "qemu:///system",
		"vol-delete", "--pool", VolumePoolName, volumeName,
	}
	if err := RunCommandOnHost(cmd); err != nil { return err }

	return nil
}

// when a user requests for vps specs to be upgraded
func VPSUpgrade() error {

	return nil
}

func VPSSnapshot(vmName string) error {

	now := time.Now().String()
	snapshotName := fmt.Sprintf("snapshot-%s-%s", vmName, now)

	cmd := []string{
		"virsh", "-c", "qemu:///system", "snapshot-create-as",
		"--domain", vmName,
		"--name", filepath.Join(SnapshotDir, snapshotName),
    }
	if err := RunCommandOnHost(cmd); err != nil { return err }

	return nil
}

func VPSStart(vmName string) error {

	cmd := []string{
		"virsh", "-c", "qemu:///system", "start", vmName,
    }
	if err := RunCommandOnHost(cmd); err != nil { return err }

	return nil
}

func VPSStop(vmName string) error {

	cmd := []string{
		"virsh", "-c", "qemu:///system", "shutdown", vmName,
    }
	if err := RunCommandOnHost(cmd); err != nil { return err }

	return nil
}

func VPSRestart(vmName string) error {

	return nil
}

