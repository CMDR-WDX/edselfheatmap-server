package img

import "image/color"

func GetColorForCount(count int) color.RGBA {
	if count == 0 {
		return color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		}
	}
	if count == 1 {
		return color.RGBA{
			R: 0,
			G: 0,
			B: 129,
			A: 255,
		}
	}
	if count >= 2 && count < 5 {
		return color.RGBA{
			R: 0,
			G: 0,
			B: 255,
			A: 255,
		}
	}
	if count >= 5 {
		return color.RGBA{
			R: 61,
			G: 61,
			B: 255,
			A: 255,
		}
	}
	return color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
}
