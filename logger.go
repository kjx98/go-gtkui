/*
 * Copyright (c) 2019-2020 Jesse Kuang <jkuang@21cn.com>
 *
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package gtkui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/op/go-logging"
	"os"
	"sync"
)

const (
	maxLines = 128
)

var log = logging.MustGetLogger("gktui")

type Logger struct {
	tv *gtk.TextView
	sw *gtk.ScrolledWindow
}

func NewGtkLogger() (*Logger, error) {
	res := Logger{}
	if tv, err := gtk.TextViewNew(); err != nil {
		log.Error("Unable to create TextView:", err)
		return nil, err
	} else {
		tv.SetWrapMode(gtk.WRAP_WORD_CHAR)
		tv.SetEditable(false)
		res.tv = tv
	}
	if sw, err := gtk.ScrolledWindowNew(nil, nil); err != nil {
		log.Error("Unable to create scroll", err)
		res.tv.Destroy()
		return nil, err
	} else {
		sw.Add(res.tv)
		res.sw = sw
	}
	return &res, nil
}

var logLock sync.Mutex

func (w *Logger) prepend_text(text string) error {
	buffer, err := w.tv.GetBuffer()
	if err != nil {
		return err
	}
	logLock.Lock()
	si := buffer.GetIterAtLine(0)
	buffer.Insert(si, text)
	if cnt := buffer.GetLineCount(); cnt > maxLines {
		// delete lines after maxLines
		bi := buffer.GetIterAtLine(maxLines)
		ei := buffer.GetEndIter()
		buffer.Delete(bi, ei)
	}
	logLock.Unlock()
	return nil
}

func (w *Logger) append_text(text string) error {
	buffer, err := w.tv.GetBuffer()
	if err != nil {
		return err
	}
	logLock.Lock()
	if cnt := buffer.GetLineCount(); cnt > maxLines+16 {
		// delete lines after maxLines
		bi := buffer.GetIterAtLine(0)
		ei := buffer.GetIterAtLine(cnt - maxLines)
		buffer.Delete(bi, ei)
	}
	si := buffer.GetEndIter()
	buffer.Insert(si, text)
	si = buffer.GetEndIter()
	mk := buffer.CreateMark("lastLine", si, false)
	w.tv.ScrollMarkOnscreen(mk)
	buffer.DeleteMark(mk)
	for gtk.EventsPending() {
		gtk.MainIteration()
	}
	logLock.Unlock()
	return nil
}

func (w *Logger) Win() *gtk.ScrolledWindow {
	return w.sw
}

func (w *Logger) Write(p []byte) (n int, err error) {
	n = len(p)
	//err = w.prepend_text(string(p))
	err = w.append_text(string(p))
	//w.tv.QueueDraw()
	//w.sw.QueueDraw()
	return
}

//	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
func init() {
	var format = logging.MustStringFormatter(
		`%{color}%{time:01-02 15:04:05}  ▶ %{level:.4s} %{color:reset} %{message}`,
	)

	logback := logging.NewLogBackend(os.Stderr, "", 0)
	logfmt := logging.NewBackendFormatter(logback, format)
	logging.SetBackend(logfmt)
}
