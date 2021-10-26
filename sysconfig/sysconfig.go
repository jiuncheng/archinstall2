package sysconfig

type SysConfig struct {
	RootPassword string
	LUKSPassword string
	InstallDisk  string
	Processor    string
	GPU          string
	Arch         string
	Hostname     string
	FS           string
	KBLayout     string
	Language     string
	Timezone     string
	BootLoader   string
	Superusers   []User
	Users        []User
	Package      Package
}

type Package struct {
	PacstrapPkg  []string
	IntelCPUPkg  []string
	AmdCPUPkg    []string
	NvidiaGPUPkg []string
	AmdGPUPkg    []string
}

type User struct {
	Username string
	Password string
	Script   string
}

func NewSysConfig() *SysConfig {
	return &SysConfig{}
}
