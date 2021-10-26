package main

import (
	"fmt"
	"log"
	"os/exec"

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

		if res == "us" || res == "" {
			s.cfg.KBLayout = res
			return nil
		}

		err = exec.Command("localectl", "set-keymap", res).Run()
		if err != nil {
			fmt.Print("Keymap is invalid. Press enter to select again : ")
			fmt.Scanln()
			continue
		}
		s.cfg.KBLayout = res
		return nil
	}
}

func (s *Selection) FileSystemSelection() {
	for {
		fmt.Print("\n\n")
		fmt.Println("----FILESYSTEM-----")
		fmt.Println("1. btrfs")
		fmt.Println("2. ext4")
		fmt.Print("Choose the filesystem for main partition : ")

		var res string
		fmt.Scanln(&res)
		if res == "1" {
			s.cfg.FS = "btrfs"
			return
		} else if res == "2" {
			s.cfg.FS = "ext4"
			return
		}

		fmt.Print("Input is invalid. Press enter to select again : ")
		fmt.Scanln()
		continue
	}
}

// func (s *Selection) diskSelection2() {}

// func (s *Selection) diskSelection2() {}

// func (s *Selection) diskSelection2() {}
