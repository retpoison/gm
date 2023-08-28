package main

import (
	Conf "gm/config"
	Gui "gm/gui"
	Pl "gm/player"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var config Conf.Config = Conf.Init()
	var player Pl.Player = Pl.Init(config.Proxy.Url + ":" + config.Proxy.Port)
	Gui.Run(player, config)
}
