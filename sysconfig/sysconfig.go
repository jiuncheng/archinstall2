package sysconfig

type SysConfig struct {
	RootPassword string
	LUKSPassword string
	InstallDisk  string
	Processor    string
	Arch         string
	Hostname     string
	FS           string
	KBLayout     string
	Language     string
	Timezone     string
	BootLoader   string
	Superusers   []User
	Users        []User
}

type User struct {
	Username string
	Password string
	Script   string
}

func NewSysConfig() *SysConfig {
	return &SysConfig{}
}
