package vps

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/KaratsubaLabs/takoyaki-backend/db"
	"github.com/KaratsubaLabs/takoyaki-backend/util"
)

type OSInfo struct {
	ImageFile string
	OSVariant string
}

var OSOptions = map[string]OSInfo{
	"ubuntu": OSInfo{ImageFile: "ubuntu-server.img", OSVariant: "ubuntu20.04"},
}

const VolumePoolName = "vps"
const CloudImageDir = "/home/pinosaur/Temp/cloud-img"
const SnapshotDir = "/snapshots"

// TODO: limitation, buildMetadataFile and buildUserdataFile can not contain new lines
// possibly move the entirety of create vps to host side and have takoyaki send the vps data over directly
func buildMetadataFile(hostname string) string {
	return fmt.Sprintf("local-hostname: %s", hostname)
}

func buildUserdataFile(username string, password string, sshKey string) string {

	// TODO: validate ssh key
	if sshKey == "" {
		return fmt.Sprintf("users: [{ name: %s, groups: sudo, shell: /bin/bash, passwd: %s, lock_passwd: false }]", username, password)
	}

	return fmt.Sprintf("users: [{ name: %s, groups: sudo, shell: /bin/bash, passwd: %s, lock_passwd: false, ssh_authorized_keys: [ %s ] }]", username, password, sshKey)
}

func makeUserPassword(rawPassword string) (string, error) {

	cmd := []string{"mkpasswd", "--method=SHA-512", "--rounds=4096", rawPassword}
	hashedPass, err := util.RunCommand(cmd)

	return strings.TrimSuffix(hashedPass, "\n"), err
}

// run these as a go routine since they block
func Create(vmName string, config db.VPSCreateRequestData) error {

	// TODO do some validation on config to make sure there is no
	// script injection or insane settings going on

	// check if vm already exists

	// create temp dir
	tempDir := fmt.Sprintf("/tmp/takoyaki-%s", vmName)
	cmd := []string{
		"mkdir", tempDir,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	// some vars
	cidataLocation := filepath.Join(tempDir, "cidata.iso")
	metadataLocation := filepath.Join(tempDir, "meta-data")
	userdataLocation := filepath.Join(tempDir, "user-data")
	volumeName := vmName + "-vol"

	osOptions, ok := OSOptions[config.OS]
	if !ok {
		return errors.New("invalid os")
	}

	cloudImg := filepath.Join(CloudImageDir, osOptions.ImageFile)
	osVariant := OSOptions[config.OS].OSVariant // determine based on image (full list from osinfo-query os)

	// generate meta-data and user-data files
	metadataFile := buildMetadataFile(config.Hostname)
	cmd = []string{
		"echo", "'" + metadataFile + "'", ">", metadataLocation,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	// encrypt user password
	// TODO make sure container has mkpasswd program installed
	hashedPass, err := makeUserPassword(config.Password)
	if err != nil {
		return err
	}

	userdataFile := buildUserdataFile(config.Username, hashedPass, config.SSHKey)
	fmt.Printf("userdata file =-=-=-=-=\n%s", userdataFile)
	cmd = []string{
		"echo", "'" + userdataFile + "'", ">", userdataLocation,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	cmd = []string{
		"genisoimage", "-output", cidataLocation, "-V",
		"cidata", "-r", "-J", userdataLocation, metadataLocation,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	// create volume
	cmd = []string{
		"virsh", "-c", "qemu:///system", "vol-create-as",
		VolumePoolName, volumeName, fmt.Sprintf("%dG", config.Disk), "--format", "qcow2",
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	// load cloud image into volume
	cmd = []string{
		"virsh", "-c", "qemu:///system", "vol-upload",
		"--pool", VolumePoolName, volumeName, cloudImg,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

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
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	cmd = []string{
		"rm", "-rf", tempDir,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	return nil
}

// possibly keep user data for recovery for a set amout of time
func Destroy(vmName string) error {

	volumeName := vmName + "-vol"

	// possibly backup vps

	cmd := []string{
		"virsh", "-c", "qemu:///system",
		"shutdown", vmName,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	cmd = []string{
		"virsh", "-c", "qemu:///system",
		"destroy", vmName,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	cmd = []string{
		"virsh", "-c", "qemu:///system",
		"undefine", "--nvram", vmName,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	cmd = []string{
		"virsh", "-c", "qemu:///system",
		"vol-delete", "--pool", VolumePoolName, volumeName,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	return nil
}

// when a user requests for vps specs to be upgraded
func Upgrade() error {

	return nil
}

func Snapshot(vmName string) error {

	now := time.Now().String()
	snapshotName := fmt.Sprintf("snapshot-%s-%s", vmName, now)

	cmd := []string{
		"virsh", "-c", "qemu:///system", "snapshot-create-as",
		"--domain", vmName,
		"--name", filepath.Join(SnapshotDir, snapshotName),
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	return nil
}

func Start(vmName string) error {

	cmd := []string{
		"virsh", "-c", "qemu:///system", "start", vmName,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	return nil
}

func Stop(vmName string) error {

	cmd := []string{
		"virsh", "-c", "qemu:///system", "shutdown", vmName,
	}
	if err := util.RunCommandOnHost(cmd); err != nil {
		return err
	}

	return nil
}

// maybe don't need this
func Restart(vmName string) error {

	return nil
}
