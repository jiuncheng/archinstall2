package main

import (
	"github.com/jiuncheng/archinstall2/cmd"
)

func RunPrerequisites() error {
	// Prerequisite
	err := cmd.NewCmd("timedatectl set-ntp true").SetDesc("Syncing datetime to ntp server...").Run()
	if err != nil {
		return err
	}

	err = cmd.NewCmd("reflector -a 48 -c sg -f 5 -l 20 --verbose --sort rate --save /etc/pacman.d/mirrorlist").SetDesc("Finding fastest mirror...").Run()
	if err != nil {
		return err
	}

	err = cmd.NewCmd("sed -i s/^#Para/Para/ /etc/pacman.conf").SetDesc("Enabling parallel download..").Run()
	if err != nil {
		return err
	}

	cmd.NewCmd("umount -a").SetDesc("Umounting all drives").Run()

	return nil
}
