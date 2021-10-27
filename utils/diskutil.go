package utils

import (
	"github.com/jiuncheng/archinstall2/sysconfig"
)

type DiskUtil struct {
	cfg *sysconfig.SysConfig
}

func NewDiskUtil(cfg *sysconfig.SysConfig) *DiskUtil {
	return &DiskUtil{cfg: cfg}
}

func (d *DiskUtil) FormatDiskToGPT() error {
	return NewCmd("sgdisk -o " + d.cfg.InstallDisk).SetDesc("Formatting selected disk...").Run()
}

func (d *DiskUtil) NewEFIPartition() error {
	return NewCmd("sgdisk -n 1:0:+500M -t 1:ef00 -c 1:'EFI' " + d.cfg.InstallDisk).Run()
}

func (d *DiskUtil) NewLinuxFSPartition(size string) error {
	return NewCmd("sgdisk -n 2:0:" + size + " -t 2:8300 -c 2:'rootfs' " + d.cfg.InstallDisk).Run()
}

func (d *DiskUtil) MkfsVfatEFI() error {
	return NewCmd("mkfs.vfat -n EFI " + d.cfg.InstallDisk + "1").Run()
}

func (d *DiskUtil) MkfsBTRFS() error {
	return NewCmd("mkfs.btrfs -f -L ROOT " + d.cfg.InstallDisk + "2").Run()
}

func (d *DiskUtil) CreateBTRFS() error {
	err := d.FormatDiskToGPT()
	if err != nil {
		return err
	}

	err = d.NewEFIPartition()
	if err != nil {
		return err
	}

	err = d.NewLinuxFSPartition("-4G")
	if err != nil {
		return err
	}

	err = d.MkfsVfatEFI()
	if err != nil {
		return err
	}

	err = d.MkfsBTRFS()
	if err != nil {
		return err
	}

	return nil
}
