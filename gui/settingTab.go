package gui

import (
	Pip "gm/piped"
	Pl "gm/player"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func getSettingTab() fyne.CanvasObject {
	var instancesLabel = widget.NewLabel("Piped Instances")
	var selectInstance = widget.NewSelect([]string{},
		func(s string) {
			config.Set("piped.instance", s)
			config.Instance = s
			log.Println("Using", s)
		})
	selectInstance.Selected = config.Instance

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		log.Println("Refreshing instances")
		go func() {
			var ins []string = Pip.GetInstances()
			selectInstance.Options = ins
			selectInstance.SetSelected(ins[0])
		}()
	})

	var urlL = widget.NewLabel("HTTP Proxy")
	var urlE = widget.NewEntry()
	urlE.Text = config.Proxy.Url
	urlE.PlaceHolder = "http://127.0.0.1"
	urlE.Wrapping = fyne.TextTruncate

	var portL = widget.NewLabel("Port")
	var portE = widget.NewEntry()
	portE.Text = config.Proxy.Port
	portE.PlaceHolder = "80"
	portE.Wrapping = fyne.TextTruncate

	var applyBtn = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(),
		func() {
			// TODO proxy validate
			config.Set("proxy.url", urlE.Text)
			config.Set("proxy.port", portE.Text)
			player.Close()
			player = Pl.Init(urlE.Text + ":" + portE.Text)
		})

	return container.NewVBox(layout.NewSpacer(),
		container.NewGridWithColumns(1, container.NewGridWithColumns(2, instancesLabel, selectInstance), refreshBtn),
		layout.NewSpacer(),
		container.NewGridWithColumns(4, urlL, urlE, portL, portE), applyBtn,
		layout.NewSpacer())
}
