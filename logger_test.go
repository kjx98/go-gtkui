package gtkui

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"testing"
	"time"
)

func init() {
	gtk.Init(nil)
}

func setup_window(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func TestLogger(t *testing.T) {
	win := setup_window("test logger")
	var logg *Logger
	if ll, err := NewGtkLogger(); err != nil {
		t.Error(err)
	} else {
		logg = ll
	}
	win.Add(logg.Win())
	ss := "line1\nline2\ntest ok\n"
	logg.Write([]byte(ss))
	win.ShowAll()
	for i := 0; i < 200; i++ {
		ss := fmt.Sprintln("log line ", i)
		logg.Write([]byte(ss))
	}
	go func() {
		time.Sleep(30 * time.Second)
		gtk.MainQuit()
	}()
	gtk.Main()
}
