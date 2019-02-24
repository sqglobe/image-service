package service

import (
	"image"
	"image/draw"
)

// Transformer - is the interface of type, that can draws image  over the background and adds watermark.
type Transformer interface {
	// Overlay draws image 'topImg' over the top of 'backgroundImg'
	Overlay(backgroundImg draw.Image, topImg image.Image) error

	// PutWatermark adds 'watermark' to the 'backgroundImg'
	PutWatermark(backgroundImg draw.Image, watermark image.Image) error
}
