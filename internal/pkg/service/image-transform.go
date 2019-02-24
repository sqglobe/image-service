package service

import (
	"image"
	"image/draw"

	"github.com/nfnt/resize"
)

// ImageTransformer - is the interface of type, that can draws image  over the background and adds watermark.
// This type expands the 'service.Transformer' interface.
type ImageTransformer struct{}

func getResizeDimention(bgRect image.Rectangle, resizeRect image.Rectangle) (uint, uint) {
	heightRatio := bgRect.Dy() / resizeRect.Dy()
	widthRatio := bgRect.Dx() / resizeRect.Dx()
	if heightRatio > widthRatio {
		return uint(bgRect.Dx()), 0
	}
	return 0, uint(bgRect.Dy())
}

func getBoundsInCenter(bgRect image.Rectangle, imgRect image.Rectangle) image.Rectangle {
	offsetX := int((bgRect.Dx() - imgRect.Dx()) / 2)
	offsetY := int((bgRect.Dy() - imgRect.Dy()) / 2)
	return imgRect.Add(image.Pt(offsetX, offsetY))
}

// Overlay draws image 'topImg' over 'backgroundImg' in center.
// 'topImg' will be scaled to size of 'backgroundImg' proportional.
func (ImageTransformer) Overlay(backgroundImg draw.Image, topImg image.Image) error {
	resizeWidth, resizeHeight := getResizeDimention(backgroundImg.Bounds(), topImg.Bounds())
	resizedImg := resize.Resize(resizeWidth, resizeHeight, topImg, resize.Bicubic)
	draw.Draw(backgroundImg, getBoundsInCenter(backgroundImg.Bounds(), resizedImg.Bounds()), resizedImg, image.ZP, draw.Over)
	return nil
}

// PutWatermark adds watermark to the 'backgroundImg' image.
// 'watermark' will draws so many times as 'backgroundImg' size allows
func (st ImageTransformer) PutWatermark(backgroundImg draw.Image, watermark image.Image) error {
	heightRatio := int(backgroundImg.Bounds().Dy() / watermark.Bounds().Dy())
	widthRatio := int(backgroundImg.Bounds().Dx() / watermark.Bounds().Dx())

	if heightRatio == 0 || widthRatio == 0 {
		return st.Overlay(backgroundImg, watermark)
	}

	offsetX := int((backgroundImg.Bounds().Dx() - watermark.Bounds().Dx()*widthRatio) / 2)
	offsetY := int((backgroundImg.Bounds().Dy() - watermark.Bounds().Dy()*heightRatio) / 2)
	for row := 0; row < heightRatio; row++ {
		for col := 0; col < widthRatio; col++ {
			offsetPt := image.Pt(offsetX+watermark.Bounds().Dx()*col, offsetY+watermark.Bounds().Dy()*row)
			draw.Draw(backgroundImg, watermark.Bounds().Add(offsetPt), watermark, image.ZP, draw.Over)
		}
	}
	return nil
}
