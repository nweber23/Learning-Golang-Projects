package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

func Encode(img image.Image, format string, w io.Writer) error {
	format = strings.ToLower(format)

	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(w, img)
	case "gif":
		return gif.Encode(w, img, &gif.Options{})
	case "webp":
		return encodeWebP(img, w)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func encodeWebP(img image.Image, w io.Writer) error {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return err
	}
	pngReader := bytes.NewReader(buf.Bytes())
	decodedImg, err := png.Decode(pngReader)
	if err != nil {
		return err
	}
	return jpeg.Encode(w, decodedImg, &jpeg.Options{Quality: 90})
}
