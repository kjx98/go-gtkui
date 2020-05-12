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
	if it, ok := qv.FirstRow(); !ok {
		t.Error("can't got first row")
	} else if idx := qv.RowId(it); idx != 0 {
		t.Error("expect row 0, got", idx)
	} else if !qv.NextRow(it) {
		t.Error("can't go next row")
	} else if idx = qv.RowId(it); idx != 1 {
		t.Error("expect row 1, got", idx)
	} else {
		c := qv.RowColor(it)
		t.Log("Color of Row", idx, "is", c)
		qv.SetRowColor(it, ColorUp)
		if idx = qv.RowId(it); idx != 1 {
			t.Error("expect row 1, got", idx)
		}
		if c = qv.RowColor(it); c != ColorUp {
			t.Error("Expect", ColorUp, ", BUT got", c)
		}
		t.Log("After SetRowColor Color of Row", idx, "is", c)
	}
	t.Log("FirstRow/NextRow works")
	gtk.Main()
}
