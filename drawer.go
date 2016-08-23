/*
 * Copyright (c) 2016 Artur Grabowski <art@blahonga.org>
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

// Package split_drawer implements an image/draw.Drawer that splits
// the Draw operation into multiple goroutines. It is intended only to
// be used for extremely large images or if the source/destination 
// images for some reason are computationally heavy (like an on-the-fly
// generated source image).
//
// The spliting is done on rows in the destination rectangle. If the
// number of rows is small, this splitting might not help or be quite
// uneven.
package split_drawer

import (
	"image"
	"image/draw"
	"sync"
	"runtime"
)

type splitDrawer struct {
	d draw.Drawer
	n int
}

// Default Drawer with sensible defaults.
var D draw.Drawer = splitDrawer{}

// Returns a draw.Drawer which will perform the d.Draw operation using
// n goroutines. If n is 0 and/or d is nil sensible defaults will be used.
func New(d draw.Drawer, n int) draw.Drawer {
	return splitDrawer{d,n}
}

func (sd splitDrawer)Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	n := sd.n
	if n == 0 {
		n = runtime.NumCPU()
	}

	d := sd.d
	if d == nil {
		d = draw.Src
	}

	rows := r.Max.Y - r.Min.Y

	if rows < n {
		n = rows
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(y0, y1 int) {
			d.Draw(dst, image.Rect(r.Min.X, r.Min.Y+y0, r.Max.X, r.Min.Y+y1), src, sp.Add(image.Pt(0,y0)))
			wg.Done()
		}(i*rows/n, (i+1)*rows/n)
	}
	wg.Wait()
}
