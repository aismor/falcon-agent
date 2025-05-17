package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dev/falcon-agent/internal/model"
)

type App struct {
	window      fyne.Window
	machineInfo *model.MachineInfo
	updateChan  chan *model.MachineInfo
	content     *fyne.Container
}

func New(machineInfo *model.MachineInfo) *App {
	a := app.New()
	window := a.NewWindow("Falcon Agent")
	window.Resize(fyne.NewSize(1024, 768))
	window.SetMaster()

	updateChan := make(chan *model.MachineInfo, 1)
	updateChan <- machineInfo

	app := &App{
		window:      window,
		machineInfo: machineInfo,
		updateChan:  updateChan,
	}

	app.setupUI()
	go app.updateLoop()

	return app
}

func (a *App) Run() {
	a.window.ShowAndRun()
}

func (a *App) setupUI() {
	// Menu lateral
	sidebar := a.createSidebar()

	// Separador vertical
	separator := canvas.NewLine(theme.ShadowColor())
	separator.StrokeWidth = 2

	// Conteúdo principal inicial (Sistema)
	a.content = a.createSystemContent()

	// Layout principal com split
	split := container.NewHSplit(
		container.NewHBox(sidebar, separator),
		container.NewPadded(a.content),
	)
	split.SetOffset(0.2) // 20% para o menu lateral

	a.window.SetContent(split)
}

func (a *App) createSidebar() *fyne.Container {
	// Botões do menu com ícones
	systemBtn := widget.NewButtonWithIcon("Sistema", theme.ComputerIcon(), func() {
		a.content.Objects = []fyne.CanvasObject{a.createSystemContent()}
		a.content.Refresh()
	})

	cpuBtn := widget.NewButtonWithIcon("CPU", theme.MediaPlayIcon(), func() {
		a.content.Objects = []fyne.CanvasObject{a.createProcessorContent()}
		a.content.Refresh()
	})

	memoryBtn := widget.NewButtonWithIcon("Memória", theme.StorageIcon(), func() {
		a.content.Objects = []fyne.CanvasObject{a.createMemoryContent()}
		a.content.Refresh()
	})

	storageBtn := widget.NewButtonWithIcon("Armazenamento", theme.FolderIcon(), func() {
		a.content.Objects = []fyne.CanvasObject{a.createStorageContent()}
		a.content.Refresh()
	})

	biosBtn := widget.NewButtonWithIcon("BIOS", theme.SettingsIcon(), func() {
		a.content.Objects = []fyne.CanvasObject{a.createBIOSContent()}
		a.content.Refresh()
	})

	usbBtn := widget.NewButtonWithIcon("USB", theme.MediaRecordIcon(), func() {
		a.content.Objects = []fyne.CanvasObject{a.createUSBContent()}
		a.content.Refresh()
	})

	// Estiliza os botões
	for _, btn := range []*widget.Button{systemBtn, cpuBtn, memoryBtn, storageBtn, biosBtn, usbBtn} {
		btn.Alignment = widget.ButtonAlignLeading
		btn.Importance = widget.LowImportance
	}

	// Container do menu
	return container.NewVBox(
		systemBtn,
		cpuBtn,
		memoryBtn,
		storageBtn,
		biosBtn,
		usbBtn,
		layout.NewSpacer(),
	)
}

func (a *App) createSystemContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações do Sistema",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	titleContainer := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
	)

	info := container.NewVBox(
		widget.NewLabelWithStyle("Sistema Operacional", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(a.machineInfo.OS),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Hostname", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(a.machineInfo.Hostname),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Serial Number", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(a.machineInfo.SerialNumber),
	)

	return container.NewVBox(
		titleContainer,
		container.NewPadded(info),
	)
}

func (a *App) createProcessorContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações do Processador",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	titleContainer := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
	)

	info := container.NewVBox(
		widget.NewLabelWithStyle("Modelo", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(a.machineInfo.Processor.Model),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Núcleos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("%d", a.machineInfo.Processor.Cores)),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Threads", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("%d", a.machineInfo.Processor.Threads)),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Frequência", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("%.2f GHz", a.machineInfo.Processor.FrequencyGHz)),
	)

	return container.NewVBox(
		titleContainer,
		container.NewPadded(info),
	)
}

func (a *App) createMemoryContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações da Memória",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	titleContainer := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
	)

	info := container.NewVBox()
	for _, mem := range a.machineInfo.Memory {
		memBox := container.NewVBox(
			widget.NewLabelWithStyle("Slot", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(mem.Slot),
			widget.NewLabelWithStyle("Capacidade", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fmt.Sprintf("%d MB", mem.SizeMB)),
			widget.NewLabelWithStyle("Fabricante", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(mem.Manufacturer),
		)
		if mem.SerialNumber != "N/A" {
			memBox.Add(widget.NewLabelWithStyle("Serial", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			memBox.Add(widget.NewLabel(mem.SerialNumber))
		}
		memBox.Add(widget.NewSeparator())
		info.Add(memBox)
	}

	scroll := container.NewVScroll(container.NewPadded(info))
	scroll.SetMinSize(fyne.NewSize(600, 400))

	return container.NewVBox(
		titleContainer,
		scroll,
	)
}

func (a *App) createStorageContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações de Armazenamento",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	titleContainer := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
	)

	info := container.NewVBox()
	for _, hd := range a.machineInfo.HDs {
		hdBox := container.NewVBox(
			widget.NewLabelWithStyle("Modelo", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(hd.Model),
			widget.NewLabelWithStyle("Capacidade", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fmt.Sprintf("%d GB", hd.SizeGB)),
		)
		if hd.Serial != "" && hd.Serial != "unknown" {
			hdBox.Add(widget.NewLabelWithStyle("Serial", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			hdBox.Add(widget.NewLabel(hd.Serial))
		}
		hdBox.Add(widget.NewSeparator())
		info.Add(hdBox)
	}

	scroll := container.NewVScroll(container.NewPadded(info))
	scroll.SetMinSize(fyne.NewSize(600, 400))

	return container.NewVBox(
		titleContainer,
		scroll,
	)
}

func (a *App) createBIOSContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações da BIOS",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	titleContainer := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
	)

	info := container.NewVBox(
		widget.NewLabelWithStyle("Fabricante", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(a.machineInfo.BIOS.Vendor),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Versão", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(a.machineInfo.BIOS.Version),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Data", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(a.machineInfo.BIOS.ReleaseDate),
	)

	return container.NewVBox(
		titleContainer,
		container.NewPadded(info),
	)
}

func (a *App) createUSBContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Dispositivos USB",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	titleContainer := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
	)

	info := container.NewVBox()
	if len(a.machineInfo.USBDevices) == 0 {
		info.Add(widget.NewLabel("Nenhum dispositivo USB encontrado"))
	} else {
		for _, usb := range a.machineInfo.USBDevices {
			usbBox := container.NewVBox(
				widget.NewLabelWithStyle("Nome", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(usb.Name),
				widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(fmt.Sprintf("%s:%s", usb.VendorID, usb.ProductID)),
			)
			if usb.Serial != "" {
				usbBox.Add(widget.NewLabelWithStyle("Serial", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
				usbBox.Add(widget.NewLabel(usb.Serial))
			}
			usbBox.Add(widget.NewSeparator())
			info.Add(usbBox)
		}
	}

	scroll := container.NewVScroll(container.NewPadded(info))
	scroll.SetMinSize(fyne.NewSize(600, 400))

	return container.NewVBox(
		titleContainer,
		scroll,
	)
}

func (a *App) updateLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case info := <-a.updateChan:
			a.machineInfo = info
			a.window.Canvas().Refresh(a.window.Content())
		case <-ticker.C:
			a.window.Canvas().Refresh(a.window.Content())
		}
	}
}
