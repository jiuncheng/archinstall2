package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"

	"github.com/jiuncheng/archinstall2/disklist"
)

var (
	installDisk string
)

func main() {
	cmd := exec.Command("lsblk", "-ldnJe", "7,11")
	res, err := cmd.Output()
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

	fmt.Println("\nFormatting selected disk...")
	cmd = exec.Command("sgdisk", "-o", installDisk)
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("sgdisk", "-n", "1:0:+500M", "-t", "1:ef00", "-c", "1:'EFI'", installDisk)
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("sgdisk", "-n", "2:0:-4G", "-t", "2:8300", "-c", "2:'rootfs'", installDisk)
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mkfs.vfat", "-n", "EFI", installDisk+"1")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mkfs.btrfs", "-f", "-L", "ROOT", installDisk+"2")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mkfs.btrfs", "-f", "-L", "ROOT", installDisk+"2")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nMounting disk for btrfs subvolume creation...")
	cmd = exec.Command("mount", installDisk+"2", "/mnt")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("btrfs", "subvolume", "create", "/mnt/@")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("btrfs", "subvolume", "create", "/mnt/@home")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("btrfs", "subvolume", "create", "/mnt/@snapshots")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("btrfs", "subvolume", "create", "/mnt/@var_log")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nUnmounting disk for remounting with BTRFS options")
	cmd = exec.Command("umount", "/mnt")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nRemounting with BTRFS options")
	cmd = exec.Command("mount", "-o", "noatime,compress=zstd,space_cache,discard=async,subvol=@", installDisk+"2", "/mnt")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mkdir", "-p", "/mnt/boot", "/mnt/home", "/mnt/var", "/mnt/.snapshots")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mkdir", "-p", "/mnt/var/log")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mount", "-o", "noatime,compress=zstd,space_cache,discard=async,subvol=@home", installDisk+"2", "/mnt/home")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mount", "-o", "noatime,compress=zstd,space_cache,discard=async,subvol=@snapshots", installDisk+"2", "/mnt/.snapshots")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("mount", "-o", "noatime,compress=zstd,space_cache,discard=async,subvol=@var_log", installDisk+"2", "/mnt/var/log")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nMounting EFI /boot...")
	cmd = exec.Command("mount", installDisk+"1", "/mnt/boot")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nDownloading packages from Pacstrap...")
	cmd = exec.Command("pacstrap", "/mnt", "base base-devel linux linux-firmware intel-ucode git neovim nano btrfs-progs")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err.Error())
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()

	fmt.Println("\nGenerating FSTAB file...")
	cmd = exec.Command("genfstab", "-U", "/mnt", ">>", "/mnt/etc/fstab")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nSymlinking Asia/Kuala Lumpur time to /etc/localtime...")
	cmd = exec.Command("arch-chroot", "/mnt", "ln", "-sf", "/usr/share/zoneinfo/Asia/Kuala_Lumpur", "/etc/localtime")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nSyncing datetime to ntp server")
	cmd = exec.Command("arch-chroot", "/mnt", "timedatectl", "set-ntp", "true")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nSetting hardware clock...")
	cmd = exec.Command("arch-chroot", "/mnt", "hwclock", "--systohc")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nGenerating en_US_utf-8 locale")
	cmd = exec.Command("arch-chroot", "/mnt", "sed", "-i", "'177s/.//'", "/etc/locale.gen")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("arch-chroot", "/mnt", "locale-gen")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("arch-chroot", "/mnt", "echo", "'LANG=en_US.UTF-8'", ">>", "/etc/locale.conf")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nSetting hostname as archlinux")
	cmd = exec.Command("arch-chroot", "/mnt", "echo", "'archlinux'", ">>", "/etc/hostname")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\nConfiguring /etc/hosts")
	cmd = exec.Command("arch-chroot", "/mnt", "echo", "'127.0.0.1    localhost'", ">>", "/etc/hosts")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("arch-chroot", "/mnt", "echo", "'::1    localhost'", ">>", "/etc/hosts")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd = exec.Command("arch-chroot", "/mnt", "echo", "'127.0.1.1    archlinux.localdomain archlinux'", ">>", "/etc/hosts")
	err = cmd.Run()
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

	cmd = exec.Command("arch-chroot", "/mnt", "echo", "root:"+pwd, "|", "chpasswd")
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Installation done. You will now be able to reboot.")
}
