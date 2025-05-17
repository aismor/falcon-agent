package service

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/google/gousb"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/dev/falcon-agent/internal/model"
)

func CollectMachineInfo() (*model.MachineInfo, error) {
	info := &model.MachineInfo{
		USBDevices: make([]model.USBDevice, 0),
	}

	// Sistema operacional
	info.OS = runtime.GOOS

	// Hostname
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	// Memória
	mems, err := ghw.Memory()
	if err == nil {
		if len(mems.Modules) > 0 {
			for _, mod := range mems.Modules {
				info.Memory = append(info.Memory, model.MemoryInfo{
					Slot:         mod.Label,
					SizeMB:       uint64(mod.SizeBytes / (1024 * 1024)),
					Manufacturer: mod.Vendor,
					SerialNumber: mod.SerialNumber,
				})
			}
		} else {
			// Fallback: mostra apenas o total de memória
			info.Memory = append(info.Memory, model.MemoryInfo{
				Slot:         "Total",
				SizeMB:       uint64(mems.TotalPhysicalBytes / (1024 * 1024)),
				Manufacturer: "N/A",
				SerialNumber: "N/A",
			})
		}
	}

	// HDs
	diskInfo, err := ghw.Block()
	if err == nil {
		for _, d := range diskInfo.Disks {
			if d.DriveType == ghw.DRIVE_TYPE_HDD || d.DriveType == ghw.DRIVE_TYPE_SSD {
				sizeGB := uint64(d.SizeBytes / (1024 * 1024 * 1024))
				modeloValido := d.Model != "" && d.Model != "unknown"
				serialValido := d.SerialNumber != "" && d.SerialNumber != "unknown"
				if sizeGB > 10 && (modeloValido || serialValido) {
					info.HDs = append(info.HDs, model.HDInfo{
						Model:  d.Model,
						Serial: d.SerialNumber,
						SizeGB: sizeGB,
					})
				}
			}
		}
	}

	// Processador
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.Processor = model.ProcessorInfo{
			Model:        cpuInfo[0].ModelName,
			Cores:        int(cpuInfo[0].Cores),
			Threads:      len(cpuInfo),
			FrequencyGHz: cpuInfo[0].Mhz / 1000.0,
		}
	}

	// BIOS
	bios, err := ghw.BIOS()
	if err == nil {
		info.BIOS = model.BIOSInfo{
			Vendor:      bios.Vendor,
			Version:     bios.Version,
			ReleaseDate: bios.Date,
		}
	}

	// Placa-mãe (Motherboard Serial)
	baseboard, err := ghw.Baseboard()
	if err == nil {
		info.MotherboardSN = baseboard.SerialNumber
	}

	// Serial Number (geral)
	hostInfo, err := host.Info()
	if err == nil {
		info.SerialNumber = hostInfo.HostID
	}

	// Dispositivos USB
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Configura debug
	ctx.Debug(3)

	// Lista todos os dispositivos
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// Ignora hubs USB
		if desc.Class == gousb.ClassHub {
			return false
		}
		return true
	})

	if err != nil {
		log.Printf("Erro ao listar dispositivos USB: %v", err)
	}

	for _, dev := range devs {
		defer dev.Close()

		// Tenta obter informações do dispositivo
		var name string
		manufacturer, _ := dev.Manufacturer()
		product, _ := dev.Product()
		serial, _ := dev.SerialNumber()

		if manufacturer != "" && product != "" {
			name = fmt.Sprintf("%s %s", manufacturer, product)
		} else {
			name = fmt.Sprintf("USB Device %04x:%04x", dev.Desc.Vendor, dev.Desc.Product)
		}

		info.USBDevices = append(info.USBDevices, model.USBDevice{
			VendorID:  fmt.Sprintf("%04x", dev.Desc.Vendor),
			ProductID: fmt.Sprintf("%04x", dev.Desc.Product),
			Name:      name,
			Serial:    serial,
		})
	}

	return info, nil
}
