package model

import (
	"fmt"
)

var (
	ErrUnsupportedMIMEType = "unsupported MIME %s"
	ErrImageTooLarge       = "image is too large!\nimage size: %vMB"
)

type Food struct {
	ID          uint
	Name        string
	Description string
	Type        string
	IsAvailable bool
	Price       float64
	BranchID    uint
	Images      []FoodImage
}

type FoodPrice struct {
	ID    uint
	Price float64
}

type FoodImage struct {
	MIMEType MIMEType
	ID       uint
	Image    []byte
	FoodID   uint
}

func (f *Food) Validate() error {
	return nil
}

type MIMEType string

var (
	image_jpeg MIMEType = "image/jpeg"
	image_png  MIMEType = "image/png"
	image_webp MIMEType = "image/webp"
	image_heic MIMEType = "image/heic"
)

var (
	supportedTypes = map[MIMEType]bool{
		image_jpeg: true, image_png: true,
		image_webp: true, image_heic: true,
	}
)

func (fi *FoodImage) Validate() error {
	if _, exists := supportedTypes[fi.MIMEType]; !exists {
		return fmt.Errorf(ErrUnsupportedMIMEType, fi.MIMEType)
	}

	if len(fi.Image) > 2.5*1024*1024 {
		return fmt.Errorf(ErrImageTooLarge, len(fi.Image)/1024/1024)
	}

	return nil
}
