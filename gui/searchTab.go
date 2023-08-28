package gui

import (
	"fmt"

	Inv "gm/invidious"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

var entry = xwidget.NewCompletionEntry([]string{})
var searchResault = container.New(layout.NewGridLayout(1))
var sc = container.NewVScroll(searchResault)

type tLabel struct {
	widget.Label
	Id       string
	OnTapped func(string) `json:"-"`
}

func (tL *tLabel) Tapped(*fyne.PointEvent) {
	if tL.OnTapped != nil {
		tL.OnTapped(tL.Id)
	}
}

func (tL *tLabel) TappedSecondary(*fyne.PointEvent) {}

func newtLabel(text string) *tLabel {
	return newtLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{})
}

func newtLabelWithStyle(text string, alignment fyne.TextAlign, style fyne.TextStyle) *tLabel {
	var tl = &tLabel{}
	tl.ExtendBaseWidget(tl)
	tl.Text = text
	tl.Alignment = alignment
	tl.TextStyle = style

	return tl
}

func getSearchTab() fyne.CanvasObject {
	entry.Wrapping = fyne.TextTruncate
	entry.OnSubmitted = func(s string) { onSubmitted(s) }
	entry.OnChanged = func(s string) { onChanged(s) }

	sc.SetMinSize(fyne.NewSize(500, 400))

	return container.NewVBox(entry, sc)
}

func onSubmitted(s string) {
	go updateSearchResault(s)
	entry.HideCompletion()
	entry.SetText("")
}

func onChanged(s string) {
	if len(s) < 3 {
		entry.HideCompletion()
		return
	}

	go func() {
		var suggestion, err = Inv.GetSuggestions(config.Instance, entry.Text)
		if err != nil {
			entry.HideCompletion()
			return
		}

		if len(suggestion) == 0 {
			entry.HideCompletion()
			return
		}

		entry.SetOptions(suggestion)
		entry.ShowCompletion()
	}()
}

func updateSearchResault(query string) {
	searchResault.RemoveAll()
	var videos, _ = Inv.Search(config.Instance, query, "video", 1)
	var uri fyne.URI
	var thumbnail *canvas.Image
	var title *tLabel
	var channel *tLabel
	var duration *widget.Label

	for _, v := range videos {
		uri, _ = storage.ParseURI(v.VideoThumbnail[len(v.VideoThumbnail)-1].Url)
		thumbnail = canvas.NewImageFromURI(uri)

		title = newtLabel(v.Title)
		title.Id = v.VideoId
		title.Wrapping = fyne.TextWrapWord
		title.TextStyle = fyne.TextStyle{Bold: true}
		title.OnTapped = func(id string) {
			var video = Inv.GetVideo(config.Instance, id)
			player.Play(video.Format[0].Url)
			player.Resume()
			go updateHomeTab(video)
		}

		channel = newtLabel(v.Author)
		channel.Wrapping = fyne.TextWrapWord

		duration = widget.NewLabel(fmt.Sprintf("%02d:%02d",
			v.Length/60, v.Length%60))
		duration.Wrapping = fyne.TextWrapWord
		searchResault.Add(container.NewGridWithColumns(2, thumbnail,
			container.NewGridWithRows(3, title, channel, duration)))
	}
}
