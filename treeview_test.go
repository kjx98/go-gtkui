package gtkui

import (
	"github.com/gotk3/gotk3/gtk"
	"testing"
	"time"
)

func TestQuoteView(t *testing.T) {
	win := setup_window("test quote")
	var qv *QuoteView
	cols := []string{" Symbol ", "  Last  ", "  Volume  ", "  Open  ",
		"  High  ", "   Low  "}
	if ll, err := NewQuoteView(cols); err != nil {
		t.Error(err)
	} else {
		qv = ll
	}
	win.Add(qv.Win())
	qv.AddRow("$SPX.X")
	qv.AddRow("/ES")
	qv.AddRow("SPY")
	win.ShowAll()
	go func() {
		time.Sleep(30 * time.Second)
		gtk.MainQuit()
	}()
	gtk.Main()
}
