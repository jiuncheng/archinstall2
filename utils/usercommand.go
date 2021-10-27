package utils

import "fmt"

type UserCommand struct {
	args string
}

func UserCmd(args string) *UserCommand {
	return &UserCommand{args: args}
}

func (u *UserCommand) Run(user string) error {
	fmt.Print("User:::")
	err := NewCmd("arch-chroot /mnt /usr/bin/runuser -u " + user + " -- " + u.args).Run()
	if err != nil {
		return err
	}
	return nil
}
