# Wireless Bitmap

[![GoDoc](https://godoc.org/github.com/bake/wbmp?status.svg)](https://godoc.org/github.com/bake/wbmp)
[![Go Report Card](https://goreportcard.com/badge/github.com/bake/wbmp)](https://goreportcard.com/report/github.com/bake/wbmp)

Package wbmp implements decoding for Wireless Bitmap images.

A Wireless Bitmap is a monochrome image that was used in WAP and MMS where each
bit represents one pixel. Bytes are read from left to right where the left most
pixel is the high bit of the first byte. Rows are always represent by whole
bytes. If the image width is not divisable by 8, the following row starts at the
next byte. Unused bits must be set to 0.
