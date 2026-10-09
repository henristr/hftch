package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/jaypipes/ghw"
	cpu2 "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

const (
	colorReset = "\033[0m"
	colorLogo  = "\033[38;2;45;226;230m"
	colorName  = "\033[38;2;168;182;201m"
	colorInfo  = "\033[38;2;186;104;200m"
)

func formatUptime(d time.Duration) string {
	days := int(d / (24 * time.Hour))
	d %= 24 * time.Hour

	hours := int(d / time.Hour)
	d %= time.Hour

	minutes := int(d / time.Minute)
	d %= time.Minute

	seconds := int(d / time.Second)

	return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
}

func main() {
	username, err := user.Current()
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	systemString := fmt.Sprintf("%s@%s", username.Username, hostname)
	sysStringLen := len(systemString)

	var line []string
	for range sysStringLen {
		line = append(line, "-")
	}

	lineStr := strings.Join(line, "")

	osName, err := host.Info()
	if err != nil {
		panic(err)
	}

	ram, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}

	ramStr := fmt.Sprintf("%.1d MB / %.1d MB", ram.Used/1024/1024, ram.Total/1024/1024)

	cpu, err := cpu2.Info()
	if err != nil {
		panic(err)
	}

	gpu, err := ghw.GPU()
	if err != nil {
		panic(err)
	}

	internalDisk, err := disk.Usage("/")
	if err != nil {
		panic(err)
	}

	diskStr := fmt.Sprintf("%.1d GB / %.1d GB (%.1f%%)", internalDisk.Used/1024/1024/1024, internalDisk.Total/1024/1024/1024, internalDisk.UsedPercent)

	fmt.Printf("%s%-48s%s%s\n", colorLogo, " ___  ___  ________  _________  ________", colorInfo, systemString)
	fmt.Printf("%s%-48s%s%s\n", colorLogo, "|\\  \\|\\  \\|\\   ____\\|\\___   ___\\\\   __  \\", colorInfo, lineStr)
	fmt.Printf("%s%-48s%s%-10s%s%s\n", colorLogo, "\\ \\  \\\\\\  \\ \\  \\___|\\|___ \\  \\_\\ \\  \\|\\  \\", colorName, "OS:", colorInfo, osName.Platform)
	fmt.Printf("%s%-48s%s%-10s%s%s\n", colorLogo, " \\ \\   __  \\ \\_____  \\   \\ \\  \\ \\ \\   _  _\\", colorName, "Kernel:", colorInfo, osName.KernelVersion)
	fmt.Printf("%s%-48s%s%-10s%s%s\n", colorLogo, "  \\ \\  \\ \\  \\|____|\\  \\   \\ \\  \\ \\ \\  \\\\  \\|", colorName, "Arch:", colorInfo, osName.KernelArch)
	fmt.Printf("%s%-48s%s%-10s%s%s\n", colorLogo, "   \\ \\__\\ \\__\\____\\_\\  \\   \\ \\__\\ \\ \\__\\\\ _\\", colorName, "Uptime:", colorInfo, formatUptime(time.Duration(osName.Uptime)*time.Second))
	fmt.Printf("%s%-48s%s%-10s%s%s\n", colorLogo, "    \\|__|\\|__|\\_________\\   \\|__|  \\|__|\\|__|", colorName, "RAM:", colorInfo, ramStr)
	fmt.Printf("%s%-48s%s%-10s%s%s\n", colorLogo, "             \\|_________|", colorName, "CPU:", colorInfo, cpu[0].ModelName)

	for _, g := range gpu.GraphicsCards {
		fmt.Printf("%-48s%s%-10s%s%s\n", "", colorName, "GPU:", colorInfo, g.DeviceInfo.Product.Name)
	}

	fmt.Printf("%-48s%s%-10s%s%s%s\n", "", colorName, "Disk:", colorInfo, diskStr, colorReset)
}
