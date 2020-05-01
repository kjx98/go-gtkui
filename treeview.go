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
	"errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	errNoSymbol = errors.New("No quote for symbol")
)

type QuoteView struct {
	treeView  *gtk.TreeView
	listStore *gtk.ListStore
	nRows     int
	nCols     int
	cols      []int // Column ID slice
	rowsIter  []*gtk.TreeIter
}

// Add a column to the tree view (during the initialization of the tree view)
func createColumn(title string, id int) (*gtk.TreeViewColumn, error) {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Error("Unable to create text cell renderer:", err)
		return nil, err
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Error("Unable to create cell column:", err)
		return nil, err
	}

	return column, nil
}

// Creates a tree view and the list store that holds its data
//	colTitle	column title
//				column 0 MUST be symbol
func NewQuoteView(colTitle []string) (*QuoteView, error) {
	res := QuoteView{}
	types := []glib.Type{}
	if treeView, err := gtk.TreeViewNew(); err != nil {
		log.Error("Unable to create tree view:", err)
		return nil, err
	} else {
		res.treeView = treeView
		for idx, col := range colTitle {
			if cc, err := createColumn(col, idx); err != nil {
				treeView.Destroy()
				return nil, err
			} else {
				treeView.AppendColumn(cc)
				types = append(types, glib.TYPE_STRING)
			}
		}
		res.nCols = len(colTitle)
	}

	// Creating a list store. This is what holds the data that will be shown on our tree view.
	if listStore, err := gtk.ListStoreNew(types...); err != nil {
		res.treeView.Destroy()
		log.Error("Unable to create list store:", err)
		return nil, err
	} else {
		res.treeView.SetModel(listStore)
		res.listStore = listStore
	}

	return &res, nil
}

func (w *QuoteView) Win() gtk.IWidget {
	return w.treeView
}

// Append a row to the list store for the tree view
func (w *QuoteView) AddRow(sym string) {
	// Get an iterator for a new row at the end of the list store
	iter := w.listStore.Append()

	// Set the contents of the list store row that the iterator represents
	if err := w.listStore.Set(iter, []int{0}, []interface{}{sym}); err != nil {
		log.Error("Unable to add row", err)
	} else {
		w.nRows++
		w.rowsIter = append(w.rowsIter, iter)
	}
}

func (w *QuoteView) UpdateRow(row int, v []string) error {
	if row < 0 || row >= w.nRows {
		return errNoSymbol
	}
	nc := w.nCols
	if nc > len(v) {
		nc = len(v)
	}
	colIds := make([]int, nc)
	vv := make([]interface{}, nc)
	for idx := 0; idx < nc; idx++ {
		colIds[idx] = idx
		vv[idx] = v[idx]
	}
	if err := w.listStore.Set(w.rowsIter[row], colIds, vv); err != nil {
		log.Error("Update quote row", err)
	}
	return nil
}
