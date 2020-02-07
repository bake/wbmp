// Package wbmp implements decoding for Wireless Bitmap images.
//
// A Wireless Bitmap is a monochrome image that was used in WAP and MMS where
// each bit represents one pixel. Bytes are read from left to right where the
// left most pixel is the high bit of the first byte. Rows are always represent
// by whole bytes. If the image width is not divisable by 8, the following row
// starts at the next byte. Unused bits must be set to 0.
//
// http://www.wapforum.org/what/technical/SPEC-WAESpec-19990524.pdf
package wbmp

import (
	"bufio"
	"encoding/binary"
	"image"
	"image/color"
	"io"
	"io/ioutil"
)

func init() {
	image.RegisterFormat("wbmp", "", Decode, DecodeConfig)
}

// WBMP represents a Wireless Bitmap and implements an image.Image interface.
type WBMP struct {
	typeField      uint64
	fixHeaderField uint8
	width, height  uint64
	data           []byte
}

// Bounds returns the bounds of the image.
func (wbmp WBMP) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(wbmp.width), int(wbmp.height))
}

// ColorModel returns the Image's color model.
func (wbmp WBMP) ColorModel() color.Model {
	return color.GrayModel
}

// At returns the color of the pixel at (x, y).
func (wbmp WBMP) At(x, y int) color.Color {
	// w contains the number of bytes per row.
	w := int(wbmp.width)
	if wbmp.width%8 > 0 {
		w += 8
	}
	w /= 8

	b := wbmp.data[y*w+x/8]
	if b>>(7-x%8)&1 == 1 {
		return color.White
	}
	return color.Black
}

// Decode decodes the image.
func Decode(r io.Reader) (image.Image, error) {
	br := bufio.NewReader(r)
	wbmp, err := decodeConfig(br)
	if err != nil {
		return nil, err
	}
	wbmp.data, err = ioutil.ReadAll(br)
	return wbmp, err
}

// DecodeConfig decodes the color model and dimensions of an image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	br := bufio.NewReader(r)
	wbmp, err := decodeConfig(br)
	if err != nil {
		return image.Config{}, err
	}
	return image.Config{
		Width:  int(wbmp.width),
		Height: int(wbmp.height),
	}, nil
}

func decodeConfig(r *bufio.Reader) (m WBMP, err error) {
	m.typeField, err = binary.ReadUvarint(r)
	if err != nil {
		return m, err
	}
	m.fixHeaderField, err = r.ReadByte()
	if err != nil {
		return m, err
	}
	m.width, err = binary.ReadUvarint(r)
	if err != nil {
		return m, err
	}
	m.height, err = binary.ReadUvarint(r)
	if err != nil {
		return m, err
	}
	return m, nil
}
