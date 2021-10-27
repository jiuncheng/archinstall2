package filesystem

import (
	"github.com/jiuncheng/archinstall2/sysconfig"
	"github.com/jiuncheng/archinstall2/utils"
)

type BtrfsHelper struct {
	cfg *sysconfig.SysConfig
}

func NewBtrfsHelper(cfg *sysconfig.SysConfig) *BtrfsHelper {
	return &BtrfsHelper{cfg: cfg}
}

func (b *BtrfsHelper) CreateSubvolume(subvol string) error {
	return utils.NewCmd("btrfs subvolume create /mnt/" + subvol).Run()
}

func (b *BtrfsHelper) MountMnt(dm *utils.DiskMount) error {
	return dm.MountWithDesc(b.cfg.InstallDisk+"2", "/mnt", "Mounting disk for btrfs subvolume creation...")
}

func (b *BtrfsHelper) CreateRootSubvol() error {
	return b.CreateSubvolume("@")
}

func (b *BtrfsHelper) CreateHomeSubvol() error {
	return b.CreateSubvolume("@home")
}

func (b *BtrfsHelper) CreateSnapshotSubvol() error {
	return b.CreateSubvolume("@snapshots")
}

func (b *BtrfsHelper) CreateVarLogSubvol() error {
	return b.CreateSubvolume("@var_log")
}

func (b *BtrfsHelper) UmountMnt(dm *utils.DiskMount) error {
	return dm.Umount("/mnt")
}

func (b *BtrfsHelper) MountRootSubvol() error {
	dm := utils.NewDiskMount()
	dm.Options("noatime,compress=zstd,space_cache,discard=async,subvol=@")
	return dm.MountWithDesc(b.cfg.InstallDisk+"2", "/mnt", "Remounting with BTRFS options...")
}

func (b *BtrfsHelper) MountHomeSubvol() error {
	dm := utils.NewDiskMount()
	dm.Options("noatime,compress=zstd,space_cache,discard=async,subvol=@home")
	return dm.Mount(b.cfg.InstallDisk+"2", "/mnt/home")
}

func (b *BtrfsHelper) MountSnapshotsSubvol() error {
	dm := utils.NewDiskMount()
	dm.Options("noatime,compress=zstd,space_cache,discard=async,subvol=@snapshots")
	return dm.Mount(b.cfg.InstallDisk+"2", "/mnt/.snapshots")
}

func (b *BtrfsHelper) MountVarLogSubvol() error {
	dm := utils.NewDiskMount()
	dm.Options("noatime,compress=zstd,space_cache,discard=async,subvol=@var_log")
	return dm.Mount(b.cfg.InstallDisk+"2", "/mnt/var/log")
}

func (b *BtrfsHelper) MountEFI() error {
	dm := utils.NewDiskMount()
	return dm.MountWithDesc(b.cfg.InstallDisk+"1", "/mnt/boot", "Mounting EFI /boot...")
}

func (b *BtrfsHelper) GenerateBTRFSSystem() error {
	dm := utils.NewDiskMount()
	b.MountMnt(dm)
	b.CreateRootSubvol()
	b.CreateHomeSubvol()
	b.CreateSnapshotSubvol()
	b.CreateVarLogSubvol()
	b.UmountMnt(dm)

	b.MountRootSubvol()

	err := utils.NewCmd("mkdir -p /mnt/boot /mnt/home /mnt/var /mnt/.snapshots").Run()
	if err != nil {
		return err
	}

	err = utils.NewCmd("mkdir -p /mnt/var/log").Run()
	if err != nil {
		return err
	}

	b.MountHomeSubvol()
	b.MountSnapshotsSubvol()
	b.MountVarLogSubvol()
	b.MountEFI()

	return nil
}
