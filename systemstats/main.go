package main

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func handleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func roundFloat(number float64, decimalPlaces int) string {
	return strconv.FormatFloat(math.Round(number*math.Pow(10.0, float64(decimalPlaces)))/math.Pow(10.0, float64(decimalPlaces)), 'f', -1, 64)
}

func main() {
	// get uptime information
	uptime, err := host.Uptime()
	handleErr(err)

	// get total and used memory
	memory, err := mem.VirtualMemory()
	handleErr(err)

	// get OS information
	osInfo, err := host.Info()
	handleErr(err)

	// get cpu information
	cpuCoresLogical, err := cpu.Counts(true)
	handleErr(err)

	cpuCoresPhysical, err := cpu.Counts(false)
	handleErr(err)

	cpuInfo, err := cpu.Info()
	handleErr(err)

	// get information about disks
	disks, err := disk.Partitions(false)
	handleErr(err)

	// print information on screen
	// host
	fmt.Printf("Host:\t %s\n", osInfo.Hostname)

	// OS
	fmt.Printf(
		"OS:\t %s %s (%s - %s) \n",
		osInfo.KernelVersion, osInfo.PlatformVersion, osInfo.OS, osInfo.KernelArch,
	)

	// uptime
	fmt.Printf("Uptime:\t %d:%d:%d\n", uptime/60/60, uptime/60%60, uptime%60%60)

	// CPU
	fmt.Printf("CPU:\t %s (%d / %d) @ %s Mhz\n",
		cpuInfo[0].ModelName, cpuCoresPhysical, cpuCoresLogical, strconv.FormatFloat(cpuInfo[0].Mhz, 'f', -1, 64))

	// memory
	fmt.Printf("Memory:\t used : %s\tGb\n\t total: %s\tGb\n",
		roundFloat(float64(memory.Used)/1000000000, 3),
		roundFloat(float64(memory.Total)/1000000000, 3))

	// disks
	fmt.Printf("Disks: ")
	for _, d := range disks {
		device, err := disk.Usage(d.Mountpoint)
		handleErr(err)

		fmt.Printf("\t %s:\t%s Gb / %s Gb (%s%%) - %s\n",
			d.Device,
			roundFloat(float64(device.Used)/1000000000, 3),
			roundFloat(float64(device.Total)/1000000000, 3),
			roundFloat(device.UsedPercent, 2),
			device.Fstype)
	}
}
