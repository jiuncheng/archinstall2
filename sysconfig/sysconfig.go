package sysconfig

type SysConfig struct {
	InstallDisk string
	Processor   string
	Arch        string
	FS          string
}

func NewSysConfig() *SysConfig {
	return &SysConfig{}
}
