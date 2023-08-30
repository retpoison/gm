package gui

import (
	"fmt"
	"time"

	Pip "gm/piped"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var video Pip.Video
var thumbnailTitleContainer *fyne.Container

func updateHomeTab(v Pip.Video) {
	video = v
	thumbnailTitleContainer.Objects[0] = canvas.NewImageFromImage(video.GetThumbnail())
	thumbnailTitleContainer.Objects[1].(*widget.Label).SetText(fmt.Sprintf("%s\n%s",
		video.Title, video.Uploader))
}

func getHomeTab() fyne.CanvasObject {
	var playpause = binding.NewString()
	playpause.Set("Pause")

	previous := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(),
		func() {
			// TODO
		})
	pause := widget.NewButtonWithIcon("", theme.MediaPauseIcon(),
		func() {
			player.Pause()
			playpause.Set("Pause")
		})
	play := widget.NewButtonWithIcon("", theme.MediaPlayIcon(),
		func() {
			player.Resume()
			playpause.Set("Play")
		})
	next := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(),
		func() {
			// TODO
		})

	var data = binding.NewFloat()
	data.Set(0)
	var slide = widget.NewSliderWithData(0, 100, data)
	slide.Step = 0.1

	var currentPosition = binding.NewString()
	currentPosition.Set("")
	var cpLabel = widget.NewLabelWithData(currentPosition)
	cpLabel.Alignment = fyne.TextAlignCenter

	var playpauseLabel = widget.NewLabelWithData(playpause)
	playpauseLabel.Alignment = fyne.TextAlignCenter

	var title *widget.Label
	var thumbnail *canvas.Image

	thumbnail = canvas.NewImageFromResource(theme.FileImageIcon())

	title = widget.NewLabel(video.Title)
	title.Wrapping = fyne.TextWrapWord
	title.Alignment = fyne.TextAlignCenter

	thumbnailTitleContainer = container.NewGridWithRows(2, thumbnail, title)

	go func() {
		for range time.Tick(time.Second) {
			updateStatus(data, currentPosition, playpause)
		}
	}()

	return container.NewBorder(nil,
		container.NewVBox(widget.NewSeparator(),
			container.NewGridWithColumns(4, previous, pause, play, next),
			container.NewGridWithColumns(3, cpLabel, layout.NewSpacer(), playpauseLabel),
			slide),
		nil, nil,
		thumbnailTitleContainer)
}

func updateStatus(slide binding.Float,
	cPos binding.String,
	playPause binding.String) {

	var pStatus bool = player.PauseStatus()
	if pStatus == false {
		slide.Set(player.PercentPosition())

		var cp float64 = player.Position()
		var cpString string = fmt.Sprintf("%s / %02d:%02d:%02d",
			video.FormatedDuration, int(cp)/60/60, int(cp)/60%60, int(cp)%60)
		cPos.Set(cpString)
	}

	switch pStatus {
	case false:
		playPause.Set("Play")
	case true:
		playPause.Set("Pause")
	}
}
