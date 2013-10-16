package main

import (
	"fmt"
	"github.com/gographics/imagick/imagick"
)

func Join(tiles matrix) {
	// for dy := tiles.BL.y; dy < tiles.TL.y; dy++ {

	// }
	imagick.Initialize()
	defer imagick.Terminate()

	for i := tiles.TL.x; i < tiles.TR.x; i++ {
		col(tiles.BL.y, tiles.TL.y, i)
	}
	row(tiles.TL.x, tiles.TR.x)
}

func row(from, to coord) {
	base := imagick.NewMagickWand()
	base.SetImageCompressionQuality(100)
	defer base.Destroy()

	for i := from; i < to; i++ {
		additive := imagick.NewMagickWand()
		filename := fmt.Sprintf("%d.png", i)
		err := additive.ReadImage(filename)
		handle(err)

		base.AddImage(additive)
		additive.Destroy()
	}

	wand := base.MontageImage(imagick.NewDrawingWand(), "5x5", "", imagick.MONTAGE_MODE_CONCATENATE, "0x0")
	wand.WriteImage("complete.png")

}

func col(from, to, x coord) {
	base := imagick.NewMagickWand()
	base.SetImageCompressionQuality(100)
	defer base.Destroy()

	for i := from; i < to; i++ {
		additive := imagick.NewMagickWand()
		filename := fmt.Sprintf("%d_%d.png", x, i)
		err := additive.ReadImage(filename)
		handle(err)

		base.AddImage(additive)
		additive.Destroy()
	}

	wand := base.MontageImage(imagick.NewDrawingWand(), "1x", "", imagick.MONTAGE_MODE_CONCATENATE, "0x0")
	wand.WriteImage(fmt.Sprintf("%d.png", x))
}
