package sysinfo

import (
	"fmt"
	"syscall"
)

func init() {
	diskUsageImpl = darwinDiskUsage
}

func darwinDiskUsage(path string) (*DiskStats, error) {
	var stfst syscall.Statfs_t
	err := syscall.Statfs(path, &stfst)
	if err != nil {
		return nil, err
	}

	free := stfst.Bfree * uint64(stfst.Bsize)
	total := stfst.Bavail * uint64(stfst.Bsize)
	return &DiskStats{
		Free:   free,
		Total:  total,
		FsType: fmt.Sprint(stfst.Type),
	}, nil
}
