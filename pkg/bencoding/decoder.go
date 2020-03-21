// https://en.wikipedia.org/wiki/Bencode
package bencoding

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r *bufio.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

func (d *Decoder) Decode() (interface{}, error) {
	start, err := d.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	if start == 'e' {
		if err := d.r.UnreadByte(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	if start == 'l' {
		return d.decodeList()
	}
	if start == 'd' {
		return d.decodeDict()
	}
	if start == 'i' {
		return d.decodeInt64('e')
	}
	if unicode.IsDigit(rune(start)) {
		if err := d.r.UnreadByte(); err != nil {
			return nil, err
		}
		return d.decodeString()
	}

	return nil, fmt.Errorf("unimplement %s", string(start))
}

func (d *Decoder) decodeInt64(until byte) (int64, error) {
	rawResult, err := d.r.ReadString(until)
	if err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("integer missing %s", string(until))
		}
		return 0, err
	}

	// strip down until
	rawResult = rawResult[:len(rawResult)-1]
	if len(rawResult) == 0 {
		return 0, fmt.Errorf("integer missing value")
	}

	result, err := strconv.ParseInt(rawResult, 10, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (d *Decoder) decodeString() (string, error) {
	lenResult, err := d.decodeInt64(':')
	if err != nil {
		return "", err
	}

	result := ""

	var temp int64
	for temp = 0; temp < lenResult; temp += 1 {
		b, err := d.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		result += string(b)
	}
	if temp != lenResult {
		return "", fmt.Errorf("string missing character")
	}

	return result, nil
}

func (d *Decoder) decodeList() ([]interface{}, error) {
	var result []interface{}

	for {
		element, err := d.Decode()
		if err != nil {
			return nil, err
		}
		if element == nil {
			break
		}

		result = append(result, element)
	}

	expect, err := d.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("list missing e")
		}
		return nil, err
	}
	if expect != 'e' {
		return nil, fmt.Errorf("list missing e")
	}

	return result, nil
}

func (d *Decoder) decodeDict() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	prevKey := ""
	for {
		rawKey, err := d.Decode()
		if err != nil {
			return nil, err
		}
		if rawKey == nil {
			break
		}

		key, ok := rawKey.(string)
		if !ok {
			return nil, fmt.Errorf("dict key must be string")
		}

		if strings.Compare(prevKey, key) == 1 {
			return nil, fmt.Errorf("dict key must appear in lexicographical order")
		}
		prevKey = key

		value, err := d.Decode()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, fmt.Errorf("dict key missing value")
		}

		result[key] = value
	}

	expect, err := d.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("dict missing e")
		}
		return nil, err
	}
	if expect != 'e' {
		return nil, fmt.Errorf("dict missing e")
	}

	return result, nil
}
