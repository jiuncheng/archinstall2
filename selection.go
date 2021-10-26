package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/jiuncheng/archinstall2/disklist"
	"github.com/jiuncheng/archinstall2/sysconfig"
)

type Selection struct {
	cfg *sysconfig.SysConfig
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
	s.FileSystemSelection()
	s.HostnameSelection()
	s.SuperUserSelection()
	s.OptionalUserSelection()

	return nil
}

func (s *Selection) DiskSelection() {
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
			s.cfg.KBLayout = strings.TrimSpace(res)
			return nil
		}

		err = exec.Command("localectl", "set-keymap", res).Run()
		if err != nil {
			fmt.Print("Keymap is invalid. Press enter to select again : ")
			fmt.Scanln()
			continue
		}
		s.cfg.KBLayout = strings.TrimSpace(res)
		return nil
	}
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
			return
		} else if strings.TrimSpace(res) == "2" {
			s.cfg.FS = "ext4"
			return
		}

		fmt.Print("Input is invalid. Press enter to select again : ")
		fmt.Scanln()
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
		return
	}
	s.cfg.Hostname = strings.TrimSpace(res)
}

func (s *Selection) RootPassSelection() {
	for {
		fmt.Print("\n\n")
		fmt.Print("Enter root password : ")

		var res string
		fmt.Scanln(&res)
		if strings.TrimSpace(res) == "" {
			fmt.Print("Password cannot be empty. Press enter to select again : ")
			fmt.Scanln()
			continue
		}
		s.cfg.RootPassword = strings.TrimSpace(res)
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
		user.Password = strings.TrimSpace(password2)
		break
	}
	s.cfg.Superusers = append(s.cfg.Superusers, user)
}

func (s *Selection) OptionalUserSelection() {
	for {
		for {
			fmt.Print("\n")
			fmt.Print("Do you want to create a new user? [y/N] : ")
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
			user.Password = strings.TrimSpace(password2)
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

// func (s *Selection) diskSelection2() {}

// func (s *Selection) diskSelection2() {}
