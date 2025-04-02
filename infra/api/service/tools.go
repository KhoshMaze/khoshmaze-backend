package service

import (
	"fmt"
	"image/color"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	qr "github.com/skip2/go-qrcode"
)

func CreateQR(req *pb.QrCodeRequest) ([]byte, error) {
	qrcode, err := qr.New(fmt.Sprintf("https://namak-app.ir/%s", req.GetUrl()), qr.Highest)
	if err != nil {
		return nil, err
	}

	bg := req.GetBackgroundColor()
	fg := req.GetForegroundColor()
	
	if fg == nil {
		fg = &pb.QrCodeRequest_Color{R: 0, G: 0, B: 0}
	}

	if bg == nil {
		bg = &pb.QrCodeRequest_Color{R: 255, G: 255, B: 255}
	}

	qrcode.BackgroundColor = color.RGBA{R: uint8(bg.GetR()), G: uint8(bg.GetG()), B: uint8(bg.GetB()), A: 255}
	qrcode.ForegroundColor = color.RGBA{R: uint8(fg.GetR()), G: uint8(fg.GetG()), B: uint8(fg.GetB()), A: 255}
	qrcode.DisableBorder = !req.GetHasBorder()

	size := req.GetSize()
	if size <= 0 {
		size = 256
	}

	data, err := qrcode.PNG(int(size))

	return data, err
}
