package api

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

// ExtractImageMetadata parses EXIF and PNG text chunks from the image at path.
// The path argument must reference a readable JPEG or PNG on disk; the function
// opens the file, extracts known metadata fields, and returns them as a map
// without modifying the underlying file.
func ExtractImageMetadata(path string) (map[string]interface{}, error) {
	meta := make(map[string]interface{})
	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".jpg" || ext == ".jpeg" {
		f, err := os.Open(path)
		if err != nil {
			return meta, err
		}
		defer f.Close()
		x, err := exif.Decode(f)
		if err == nil {
			if tag, err := x.Get(exif.UserComment); err == nil {
				if s, err := tag.StringVal(); err == nil && s != "" {
					meta["UserComment"] = s
				}
			}
			if tag, err := x.Get(exif.ImageDescription); err == nil {
				if s, err := tag.StringVal(); err == nil && s != "" {
					meta["ImageDescription"] = s
				}
			}
		}
	} else if ext == ".png" {
		f, err := os.Open(path)
		if err != nil {
			return meta, err
		}
		defer f.Close()
		// skip PNG signature
		sig := make([]byte, 8)
		if _, err := f.Read(sig); err != nil {
			return meta, err
		}
		for {
			var length uint32
			if err := binary.Read(f, binary.BigEndian, &length); err != nil {
				break
			}
			chunkType := make([]byte, 4)
			if _, err := f.Read(chunkType); err != nil {
				break
			}
			data := make([]byte, length)
			if _, err := f.Read(data); err != nil {
				break
			}
			// skip CRC
			if _, err := f.Read(make([]byte, 4)); err != nil {
				break
			}
			t := string(chunkType)
			if t == "tEXt" || t == "iTXt" {
				parts := strings.SplitN(string(data), "\x00", 2)
				if len(parts) == 2 {
					key := parts[0]
					val := parts[1]
					if key != "" && val != "" {
						meta[key] = val
					}
				}
			}
		}
	}

	return meta, nil
}
