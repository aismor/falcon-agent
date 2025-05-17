package model

// MachineInfo representa as informações coletadas da máquina
type MachineInfo struct {
	OS            string
	Hostname      string
	Processor     ProcessorInfo
	BIOS          BIOSInfo
	Memory        []MemoryInfo
	HDs           []HDInfo
	USBDevices    []USBDevice
	MotherboardSN string
	SerialNumber  string
}

// ProcessorInfo representa as informações do processador
type ProcessorInfo struct {
	Model        string
	Cores        int
	Threads      int
	FrequencyGHz float64
}

// BIOSInfo representa as informações da BIOS
type BIOSInfo struct {
	Vendor      string
	Version     string
	ReleaseDate string
}

// MemoryInfo representa as informações de memória
type MemoryInfo struct {
	Slot         string
	SizeMB       uint64
	Manufacturer string
	SerialNumber string
}

// HDInfo representa as informações do disco rígido
type HDInfo struct {
	Model  string
	Serial string
	SizeGB uint64
}

// USBDevice representa as informações de um dispositivo USB
type USBDevice struct {
	VendorID  string
	ProductID string
	Name      string
	Serial    string
}
