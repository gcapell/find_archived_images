package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"log"
	"fmt"
	"flag"
	"time"
	"path"
	"path/filepath"
)

var (
	archive = flag.String("archive", "/Volumes/Pictures", "top-level archive directory")
	srcDir = flag.String("d", ".", "Directory to search for images")
)

func checkCopied(src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	srcSize := fi.Size()

	x, err := exif.Decode(f)
	if err != nil {
		return err
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
	defer dstF.Close()
	fi, err = dstF.Stat()
	dstSize := fi.Size()

	if dstSize != srcSize {
		return fmt.Errorf("src:%q (%d bytes), dst:%q (%d bytes)", src, srcSize, dst, dstSize)
	}
	// log.Printf("%q->%q", src, dst)
	return nil
}

func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("walker", path, err)
		return nil
	}
	if info.IsDir() {
		return nil
	}
	err = checkCopied(path)
	if err != nil {
		log.Println(path, err)
		return nil
	}
	err = os.Remove(path)
	if err != nil {
		log.Println("removing", path, err)
		return err	// worth stopping because this is weird
	}
	log.Println("deleted", path)
	return nil
}

func main() {
	flag.Parse()
	
	filepath.Walk(*srcDir, walker)
}
