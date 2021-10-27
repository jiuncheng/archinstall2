package utils

import (
	"encoding/json"
)

type DiskList struct {
	BlockDevices []Disk `json:"blockdevices"`
}

type Disk struct {
	Name        string   `json:"name"`
	MajMin      string   `json:"maj:min"`
	RM          bool     `json:"rm"`
	Size        string   `json:"size"`
	RO          bool     `json:"ro"`
	Type        string   `json:"type"`
	MountPoints []string `json:"mountpoints"`
}

func NewDiskList() (*DiskList, error) {
	output, err := NewCmd("lsblk -ldnJe 7,11").Output()
	if err != nil {
		return nil, err
	}

	var newDL DiskList
	err = json.Unmarshal(output, &newDL)
	if err != nil {
		return nil, err
	}
	return &newDL, nil
}

func GetDiskList() ([]Disk, error) {
	dl, err := NewDiskList()
	if err != nil {
		return nil, err
	}

	return dl.BlockDevices, nil
}
