package main

import (
	"fmt"
	"github.com/gographics/imagick/imagick"
)

func Join(tiles matrix) *imagick.MagickWand {
	imagick.Initialize()
	defer imagick.Terminate()

	for i := tiles.TL.x; i < tiles.TR.x; i++ {
		col(tiles.BL.y, tiles.TL.y, i)
	}

	return row(tiles.TL.x, tiles.TR.x)
}

func row(from, to coord) *imagick.MagickWand {
	base := newBase()
	defer base.Destroy()

	for i := from; i < to; i++ {
		addImage(base, fmt.Sprintf("%d.png", i))
	}

	wand := base.MontageImage(imagick.NewDrawingWand(), "x1", "", imagick.MONTAGE_MODE_CONCATENATE, "0x0")
	return wand
}

func col(from, to, x coord) {
	base := newBase()
	defer base.Destroy()

	for i := from; i < to; i++ {
		addImage(base, fmt.Sprintf("tile_%d:%d.png", x, i))
	}

	wand := base.MontageImage(imagick.NewDrawingWand(), "1x", "", imagick.MONTAGE_MODE_CONCATENATE, "0x0")
	wand.WriteImage(fmt.Sprintf("%d.png", x))
}

/**
 * Add image to base wand, I know there is way to
 * add them more at time, but that would constitute
 * making an array
 */
func addImage(base *imagick.MagickWand, filename string) {
	additive := imagick.NewMagickWand()
	err := additive.ReadImage(filename)
	if err != nil {
		panic(err)
	}
	defer additive.Destroy()
	base.AddImage(additive)
}

/**
 * Common place for config
 */
func newBase() *imagick.MagickWand {
	base := imagick.NewMagickWand()
	base.SetImageCompressionQuality(100)
	return base
}
