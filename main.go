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
	globalConf := viper.New()
	globalConf.SetConfigName("config")
	globalConf.SetConfigType("yaml")
	globalConf.AddConfigPath(".")
	err := globalConf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	cfg := sysconfig.NewSysConfig()
	cfg.Package.PacstrapPkg = globalConf.GetStringSlice("pacstrap_pkg")
	cfg.Package.GrubPkg = globalConf.GetStringSlice("grub_pkg")
	cfg.Package.IntelCPUPkg = globalConf.GetStringSlice("intel_cpu_pkg")
	cfg.Package.AmdCPUPkg = globalConf.GetStringSlice("amd_cpu_pkg")
	cfg.Package.NvidiaGPUPkg = globalConf.GetStringSlice("nvidia_gpu_pkg")
	cfg.Package.AmdGPUPkg = globalConf.GetStringSlice("amd_gpu_pkg")
	cfg.Package.ExtraPkg = globalConf.GetStringSlice("extra_pkg")

	err = RunPrerequisites()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = NewSelection(cfg).PerformSelection()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if cfg.Desktop == "gnome" {
		desktopConf := viper.New()
		desktopConf.SetConfigName("gnome")
		desktopConf.SetConfigType("yaml")
		desktopConf.AddConfigPath("./desktop")
		err = desktopConf.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		cfg.Package.DesktopPkg = desktopConf.GetStringSlice("desktop_pkg")
		cfg.Services = desktopConf.GetStringSlice("services")
	} else if cfg.Desktop == "plasma" {
		desktopConf := viper.New()
		desktopConf.SetConfigName("plasma")
		desktopConf.SetConfigType("yaml")
		desktopConf.AddConfigPath("./desktop")
		err = desktopConf.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		cfg.Package.DesktopPkg = desktopConf.GetStringSlice("desktop_pkg")
		cfg.Services = desktopConf.GetStringSlice("services")
	} else if cfg.Desktop == "cutefish" {
		desktopConf := viper.New()
		desktopConf.SetConfigName("cutefish")
		desktopConf.SetConfigType("yaml")
		desktopConf.AddConfigPath("./desktop")
		err = desktopConf.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		cfg.Package.DesktopPkg = desktopConf.GetStringSlice("desktop_pkg")
		cfg.Services = desktopConf.GetStringSlice("services")
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

	guiArgs := strings.Join(cfg.Package.DesktopPkg, " ")

	var bootloaderArgs string
	if cfg.BootLoader == "grub" {
		bootloaderArgs = strings.Join(cfg.Package.GrubPkg, " ")
	}

	cmd2 := cmd.NewCmd("pacstrap /mnt " + strings.Join(cfg.Package.PacstrapPkg, " ") + " " + strings.Join(cfg.Package.ExtraPkg, " ") + " " + cpuArgs + " " + gpuArgs + " " + guiArgs + " " + bootloaderArgs)
	err = cmd2.SetDesc("Downloading packages from Pacstrap...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("/bin/bash -c \"genfstab -U /mnt >> /mnt/etc/fstab\"").SetDesc("Generating FSTAB file...").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = cmd.NewCmd("arch-chroot /mnt ln -sf /usr/share/zoneinfo/" + cfg.Timezone + " /etc/localtime").SetDesc("Symlinking Asia/Kuala Lumpur time to /etc/localtime...").Run()
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

	// err = cmd.NewCmd("/bin/bash -c \"arch-chroot /mnt echo root:" + cfg.RootPassword + " | " + "chpasswd\"").Run()
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	err = cmd.NewCmd(fmt.Sprintf(`arch-chroot /mnt sh -c "echo root:%s | chpasswd"`, cfg.RootPassword)).Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, superuser := range cfg.Superusers {
		err = cmd.NewCmd("arch-chroot /mnt useradd -mG wheel -s /bin/fish " + superuser.Username).Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = cmd.NewCmd(fmt.Sprintf(`arch-chroot /mnt sh -c "echo %s:%s | chpasswd"`, superuser.Username, superuser.Password)).Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	for _, user := range cfg.Users {
		err = cmd.NewCmd("arch-chroot /mnt useradd -m -s /bin/fish " + user.Username).Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = cmd.NewCmd(fmt.Sprintf(`arch-chroot /mnt sh -c "echo %s:%s | chpasswd"`, user.Username, user.Password)).Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	if cfg.BootLoader == "grub" {
		err = cmd.NewCmd("arch-chroot /mnt grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=GRUB --removable").SetDesc("Creating grub bootloader...").Run()
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = cmd.NewCmd("arch-chroot /mnt grub-mkconfig -o /boot/grub/grub.cfg").SetDesc("Creating grub config...").Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else if cfg.BootLoader == "systemd-boot" {
		err = cmd.NewCmd("arch-chroot /mnt bootctl --path=/boot install").SetDesc("Creating bootloader...").Run()
		if err != nil {
			log.Fatalln(err.Error())
		}

		loaderContent := "default\t\tarch.conf\n# timeout\t\t4\nconsole-mode\tmax\neditor\t\t\tno\n"
		err = ioutil.WriteFile("/mnt/boot/loader/loader.conf", []byte(loaderContent), 0644)
		if err != nil {
			log.Fatalln(err.Error())
		}

		uuid, err := cmd.NewCmd("findmnt -fn -o UUID " + cfg.InstallDisk + "2").SetDesc("Finding partition UUID...").Output()
		if err != nil {
			log.Fatalln(err.Error())
		}
		archContent := "title Arch Linux\nlinux /vmlinuz-linux\ninitrd /" + cfg.Processor + "-ucode.img\ninitrd /initramfs-linux.img\noptions root=UUID=" + strings.Trim(string(uuid), "\n") + " rootflags=\"subvol=@\" rw\n"
		err = ioutil.WriteFile("/mnt/boot/loader/entries/arch.conf", []byte(archContent), 0644)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	wheelContent := "%wheel\tALL=(ALL)\tALL\n"
	fmt.Println("Writing sudoers.d file to enable wheel group...")
	err = ioutil.WriteFile("/mnt/etc/sudoers.d/wheel", []byte(wheelContent), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = EnableServices(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Installation done. You will now be able to reboot.")
}

func EnableServices(cfg *sysconfig.SysConfig) error {
	for _, service := range cfg.Services {
		err := cmd.NewCmd("arch-chroot /mnt systemctl enable " + service).SetDesc("Enabling " + service + " service...").Run()
		if err != nil {
			return err
		}
	}
	return nil
}
