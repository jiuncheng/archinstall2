package main

import (
	"fmt"
	"log"

	"github.com/jiuncheng/archinstall2/cmd"
	"github.com/jiuncheng/archinstall2/disklist"
	"github.com/jiuncheng/archinstall2/diskutil"
	"github.com/jiuncheng/archinstall2/filesystem"
	"github.com/jiuncheng/archinstall2/sysconfig"
)

func main() {
	err := cmd.NewCmd("timedatectl set-ntp true").SetDesc("Syncing datetime to ntp server...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cfg := sysconfig.NewSysConfig()

	dl, err := disklist.GetDiskList()
	if err != nil {
		log.Fatalln(err.Error())
	}

	var result int
	for {
		for i, disk := range dl {
			fmt.Printf("%d.    /dev/%s    %s\n", i+1, disk.Name, disk.Size)
		}
		fmt.Print("\n Please select the number which the disk will be used for installation (e.g. 1): ")

		_, err = fmt.Scanf("%d", &result)
		if err == nil {
			if result <= len(dl) && result > 0 {
				break
			}
			fmt.Println("\n\nThe number must be between 1 and ", len(dl), ".")
			fmt.Print("Press enter to choose again : ")
			fmt.Scanln()
			continue
		}
		fmt.Println("\n\nOnly number between 1 and ", len(dl), " is allowed.")
	}

	cfg.InstallDisk = "/dev/" + dl[result-1].Name

	err = diskutil.NewDiskUtil(cfg).CreateBTRFS()
	if err != nil {
		log.Fatalln(err.Error())
	}
	filesystem.NewBtrfsHelper(cfg).GenerateBTRFSSystem()

	cmd2 := cmd.NewCmd("pacstrap /mnt base base-devel linux linux-firmware intel-ucode git neovim nano btrfs-progs")
	err = cmd2.SetDesc("Downloading packages from Pacstrap...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/sh -c \"genfstab -U /mnt >> /mnt/etc/fstab\"").SetDesc("Generating FSTAB file...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt ln -sf /usr/share/zoneinfo/Asia/Kuala_Lumpur /etc/localtime").SetDesc("Symlinking Asia/Kuala Lumpur time to /etc/localtime...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt timedatectl set-ntp true").SetDesc("Syncing datetime to ntp server...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt hwclock --systohc").SetDesc("Setting hardware clock...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt sed -i s/^#*\\(en_US.UTF-8\\)/\\1/ /etc/locale.gen").SetDesc("Generating en_US_utf-8 locale...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt locale-gen").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/sh -c \"arch-chroot /mnt echo 'LANG=en_US.UTF-8' >> /etc/locale.conf\"").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/sh -c \"arch-chroot /mnt echo 'archlinux' >> /etc/hostname\"").SetDesc("Setting hostname as archlinux...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/sh -c \"arch-chroot /mnt echo '127.0.0.1    localhost' >> /etc/hosts\"").SetDesc("Configuring /etc/hosts...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/sh -c \"arch-chroot /mnt echo '::1    localhost' >> /etc/hosts\"").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/sh -c \"arch-chroot /mnt echo '127.0.1.1    archlinux.localdomain archlinux' >> /etc/hosts\"").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Print("\nPlease enter a root password: ")
	var pwd string
	_, err = fmt.Scanln(&pwd)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("New root password is ", pwd)

	err = cmd.NewCmd("/bin/sh -c \"arch-chroot /mnt echo root:" + pwd + " | " + "chpasswd\"").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt useradd -mG wheel -s /bin/bash -p 12345 home3").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// fmt.Println(cfg.InstallDisk)
	fmt.Println("Installation done. You will now be able to reboot.")
}
