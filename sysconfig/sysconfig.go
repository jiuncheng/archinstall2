package sysconfig

type SysConfig struct {
	InstallDisk string
	Processor   string
	Arch        string
	Hostname    string
	FS          string
	KBLayout    string
	Language    string
	Timezone    string
	Superusers  []User
	Users       []User
}

type User struct {
	Username string
	Password string
	Script   string
}

func NewSysConfig() *SysConfig {
	return &SysConfig{}
}
