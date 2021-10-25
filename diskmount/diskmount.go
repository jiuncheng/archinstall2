package diskmount

import (
	"fmt"

	"github.com/jiuncheng/archinstall2/cmd"
)

type DiskMount struct {
	options string
}

func NewDiskMount() *DiskMount {
	return &DiskMount{}
}

func (d *DiskMount) Options(options string) *DiskMount {
	d.options = options
	return d
}

func (d *DiskMount) Mount(source string, target string) error {
	var command string
	if d.options != "" {
		command = fmt.Sprintf("mount -o %s %s %s", d.options, source, target)
	} else {
		command = fmt.Sprintf("mount %s %s", source, target)
	}

	err := cmd.NewCmd(command).Run()
	return err
}

func (d *DiskMount) MountWithDesc(source string, target string, desc string) error {
	var command string
	if d.options != "" {
		command = fmt.Sprintf("mount -o %s %s %s", d.options, source, target)
	} else {
		command = fmt.Sprintf("mount %s %s", source, target)
	}

	err := cmd.NewCmd(command).SetDesc(desc).Run()
	return err
}

func (d *DiskMount) Umount(target string) error {
	err := cmd.NewCmd("umount " + target).Run()
	return err
}
