// https://www.bittorrent.org/beps/bep_0003.html
// https://en.wikipedia.org/wiki/Torrent_file
package metainfo

import (
	"bufio"
	"fmt"
	"io"
	"torora/pkg/bencoding"
)

type Metainfo struct {
	Announce string // announce
	Info     Info   // info
}

type Info struct {
	Name        string // name
	PieceLength int64  // piece length
	Pieces      string // pieces
	Length      int64  // length
	Files       []File // files
}

type File struct {
	Length int64  // length
	Path   string // path
}

func Parse(r io.Reader) (Metainfo, error) {
	br := bufio.NewReader(r)
	decoder := bencoding.NewDecoder(br)

	raw, err := decoder.Decode()
	if err != nil {
		return Metainfo{}, err
	}

	return parseMetainfo(raw)
}

func parseMetainfo(raw interface{}) (Metainfo, error) {
	var result Metainfo
	var ok bool

	rawMetainfo, ok := raw.(map[string]interface{})
	if !ok {
		return Metainfo{}, fmt.Errorf("metainfo is not dict")
	}

	if _, ok = rawMetainfo["announce"]; !ok {
		return Metainfo{}, fmt.Errorf("metainfo missing announce")
	}
	if result.Announce, ok = rawMetainfo["announce"].(string); !ok {
		return Metainfo{}, fmt.Errorf("metainfo/announce is not string")
	}

	var err error
	if result.Info, err = parseInfo(rawMetainfo["info"]); err != nil {
		return Metainfo{}, err
	}

	return result, nil
}

func parseInfo(raw interface{}) (Info, error) {
	var result Info
	var ok bool

	rawInfo, ok := raw.(map[string]interface{})
	if !ok {
		return Info{}, fmt.Errorf("info is not dict")
	}

	if _, ok = rawInfo["name"]; !ok {
		return Info{}, fmt.Errorf("info missing name")
	}
	if result.Name, ok = rawInfo["name"].(string); !ok {
		return Info{}, fmt.Errorf("info/name is not string")
	}

	if _, ok = rawInfo["piece length"]; !ok {
		return Info{}, fmt.Errorf("info missing piece length")
	}
	if result.PieceLength, ok = rawInfo["piece length"].(int64); !ok {
		return Info{}, fmt.Errorf("info/piece length is not int64")
	}

	if _, ok = rawInfo["pieces"]; !ok {
		return Info{}, fmt.Errorf("info missing pieces")
	}
	if result.Pieces, ok = rawInfo["pieces"].(string); !ok {
		return Info{}, fmt.Errorf("info/pieces is not string")
	}

	// exist key length or a key files, but not both or neither
	_, existLength := rawInfo["length"]
	_, existFiles := rawInfo["files"]
	if existLength == existFiles {
		return Info{}, fmt.Errorf("info/length exist same as info/files")
	}

	if existLength {
		if result.Length, ok = rawInfo["length"].(int64); !ok {
			return Info{}, fmt.Errorf("info/length is not int64")
		}
	} else {
		rawFiles, ok := rawInfo["files"].([]interface{})
		if !ok {
			return Info{}, fmt.Errorf("info/files is not list")
		}

		files := make([]File, 0, len(rawFiles))
		for _, rawFile := range rawFiles {
			file, err := parseFile(rawFile)
			if err != nil {
				return Info{}, err
			}
			files = append(files, file)
		}
		result.Files = files
	}

	return result, nil
}

func parseFile(raw interface{}) (File, error) {
	var result File
	var ok bool

	rawFile, ok := raw.(map[string]interface{})
	if !ok {
		return File{}, fmt.Errorf("file is not dict")
	}

	if _, ok = rawFile["length"]; !ok {
		return File{}, fmt.Errorf("file missing length")
	}
	if result.Length, ok = rawFile["length"].(int64); !ok {
		return File{}, fmt.Errorf("file/length is not int64")
	}

	if _, ok = rawFile["path"]; !ok {
		return File{}, fmt.Errorf("file missing path")
	}
	if result.Path, ok = rawFile["path"].(string); !ok {
		return File{}, fmt.Errorf("file/path is not string")
	}

	return result, nil
}
