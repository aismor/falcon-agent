package model

// MachineInfo representa as informações coletadas da máquina
type MachineInfo struct {
	OS            string
	Hostname      string
	Processor     ProcessorInfo
	BIOS          BIOSInfo
	MotherboardSN string
	SerialNumber  string
	Memory        []MemoryInfo
	HDs           []HDInfo
	USBDevices    []USBDevice
	Temperatures  []TemperatureInfo
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

// TemperatureInfo representa as informações de temperatura
type TemperatureInfo struct {
	Sensor string
	ValueC float64
}
