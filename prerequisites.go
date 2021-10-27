package main

import "github.com/jiuncheng/archinstall2/utils"

func RunPrerequisites() error {
	// Prerequisite
	err := utils.NewCmd("timedatectl set-ntp true").SetDesc("Syncing datetime to ntp server...").Run()
	if err != nil {
		return err
	}

	err = utils.NewCmd("reflector -a 48 -c sg -f 5 -l 20 --verbose --sort rate --save /etc/pacman.d/mirrorlist").SetDesc("Finding fastest mirror...").Run()
	if err != nil {
		return err
	}

	err = utils.NewCmd("sed -i s/^#Para/Para/ /etc/pacman.conf").SetDesc("Enabling parallel download..").Run()
	if err != nil {
		return err
	}

	utils.NewCmd("umount -a").SetDesc("Umounting all drives").Run()

	return nil
}
