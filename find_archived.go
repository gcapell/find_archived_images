package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"log"
	"fmt"
	"flag"
	"time"
	"path"
)

var (
	archive = flag.String("archive", "/Volumes/Pictures", "top-level archive directory")
)

func checkCopied(src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	srcSize := fi.Size()

	x, err := exif.Decode(f)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		log.Printf("err %s closing %q", err, src)
	}

	date, err := x.Get(exif.DateTimeOriginal)
	if err != nil {
		return err
	}
	t, err := time.Parse("2006:01:02 15:04:05", date.StringVal())
	if err != nil {
		return err
	}
	_, base := path.Split(src)
	dst := path.Join(*archive, t.Format("2006/01/02"), base)

	dstF, err := os.Open(dst)
	if err != nil {
		return err
	}
	fi, err = dstF.Stat()
	dstSize := fi.Size()

	if dstSize != srcSize {
		return fmt.Errorf("src:%q (%d bytes), dst:%q (%d bytes)", src, srcSize, dst, dstSize)
	}
	log.Printf("%q->%q", src, dst)
	return nil
}

func main() {
	flag.Parse()
	fname := flag.Arg(0)
	err := checkCopied(fname)
	if err != nil {
		log.Println(fname, err)
		return
	}
	log.Println("safe to delete", fname)
}
