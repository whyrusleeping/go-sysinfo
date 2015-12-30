package sysinfo

import (
	"bytes"
	"fmt"
	"os/exec"
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

func darwinMemInfo() (*MemStats, error) {
	out, err := exec.Command("vm_stat").Output()
	if err != nil {
		return nil, err
	}

	var mem MemStats
	for _, l := range bytes.Split(out, []byte("\n")) {
		parts := strings.Split(string(l))
		if len(parts) != 2 {
		}
	}
}
