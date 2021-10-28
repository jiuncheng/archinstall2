package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jiuncheng/archinstall2/sysconfig"
	"github.com/jiuncheng/archinstall2/utils"
)

type Selection struct {
	cfg *sysconfig.SysConfig
}

type Profile struct {
	Name string
	Desc string
}

func NewSelection(cfg *sysconfig.SysConfig) *Selection {
	return &Selection{cfg: cfg}
}

func (s *Selection) PerformSelection() error {
	s.DiskSelection()
	err := s.LayoutSelection()
	if err != nil {
		return err
	}
	err = s.TimezoneSelection()
	if err != nil {
		return err
	}

	s.FileSystemSelection()
	s.BootloaderSelection()
	s.HostnameSelection()
	s.RootPassSelection()
	s.SuperUserSelection()
	s.OptionalUserSelection()
	s.ProcessorSelection()
	s.GPUSelection()
	s.ProfileSelection()

	return nil
}

func (s *Selection) DiskSelection() {
	dl, err := utils.GetDiskList()
	if err != nil {
		log.Fatalln(err.Error())
	}

	var result int
	for {
		for i, disk := range dl {
			fmt.Printf("%d.    /dev/%s    %s\n", i+1, disk.Name, disk.Size)
		}
		fmt.Print("\n Please select the number which the disk will be used for installation (e.g. 1): ")

		_, err := fmt.Scanf("%d", &result)
		if err == nil {
			if result <= len(dl) && result > 0 {
				break
			}
			fmt.Println("\nThe number must be between 1 and ", len(dl), ".")
			fmt.Print("Press enter to choose again : ")
			fmt.Scanln()
			continue
		}
		fmt.Println("\nOnly number between 1 and ", len(dl), " is allowed.")
		fmt.Print("Press enter to choose again : ")
		fmt.Scanln()
		continue
	}

	s.cfg.InstallDisk = "/dev/" + dl[result-1].Name
	fmt.Println(s.cfg.InstallDisk)
}

func (s *Selection) LayoutSelection() error {
	for {
		out, err := exec.Command("localectl", "list-keymaps").Output()
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		fmt.Print("Select one of the following keyboard layout or skip (default: us) : ")
		var res string
		fmt.Scanln(&res)

		if strings.TrimSpace(res) == "us" || strings.TrimSpace(res) == "" {
			s.cfg.KBLayout = strings.TrimSpace("us")
		} else {
			s.cfg.KBLayout = strings.TrimSpace(res)
		}

		err = exec.Command("localectl", "set-keymap", s.cfg.KBLayout).Run()
		if err != nil {
			fmt.Print("Keymap is invalid. Press enter to select again : ")
			fmt.Scanln()
			continue
		}
		fmt.Println(s.cfg.KBLayout)
		return nil
	}
}

func (s *Selection) TimezoneSelection() error {
	for {
		fmt.Print("\n\n")
		fmt.Println("Select and choose your timezone. Enter ? to search for available timezones.")
		fmt.Print("Enter your choice (default: Asia/Kuala_Lumpur) : ")
		var res string
		fmt.Scanln(&res)
		res = strings.TrimSpace(res)
		if res == "?" {
			var searchString string
			fmt.Print("\nType here to search for timezones : ")
			fmt.Scanln(&searchString)
			out, err := exec.Command("sh", "-c", "timedatectl list-timezones | grep -i "+searchString).Output()
			if err != nil {
				continue
			}
			fmt.Println(string(out))
			continue
		} else if res == "" {
			s.cfg.Timezone = "Asia/Kuala_Lumpur"
			fmt.Println(s.cfg.Timezone)
			break
		} else {
			info, err := os.Stat("/usr/share/zoneinfo/" + res)
			if os.IsNotExist(err) {
				fmt.Print("The timezone entered is not valid. Press enter to try again.")
				fmt.Scanln()
				continue
			}

			if info.IsDir() {
				fmt.Print("The timezone entered is not valid. Press enter to try again.")
				fmt.Scanln()
				continue
			}

			// err := exec.Command("sh", "-c", "timedatectl list-timezones | grep "+res).Run()
			// if err != nil {
			// 	fmt.Print("The timezone entered is not valid. Press enter to try again.")
			// 	fmt.Scanln()
			// 	continue
			// }
			s.cfg.Timezone = res
			fmt.Println(s.cfg.Timezone)
			break
		}
	}

	return nil
}

func (s *Selection) FileSystemSelection() {
	for {
		fmt.Print("\n\n")
		fmt.Println("----FILESYSTEM-----")
		fmt.Println("1. btrfs")
		fmt.Println("2. ext4")
		fmt.Print("Choose the filesystem for main partition (default: btrfs) : ")

		var res string
		fmt.Scanln(&res)
		if strings.TrimSpace(res) == "1" || strings.TrimSpace(res) == "" {
			s.cfg.FS = "btrfs"
			fmt.Println(s.cfg.FS)
			return
		} else if strings.TrimSpace(res) == "2" {
			s.cfg.FS = "ext4"
			fmt.Println(s.cfg.FS)
			return
		}

		fmt.Print("Input is invalid. Press enter to select again : ")
		fmt.Scanln()
		continue
	}
}

// func (s *Selection) EncryptionSelection() {
// }

func (s *Selection) BootloaderSelection() {
	for {
		fmt.Print("\n\n")
		fmt.Println("-----Bootloader-----")
		fmt.Println("1. systemd-boot")
		fmt.Println("2. grub")
		fmt.Print("Select desired bootloader : ")

		var res string
		fmt.Scanln(&res)
		if strings.TrimSpace(res) == "1" {
			s.cfg.BootLoader = "systemd-boot"
			fmt.Println(s.cfg.BootLoader)
			break
		} else if strings.TrimSpace(res) == "2" {
			s.cfg.BootLoader = "grub"
			fmt.Println(s.cfg.BootLoader)
			break
		}
		continue
	}
}

func (s *Selection) HostnameSelection() {
	fmt.Print("\n\n")
	fmt.Print("Enter desired hostname for installation (default: archlinux) : ")

	var res string
	fmt.Scanln(&res)
	if strings.TrimSpace(res) == "" {
		s.cfg.Hostname = "archlinux"
		fmt.Println(s.cfg.Hostname)
		return
	}
	s.cfg.Hostname = strings.TrimSpace(res)
	fmt.Println(s.cfg.Hostname)
}

func (s *Selection) RootPassSelection() {
	for {
		fmt.Print("\n")
		fmt.Print("Enter password for root user : ")
		var password string
		fmt.Scanln(&password)
		if strings.TrimSpace(password) == "" {
			fmt.Println("Password cannot be empty.")
			continue
		}

		fmt.Print("Enter password again for confirmation : ")
		var password2 string
		fmt.Scanln(&password2)
		if strings.TrimSpace(password2) != strings.TrimSpace(password) {
			fmt.Println("Password does not match. Try again.")
			continue
		}
		s.cfg.RootPassword = strings.TrimSpace(password)
		break
	}
}

func (s *Selection) SuperUserSelection() {
	var user sysconfig.User

	for {
		fmt.Print("\n\n")
		fmt.Print("Enter username for superuser : ")
		var username string
		fmt.Scanln(&username)
		if strings.TrimSpace(username) == "" {
			fmt.Println("Username cannot be empty.")
			continue
		}
		user.Username = strings.TrimSpace(username)
		break
	}

	for {
		fmt.Print("\n")
		fmt.Print("Enter password for superuser : ")
		var password string
		fmt.Scanln(&password)
		if strings.TrimSpace(password) == "" {
			fmt.Println("Password cannot be empty.")
			continue
		}

		fmt.Print("Enter password again for confirmation : ")
		var password2 string
		fmt.Scanln(&password2)
		if strings.TrimSpace(password2) != strings.TrimSpace(password) {
			fmt.Println("Password does not match. Try again.")
			continue
		}
		user.Password = strings.TrimSpace(password)
		break
	}
	s.cfg.Superusers = append(s.cfg.Superusers, user)
	fmt.Println(strings.TrimSpace(user.Username))
}

func (s *Selection) OptionalUserSelection() {
	for {
		for {
			fmt.Print("\n")
			fmt.Print("Do you want to create additional user? [y/N] : ")
			var res string
			fmt.Scanln(&res)
			if strings.TrimSpace(res) == "y" || strings.TrimSpace(res) == "Y" {
				break
			} else if strings.TrimSpace(res) == "n" || strings.TrimSpace(res) == "N" || strings.TrimSpace(res) == "" {
				return
			} else {
				continue
			}
		}

		var user sysconfig.User
		for {
			fmt.Print("\n\n")
			fmt.Print("Enter username for user : ")
			var username string
			fmt.Scanln(&username)
			if strings.TrimSpace(username) == "" {
				fmt.Println("Username cannot be empty.")
				continue
			}
			user.Username = strings.TrimSpace(username)
			fmt.Println(strings.TrimSpace(user.Username))
			break
		}

		for {
			fmt.Print("\n")
			fmt.Print("Enter password for user : ")
			var password string
			fmt.Scanln(&password)
			if strings.TrimSpace(password) == "" {
				fmt.Println("Password cannot be empty.")
				continue
			}

			fmt.Print("Enter password again for confirmation : ")
			var password2 string
			fmt.Scanln(&password2)
			if strings.TrimSpace(password2) != strings.TrimSpace(password) {
				fmt.Println("Password does not match. Try again.")
				continue
			}
			user.Password = strings.TrimSpace(password)
			break
		}

		for {
			fmt.Print("\n")
			fmt.Print("Set this user as superuser (sudoer)? [y/N] : ")
			var res string
			fmt.Scanln(&res)
			if strings.TrimSpace(res) == "y" || strings.TrimSpace(res) == "Y" {
				s.cfg.Superusers = append(s.cfg.Superusers, user)
				break
			} else if strings.TrimSpace(res) == "n" || strings.TrimSpace(res) == "N" || strings.TrimSpace(res) == "" {
				s.cfg.Users = append(s.cfg.Users, user)
				break
			} else {
				continue
			}
		}

	}

}

func (s *Selection) ProcessorSelection() {
	for {
		fmt.Print("\n\n")
		fmt.Println("-----Processor-----")
		fmt.Println("1. Intel cpu")
		fmt.Println("2. AMD cpu")
		fmt.Print("Select processor model : ")

		var res string
		fmt.Scanln(&res)
		if strings.TrimSpace(res) == "1" {
			s.cfg.Processor = "intel"
			fmt.Println(s.cfg.Processor)
			break
		} else if strings.TrimSpace(res) == "2" {
			s.cfg.Processor = "amd"
			fmt.Println(s.cfg.Processor)
			break
		}
		continue
	}
}

func (s *Selection) GPUSelection() {
	for {
		fmt.Print("\n\n")
		fmt.Println("-----Graphics Model-----")
		fmt.Println("1. Intel gpu")
		fmt.Println("2. Nvidia gpu")
		fmt.Println("3. AMD gpu")
		fmt.Print("Select graphics model : ")

		var res string
		fmt.Scanln(&res)
		if strings.TrimSpace(res) == "1" {
			s.cfg.GPU = "intel"
			fmt.Println(s.cfg.GPU)
			break
		} else if strings.TrimSpace(res) == "2" {
			s.cfg.GPU = "nvidia"
			fmt.Println(s.cfg.GPU)
			break
		} else if strings.TrimSpace(res) == "3" {
			s.cfg.GPU = "amd"
			fmt.Println(s.cfg.GPU)
			break
		}
		continue
	}
}

func (s *Selection) ProfileSelection() {
	var list []*Profile
	log.Println(s.cfg.ProfileList)

	for name, desc := range s.cfg.ProfileList {
		log.Println(name, desc)
		newProfile := &Profile{Name: name, Desc: desc}
		list = append(list, newProfile)
	}
	log.Println(list)
	for {
		fmt.Print("\n\n")
		fmt.Println("-----Install Profile-----")
		for i, profile := range list {
			fmt.Printf("%d. %s", i+1, profile.Desc)
		}
		fmt.Print("Select desktop environment : ")

		var res int
		_, err := fmt.Scanln(&res)
		if err != nil {
			continue
		}

		if res > len(list) || res < 1 {
			continue
		}

		s.cfg.Profile = list[res-1].Name
		fmt.Println(s.cfg.Profile)
		break
	}
}
