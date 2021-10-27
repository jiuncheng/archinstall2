package main

import (
	"fmt"

	"github.com/jiuncheng/archinstall2/sysconfig"
	"github.com/jiuncheng/archinstall2/utils"
)

type PostInstall struct {
	cfg  *sysconfig.SysConfig
	user string
}

func NewPostInstall(cfg *sysconfig.SysConfig) *PostInstall {
	return &PostInstall{cfg: cfg, user: cfg.Superusers[0].Username}
}

func (p *PostInstall) PerformPostInstall() error {
	err := p.InstallParu()
	if err != nil {
		return err
	}

	err = p.InstallZramd()
	if err != nil {
		return err
	}

	// err = p.InstallCutefishBeta()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (p *PostInstall) InstallParu() error {
	err := utils.UserCmd(fmt.Sprintf("git clone https://aur.archlinux.org/paru-bin /home/%s/paru-bin", p.user)).Run(p.user)
	if err != nil {
		return err
	}
	err = utils.UserCmd(fmt.Sprintf(`bash -c "cd /home/%s/paru-bin && makepkg -si --noconfirm"`, p.user)).Run(p.user)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostInstall) InstallZramd() error {
	err := utils.UserCmd(`paru -S --noconfirm zramd`).Run(p.user)
	if err != nil {
		return err
	}
	// err = utils.NewCmd(`sed -i "s/^# MAX_SIZE/MAX_SIZE/" /etc/default/zramd`).SetDesc("Enabling parallel download..").Run()
	// if err != nil {
	// 	return err
	// }

	err = utils.NewCmd("arch-chroot /mnt systemctl enable zramd").SetDesc("Enabling zramd service..").Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostInstall) InstallCutefishBeta() error {
	err := utils.UserCmd(`paru -S --noconfirm cutefish-git sddm`).Run(p.user)
	if err != nil {
		return err
	}

	err = utils.NewCmd("arch-chroot /mnt systemctl enable sddm").SetDesc("Enabling sddm service..").Run()
	if err != nil {
		return err
	}

	return nil
}
