package disklist

import "encoding/json"

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

func NewDiskList() *DiskList {
	return &DiskList{}
}

func NewDiskListFromJSON(input []byte) (*DiskList, error) {
	var newDL DiskList
	err := json.Unmarshal(input, &newDL)
	if err != nil {
		return nil, err
	}
	return &newDL, nil
}
