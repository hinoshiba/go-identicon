package main

import (
	"os"
	"log"
	"fmt"
	"flag"
)

func die(s string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, s + "\n" , msg...)
	os.Exit(1)
}

var Seed  string
var SizeH int64
var SizeW int64
var Fpath string

type RGB struct {
	r uint8
	g uint8
	b uint8
}

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

	img, err := array2img(forest, rgb, SizeH, SizeW)
	if err != nil {
		return err
	}
	if err := writter(Fpath, img); err != nil {
		return err
	}
	return nil
}

func seed2bits(seed string) ([]byte, []byte, error) {
	log.Println(seed)
	return []byte(""), []byte(""), nil
}

func b60forest(b []byte) ([][]bool, error) {
	log.Println(b)
	return [][]bool{}, nil
}

func b28rgb(b []byte) (RGB, error) {
	log.Println(b)
	return RGB{}, nil
}

func array2img(f [][]bool, rgb RGB, h int64, w int64) ([]byte, error) {
	log.Println(f, rgb, h, w)
	return []byte(""), nil
}

func writter(fpath string, img []byte) error {
	log.Println(fpath, img)
	return nil
}

func init() {
	var size string
	var fpath string
	flag.StringVar(&size, "s", "64x64", "export image size.")
	flag.StringVar(&fpath, "e", "./idicon", "export image file path.")
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

	if size == "" {
		die("empty export image size.")
	}

	h, ok := SIZES_H[size]
	if !ok {
		die("undefined size : %s.", size)
	}
	SizeH = h

	w, ok := SIZES_W[size]
	if !ok {
		die("undefined size : %s.", size)
	}
	SizeW = w
}

func main() {
	if err := iDicon(); err != nil{
		die("failed: %s", err)
	}
}
