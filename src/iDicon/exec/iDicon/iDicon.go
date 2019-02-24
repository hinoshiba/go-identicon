package main

import (
	"os"
	"fmt"
	"flag"
	"image"
	"image/color"
	"image/png"
	"errors"
	"crypto/md5"
)

import (
	"github.com/lucasb-eyer/go-colorful"
)

func die(s string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, s + "\n" , msg...)
	os.Exit(1)
}

var Seed  string
var Size int64
var Fpath string
var BG_COLOR color.RGBA = color.RGBA{212, 212, 212, 255}

func iDicon() error {
	fbit, cbit, err := seed2bits(Seed)
	if err != nil {
		return err
	}
	forest, err := b60forest(fbit)
	if err != nil {
		return err
	}
	rgb, err := b28rgb(cbit)
	if err != nil {
		return err
	}

	img, err := array2img(forest, rgb, Size)
	if err != nil {
		return err
	}
	if err := writter(Fpath, img); err != nil {
		return err
	}
	return nil
}

func seed2bits(seed string) ([]byte, []byte, error) {
	h := md5.Sum([]byte(seed))

	fbit := rShift4bit(h[0:8])

	cbit := h[12:16]
	cbit[0] = cbit[0] & 15

	return fbit, cbit, nil
}

func rShift4bit(b []byte) []byte {
	var buf byte
	for i, _ := range b {
		buff := b[i] << 4
		b[i] = b[i] >> 4
		b[i] = b[i] | buf
		buf = buff
	}
	return b
}

func b60forest(b []byte) ([][]bool, error) {
	if len(b) < 8 {
		return [][]bool{}, errors.New("less than 8 byte.")
	}

	var f []bool
	for _, v := range b {
		h := v >> 4
		f = append(f, odd(h))

		t := v & 15
		f = append(f, odd(t))
	}
	if len(f) < 16 {
		return [][]bool{}, errors.New("bit mapping.")
	}

	var forest [][]bool
	forest = append(forest, f[11:16])
	forest = append(forest, f[6:11])
	forest = append(forest, f[1:6])
	forest = append(forest, f[6:11])
	forest = append(forest, f[11:16])

	return forest, nil
}

func odd(i uint8) bool {
	if i % 2 == 0 {
		return true
	}
	return false
}

func b28rgb(by []byte) (color.RGBA, error) {
	if len(by) < 4 {
		return color.RGBA{}, errors.New("less than 4 byte.")
	}

	hue := float64(uint8(by[0]) + uint8(by[1])) * 360 / 4095
	sat := 65 - (float64(uint8(by[1])) * 20 / 255)
	lum := 75 - (float64(uint8(by[2])) * 20 / 255)

	cl := colorful.Hcl(hue, sat, lum)
	r, g, b := cl.RGB255()

	//rgba := hsl2rgba(hue, sat, lum)
	//return rgba, nil
	return color.RGBA{r, g, b, 255}, nil
}

func hsl2rgba(h, s, l float64) color.RGBA {
	if s == 0 {
		r := uint8(l * 255)
		g := uint8(l * 255)
		b := uint8(l * 255)
		return color.RGBA{r, g, b, 240}
	}

	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}
	v1 = 2*l - v2

	r := uint8(255 * hue2rgb(v1, v2, h + (1/3)))
	g := uint8(255 * hue2rgb(v1, v2, h))
	b := uint8(255 * hue2rgb(v1, v2, h - (1/3)))

	return color.RGBA{r, g, b, 255}
}

func hue2rgb(v1, v2, vH float64) float64 {
	if vH < 0 {
		vH += 1
	}
	if vH > 1 {
		vH -= 1
	}
	if (6 * vH) < 1 {
		return (v1 + (v2-v1) * 6 * vH)
	}
	if (2 * vH) < 1 {
		return v2
	}
	if (3 * vH) < 2 {
		return (v1 + (v2-v1) * ((2/3) - vH) * 6)
	}
	return v1
}

func array2img(f [][]bool, front color.RGBA, hw int64) (*image.RGBA, error) {
	cnt := int64(len(f) + 1)
	px := hw / cnt
	frame := px / 2
	if hw % cnt != 0 {
		frame += (hw % cnt) / 2
	}

	img := image.NewRGBA(image.Rect(0, 0, int(hw), int(hw)))
	setColor(img, 0, 0, hw, hw, BG_COLOR)

	for h, v := range f {
		for w, flg := range v {
			wpx := frame + (px * int64(w))
			hpx := frame + (px * int64(h))
			if !flg {
				continue
			}
			setColor(img, hpx, wpx, hpx+px, wpx+px, front)
		}
	}

	return img, nil
}

func setColor(img *image.RGBA, hs, ws, he, we int64, clr color.RGBA) {
	for h := hs; h <= he; h++ {
		for w := ws; w <= we; w++ {
			img.Set(int(h), int(w), clr)
		}
	}
}

func writter(fpath string, img *image.RGBA) error {
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	var size int
	var fpath string
	flag.IntVar(&size, "s", 64, "export image size.")
	flag.StringVar(&fpath, "e", "./idicon.png", "export image file path.")
	flag.Parse()

	if flag.NArg() < 1 {
		die("usage: iDicon [-s <size>] [-e <fpath>] <username>")
	}
	if flag.Arg(0) == "" {
		die("usage: iDicon [-s <size>] [-e <fpath>] <username>")
	}
	Seed = flag.Arg(0)

	if fpath == "" {
		die("empty export image file path.")
	}
	Fpath = fpath

	if size <= 0 {
		die("empty export image size.")
	}
	Size = int64(size)
}

func main() {
	if err := iDicon(); err != nil{
		die("failed: %s", err)
	}
}
