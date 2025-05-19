package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
	"github.com/dev/falcon-agent/internal/model"
	"github.com/dev/falcon-agent/internal/service"
)

type App struct {
	window      fyne.Window
	machineInfo *model.MachineInfo
	updateChan  chan *model.MachineInfo
	content     *fyne.Container
	systemTray  fyne.App
}

func New(machineInfo *model.MachineInfo) *App {
	a := app.New()
	window := a.NewWindow("Falcon Agent")
	window.Resize(fyne.NewSize(600, 400))
	window.SetMaster()

	// Configurar o comportamento de minimização
	window.SetCloseIntercept(func() {
		window.Hide()
	})

	// Configurar para não mostrar na barra de tarefas
	window.SetIcon(theme.ComputerIcon())
	window.CenterOnScreen()

	// Configurar o comportamento da janela
	window.SetOnClosed(func() {
		window.Hide()
	})

	updateChan := make(chan *model.MachineInfo, 1)
	updateChan <- machineInfo

	app := &App{
		window:      window,
		machineInfo: machineInfo,
		updateChan:  updateChan,
		systemTray:  a,
	}

	app.setupUI()
	go app.setupSystemTray()
	go app.updateLoop()

	return app
}

func (a *App) setupSystemTray() {
	systray.Run(func() {
		systray.SetIcon(theme.ComputerIcon().Content())
		systray.SetTitle("Falcon Agent")
		systray.SetTooltip("Falcon Agent")

		mShow := systray.AddMenuItem("Mostrar", "Mostrar janela")
		mQuit := systray.AddMenuItem("Sair", "Sair do aplicativo")

		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					a.window.Show()
					a.window.RequestFocus()
				case <-mQuit.ClickedCh:
					systray.Quit()
					a.systemTray.Quit()
					return
				}
			}
		}()
	}, nil)
}

func (a *App) Run() {
	a.window.ShowAndRun()
}

func (a *App) setupUI() {
	sidebar := a.createSidebar()
	separator := canvas.NewRectangle(theme.ShadowColor())
	separator.SetMinSize(fyne.NewSize(2, 0))
	a.content = a.createSystemContent()
	split := container.NewHSplit(
		container.NewHBox(sidebar, separator),
		container.NewPadded(a.content),
	)
	split.SetOffset(0.2)
	a.window.SetContent(split)
}

// Novo: Card moderno com fundo diferenciado
func createModernCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.NRGBA{R: 36, G: 37, B: 46, A: 255})
	bg.SetMinSize(fyne.NewSize(340, 120))

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	cardContent := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		content,
	)

	return container.NewMax(
		bg,
		container.NewPadded(cardContent),
	)
}

// Novo: Botão moderno destacado
func createModernButton(text string, icon fyne.Resource, tapped func()) fyne.CanvasObject {
	btn := widget.NewButtonWithIcon(text, icon, tapped)
	btn.Importance = widget.HighImportance
	return container.NewPadded(btn)
}

func (a *App) createSidebar() *fyne.Container {
	menuContainer := container.NewVBox()
	menuContainer.Resize(fyne.NewSize(200, 0))

	buttons := []struct {
		icon   fyne.Resource
		text   string
		action func()
	}{
		{theme.ComputerIcon(), "Sistema", func() {
			a.content.Objects = []fyne.CanvasObject{a.createSystemContent()}
			a.content.Refresh()
		}},
		{theme.MediaPlayIcon(), "CPU", func() {
			a.content.Objects = []fyne.CanvasObject{a.createProcessorContent()}
			a.content.Refresh()
		}},
		{theme.StorageIcon(), "Memória", func() {
			a.content.Objects = []fyne.CanvasObject{a.createMemoryContent()}
			a.content.Refresh()
		}},
		{theme.FolderIcon(), "Armazenamento", func() {
			a.content.Objects = []fyne.CanvasObject{a.createStorageContent()}
			a.content.Refresh()
		}},
		{theme.SettingsIcon(), "BIOS", func() {
			a.content.Objects = []fyne.CanvasObject{a.createBIOSContent()}
			a.content.Refresh()
		}},
		{theme.MediaRecordIcon(), "USB", func() {
			a.content.Objects = []fyne.CanvasObject{a.createUSBContent()}
			a.content.Refresh()
		}},
	}

	for _, b := range buttons {
		menuContainer.Add(createModernButton(b.text, b.icon, b.action))
	}

	menuContainer.Add(layout.NewSpacer())
	return menuContainer
}

func (a *App) createSystemContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações do Sistema",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	grid := container.NewGridWithColumns(2,
		createModernCard("Sistema Operacional", widget.NewLabel(a.machineInfo.OS)),
		createModernCard("Hostname", widget.NewLabel(a.machineInfo.Hostname)),
		createModernCard("Serial Number", widget.NewLabel(a.machineInfo.SerialNumber)),
	)
	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(grid),
	)
}

func (a *App) createProcessorContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações do Processador",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	grid := container.NewGridWithColumns(2,
		createModernCard("Modelo", widget.NewLabel(a.machineInfo.Processor.Model)),
		createModernCard("Núcleos", widget.NewLabel(fmt.Sprintf("%d", a.machineInfo.Processor.Cores))),
		createModernCard("Threads", widget.NewLabel(fmt.Sprintf("%d", a.machineInfo.Processor.Threads))),
		createModernCard("Frequência", widget.NewLabel(fmt.Sprintf("%.2f GHz", a.machineInfo.Processor.FrequencyGHz))),
	)
	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(grid),
	)
}

func (a *App) createMemoryContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações da Memória",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	content := container.NewVBox()
	for _, mem := range a.machineInfo.Memory {
		card := createModernCard("Slot "+mem.Slot, container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Capacidade: %d MB", mem.SizeMB)),
			widget.NewLabel(fmt.Sprintf("Fabricante: %s", mem.Manufacturer)),
			widget.NewLabel(fmt.Sprintf("Serial: %s", mem.SerialNumber)),
		))
		content.Add(card)
	}
	scroll := container.NewVScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 400))
	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(scroll),
	)
}

func (a *App) createStorageContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações de Armazenamento",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	content := container.NewVBox()
	for _, hd := range a.machineInfo.HDs {
		card := createModernCard(hd.Model, container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Capacidade: %d GB", hd.SizeGB)),
			widget.NewLabel(fmt.Sprintf("Serial: %s", hd.Serial)),
		))
		content.Add(card)
	}
	scroll := container.NewVScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 400))
	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(scroll),
	)
}

func (a *App) createBIOSContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Informações da BIOS",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	grid := container.NewGridWithColumns(2,
		createModernCard("Fabricante", widget.NewLabel(a.machineInfo.BIOS.Vendor)),
		createModernCard("Versão", widget.NewLabel(a.machineInfo.BIOS.Version)),
		createModernCard("Data de Lançamento", widget.NewLabel(a.machineInfo.BIOS.ReleaseDate)),
	)
	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(grid),
	)
}

func (a *App) createUSBContent() *fyne.Container {
	title := widget.NewLabelWithStyle(
		"Dispositivos USB",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	content := container.NewVBox()
	for _, usb := range a.machineInfo.USBDevices {
		card := createModernCard(usb.Name, container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Vendor ID: %s", usb.VendorID)),
			widget.NewLabel(fmt.Sprintf("Product ID: %s", usb.ProductID)),
			widget.NewLabel(fmt.Sprintf("Serial: %s", usb.Serial)),
		))
		content.Add(card)
	}
	scroll := container.NewVScroll(content)
	scroll.SetMinSize(fyne.NewSize(600, 400))
	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(scroll),
	)
}

func (a *App) updateLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if newInfo, err := service.CollectMachineInfo(); err == nil {
				a.updateChan <- newInfo
				a.machineInfo = newInfo
				a.content.Refresh()
			}
		}
	}
}
