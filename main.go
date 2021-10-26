package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jiuncheng/archinstall2/cmd"
	"github.com/jiuncheng/archinstall2/diskutil"
	"github.com/jiuncheng/archinstall2/filesystem"
	"github.com/jiuncheng/archinstall2/sysconfig"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// Prerequisite
	err = cmd.NewCmd("timedatectl set-ntp true").SetDesc("Syncing datetime to ntp server...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}
	cmd.NewCmd("umount -a").SetDesc("Umounting all drives").Run()

	cfg := sysconfig.NewSysConfig()
	cfg.Package.PacstrapPkg = viper.GetStringSlice("pacstrap_pkg")
	cfg.Package.IntelCPUPkg = viper.GetStringSlice("intel_cpu_pkg")
	cfg.Package.AmdCPUPkg = viper.GetStringSlice("amd_cpu_pkg")
	cfg.Package.NvidiaGPUPkg = viper.GetStringSlice("nvidia_gpu_pkg")
	cfg.Package.AmdGPUPkg = viper.GetStringSlice("amd_gpu_pkg")

	err = NewSelection(cfg).PerformSelection()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = diskutil.NewDiskUtil(cfg).CreateBTRFS()
	if err != nil {
		log.Fatalln(err.Error())
	}
	filesystem.NewBtrfsHelper(cfg).GenerateBTRFSSystem()

	var cpuArgs string
	if cfg.Processor == "intel" {
		cpuArgs = strings.Join(cfg.Package.IntelCPUPkg, " ")
	} else {
		cpuArgs = strings.Join(cfg.Package.AmdCPUPkg, " ")
	}

	var gpuArgs string
	if cfg.GPU == "nvidia" {
		gpuArgs = strings.Join(cfg.Package.NvidiaGPUPkg, " ")
	} else if cfg.GPU == "amd" {
		gpuArgs = strings.Join(cfg.Package.AmdGPUPkg, " ")
	}
	cmd2 := cmd.NewCmd("pacstrap /mnt " + strings.Join(cfg.Package.PacstrapPkg, " ") + cpuArgs + gpuArgs)
	err = cmd2.SetDesc("Downloading packages from Pacstrap...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/bash -c \"genfstab -U /mnt >> /mnt/etc/fstab\"").SetDesc("Generating FSTAB file...").Run()
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

	err = ioutil.WriteFile("/mnt/etc/locale.conf", []byte("LANG=en_US.UTF-8\n"), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = ioutil.WriteFile("/mnt/etc/hostname", []byte(cfg.Hostname+"\n"), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = ioutil.WriteFile("/mnt/etc/hosts", []byte("127.0.0.1\tlocalhost\n"), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}

	file, err := os.OpenFile("/mnt/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	_, err = file.WriteString("::1\t\tlocalhost\n")
	if err != nil {
		log.Fatal(err)
	}

	file2, err := os.OpenFile("/mnt/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file2.Close()
	_, err = file2.WriteString("127.0.1.1\t" + cfg.Hostname + ".localdomain\t" + cfg.Hostname + "\n")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Print("\nPlease enter a root password: ")
	// var pwd string
	// _, err = fmt.Scanln(&pwd)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// fmt.Println("New root password is ", pwd)

	err = cmd.NewCmd("/bin/bash -c \"arch-chroot /mnt echo root:" + cfg.RootPassword + " | " + "chpasswd\"").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, superuser := range cfg.Superusers {
		err = cmd.NewCmd("arch-chroot /mnt useradd -mG wheel -s /bin/bash -p " + superuser.Password + " " + superuser.Username).Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	for _, user := range cfg.Users {
		err = cmd.NewCmd("arch-chroot /mnt useradd -m -s /bin/bash -p " + user.Password + " " + user.Username).Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	// fmt.Println(cfg.InstallDisk)
	fmt.Println("Installation done. You will now be able to reboot.")
}
