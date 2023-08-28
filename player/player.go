package player

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/blang/mpv"
)

type Player struct {
	cmd    *exec.Cmd
	client *mpv.Client
}

func Init(proxy string) Player {
	var player Player = Player{}
	runMpv(&player, proxy)

	for {
		player.client = getClient()
		if player.client == nil {
			time.Sleep(time.Second * 2)
			log.Println("Trying again.")
		} else {
			break
		}
	}

	log.Println("connected to mpv.")
	player.Pause()
	return player
}

func (pl *Player) Close() {
	killMpv(pl)
}

func runMpv(pl *Player, proxy string) {
	if proxy != ":" {
		log.Println("Set http proxy.")
		os.Setenv("http_proxy", proxy)
	}

	var mpvSocket string = "/tmp/mpvsocket"
	log.Println("Run mpv ipc server on", mpvSocket+".")
	pl.cmd = exec.Command("mpv", "--idle", "--no-terminal", fmt.Sprintf("%s%s", "--input-ipc-server=", mpvSocket))
	var err error = pl.cmd.Start()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Running mpv done.")
}

func killMpv(pl *Player) {
	pl.cmd.Process.Kill()
}

func getClient() *mpv.Client {
	defer handlePanic()
	var ipcc = mpv.NewIPCClient("/tmp/mpvsocket")
	var client = mpv.NewClient(ipcc)
	return client
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Println(r, "trying in 2 seconds.")
	}
}

func (pl *Player) Play(url string) {
	pl.client.Loadfile(url, mpv.LoadFileModeReplace)
}

func (pl *Player) Pause() {
	pl.client.SetPause(true)
}

func (pl *Player) Resume() {
	pl.client.SetPause(false)
}

func (pl *Player) PlayPause() {
	pl.client.SetPause(!(pl.PauseStatus()))
}

func (pl *Player) Seek(i int) {
	pl.client.Seek(i, mpv.SeekModeRelative)
}

func (pl *Player) Position() float64 {
	var pos, _ = pl.client.Position()
	return pos
}

func (pl *Player) PercentPosition() float64 {
	var posPercent, _ = pl.client.PercentPosition()
	return posPercent
}

func (pl *Player) PauseStatus() bool {
	var pause, _ = pl.client.Pause()
	return pause
}
