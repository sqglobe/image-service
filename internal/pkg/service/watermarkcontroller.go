package service

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// WatermarkController handles request to resize image and add watermark.
// Supposed to request has type 'multipart/form-data' and consist of two images:
// 1. 'image' will be resized and used as background for watermark. It should has jpeg format
// 2. 'watermark' will be used as watermark
// Result image will be saved in 'imagePath'
type WatermarkController struct {
	transformer  Transformer
	resizeHeight int
	resizeWidth  int
	imagePath    string
}

// NewWatermarkController creates new controller. It takes possible to specify resize width and height, and set up store path for result images.
func NewWatermarkController(transformer Transformer, resizeWidth, resizeHeight int, imagePath string) *WatermarkController {
	return &WatermarkController{
		transformer:  transformer,
		resizeWidth:  resizeWidth,
		resizeHeight: resizeHeight,
		imagePath:    imagePath,
	}
}

// Function handle request to resize image and put watermark.
func (contr *WatermarkController) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(response, fmt.Sprintf("undefiled error %s", err), http.StatusInternalServerError)
		}
	}()

	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	successChan := make(chan string)
	errorChan := make(chan string)

	go contr.exec(request, successChan, errorChan)

	select {
	case err := <-errorChan:
		http.Error(response, err, http.StatusInternalServerError)
	case succ := <-successChan:
		io.WriteString(response, succ)
	}
}

func (contr *WatermarkController) exec(request *http.Request, successChan, errChan chan string) {
	if err := request.ParseMultipartForm(32 << 20); err != nil {
		errChan <- fmt.Sprintf("Needed multipart data %s", err)
		return
	}

	imageFile, _, err := request.FormFile("image")

	if err != nil {
		errChan <- fmt.Sprintf("Failed get form part 'image': %s", err)
		return
	}

	defer imageFile.Close()

	imageObj, err := jpeg.Decode(imageFile)

	if err != nil {
		errChan <- fmt.Sprintf("Failed deecode jpeg: %s", err)
		return
	}

	watermarkFile, _, err := request.FormFile("watermark")
	if err != nil {
		errChan <- fmt.Sprintf("Failed get form part 'watermark': %s", err)
		return
	}

	defer watermarkFile.Close()

	watermarkObj, err := png.Decode(watermarkFile)

	if err != nil {
		errChan <- fmt.Sprintf("Failed deecode png: %s", err)
		return
	}

	bg := image.NewRGBA(image.Rect(0, 0, contr.resizeWidth, contr.resizeHeight))
	if err = contr.transformer.Overlay(bg, imageObj); err != nil {
		errChan <- fmt.Sprintf("Failed add overlay: %s", err)
		return
	}

	if err = contr.transformer.PutWatermark(bg, watermarkObj); err != nil {
		errChan <- fmt.Sprintf("Failed add watermark: %s", err)
		return
	}

	fileName := contr.imagePath + "/" + uuid.New().String() + ".png"
	imgOut, err := os.Create(fileName)

	if err != nil {
		errChan <- fmt.Sprintf("Failed open out file: %s, error: %s", fileName, err)
		return
	}

	defer imgOut.Close()

	if err = png.Encode(imgOut, bg); err != nil {
		errChan <- fmt.Sprintf("Failed open out file: %s, error: %s", fileName, err)
		return
	}
	successChan <- fmt.Sprintf("Success. File saved as %s", fileName)
}
