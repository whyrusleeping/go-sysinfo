package sysinfo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"syscall"

	humanize "github.com/dustin/go-humanize"
)

func init() {
	diskUsageImpl = linuxDiskUsage
	memInfoImpl = linuxMemInfo
}

func linuxDiskUsage(path string) (*DiskStats, error) {
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

func linuxMemInfo() (*MemStats, error) {
	info, err := ioutil.ReadFile("/proc/self/status")
	if err != nil {
		return nil, err
	}

	var stats MemStats
	for _, e := range bytes.Split(info, []byte("\n")) {
		if !bytes.HasPrefix(e, []byte("Vm")) {
			continue
		}

		parts := bytes.Split(e, []byte(":"))
		if len(parts) != 2 {
			return nil, fmt.Errorf("unexpected line in proc stats: %q", string(e))
		}

		switch string(parts[0]) {
		case "VmSize":
			vmsize, err := humanize.ParseBytes(string(parts[1]))
			if err != nil {
				return nil, err
			}

			stats.Used = vmsize
		case "VmSwap":
			swapsize, err := humanize.ParseBytes(string(parts[1]))
			if err != nil {
				return nil, err
			}

			stats.Swap = swapsize
		}
	}
	return &stats, nil
}
