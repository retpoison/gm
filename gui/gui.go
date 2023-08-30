package gui

import (
	Conf "gm/config"
	Pl "gm/player"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var player Pl.Player
var config Conf.Config

func Run(p Pl.Player, c Conf.Config) {
	config = c
	player = p

	log.Println("Initializing gui.")
	app := app.New()
	window := app.NewWindow("gm")
	tabs := container.NewAppTabs(
		container.NewTabItem("Home", getHomeTab()),
		container.NewTabItem("Search", getSearchTab()),
		container.NewTabItem("Setting", getSettingTab()),
	)

	window.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if string(e.Name) == "Space" && (*tabs).CurrentTab().Text == "Home" {
			player.PlayPause()
		}
	})

	tabs.SetTabLocation(container.TabLocationLeading)
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(500, 400))
	window.SetContent(tabs)

	log.Println("Runing gui.")
	window.ShowAndRun()
}
