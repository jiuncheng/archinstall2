package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/jiuncheng/archinstall2/cmd"
	"github.com/jiuncheng/archinstall2/disklist"
)

var (
	installDisk string
)

func main() {
	command := exec.Command("lsblk", "-ldnJe", "7,11")
	res, err := command.Output()
	if err != nil {
		log.Fatalln(err.Error())
	}

	dl, err := disklist.NewDiskListFromJSON(res)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var result int
	for {
		for i, disk := range dl.BlockDevices {
			fmt.Printf("%d.    /dev/%s    %s\n", i+1, disk.Name, disk.Size)
		}
		fmt.Println("\n Please select the number which the disk will be used for installation (e.g. 3): ")

		_, err = fmt.Scan(&result)
		if err == nil {
			break
		}
		fmt.Println("\n\nOnly number is allowed.")
	}

	installDisk = "/dev/" + dl.BlockDevices[result-1].Name

	err = cmd.NewCmd("sgdisk -o " + installDisk).SetDesc("Formatting selected disk...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("sgdisk -n 1:0+500M -t 1:ef00 -c 1:'EFI' " + installDisk).Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("sgdisk -n 2:0:-4G -t 2:8300 -c 2:'rootfs' " + installDisk).Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mkfs.vfat -n EFI " + installDisk + "1").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mkfs.btrfs -f -L ROOT " + installDisk + "2").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mkfs.btrfs -f -L ROOT " + installDisk + "2").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mount " + installDisk + "2" + " /mnt").SetDesc("Mounting disk for btrfs subvolume creation...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("btrfs subvolume create /mnt/@").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("btrfs subvolume create /mnt/@home").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("btrfs subvolume create /mnt/@snapshots").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("btrfs subvolume create /mnt/@var_log").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("umount /mnt").SetDesc("Unmounting disk for remounting with BTRFS options...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mount -o noatime,compress=zstd,space_cache,discard=async,subvol=@ " + installDisk + "2" + " /mnt").SetDesc("Remounting with BTRFS options...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mkdir -p /mnt/boot /mnt/home /mnt/var /mnt/.snapshots").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mkdir -p /mnt/var/log").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mount -o noatime,compress=zstd,space_cache,discard=async,subvol=@home " + installDisk + "2" + " /mnt/home").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mount -o noatime,compress=zstd,space_cache,discard=async,subvol=@snapshots " + installDisk + "2" + " /mnt/.snapshots").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mount -o noatime,compress=zstd,space_cache,discard=async,subvol=@var_log " + installDisk + "2" + " /mnt/var/log").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("mount " + installDisk + "1" + " /mnt/boot").SetDesc("Mounting EFI /boot...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd2 := cmd.NewCmd("pacstrap /mnt base base-devel linux linux-firmware intel-ucode git neovim nano btrfs-progs")
	err = cmd2.SetDesc("Downloading packages from Pacstrap...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("genfstab -U /mnt >> /mnt/etc/fstab").SetDesc("Generating FSTAB file...").Run()
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

	err = cmd.NewCmd("arch-chroot /mnt sed -i '177s/.//' /etc/locale.gen").SetDesc("Generating en_US_utf-8 locale...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt locale-gen").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt echo 'LANG=en_US.UTF-8' >> /etc/locale.conf").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt echo 'archlinux' >> /etc/hostname").SetDesc("Setting hostname as archlinux...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt echo '127.0.0.1    localhost' >> /etc/hosts").SetDesc("Configuring /etc/hosts...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt echo '::1    localhost' >> /etc/hosts").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt echo '127.0.1.1    archlinux.localdomain archlinux' >> /etc/hosts").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Print("\nPlease enter a root password: ")
	var pwd string
	_, err = fmt.Scanln(&result)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("New root password is ", pwd)

	err = cmd.NewCmd("arch-chroot /mnt echo root:" + pwd + " | " + "chpasswd").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Installation done. You will now be able to reboot.")
}
