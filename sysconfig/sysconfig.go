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
	Desktop      string
	Superusers   []User
	Users        []User
	Package      Package
	Services     []string
}

type Package struct {
	PacstrapPkg  []string
	ExtraPkg     []string
	IntelCPUPkg  []string
	AmdCPUPkg    []string
	NvidiaGPUPkg []string
	AmdGPUPkg    []string
	GnomePkg     []string
	PlasmaPkg    []string
	GrubPkg      []string
}

type User struct {
	Username string
	Password string
	Script   string
}

func NewSysConfig() *SysConfig {
	return &SysConfig{}
}
