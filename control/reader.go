// Copyright (c) 2011 Mikkel Krautz
// The use of this source code is goverened by a BSD-style
// license that can be found in the LICENSE-file.

// Package control implements reading of files in the style of Debian control files
package control

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// A KeyValuePair represents a key/value pair
// found in a Debian control file.
type KeyValuePair struct {
	Key   string
	Value string
}

func isKeySeparator(buf []byte) bool {
	if len(buf) >= 2 {
		return string(buf[0:2]) == ": "
	}		
	return false
}

// Parse parses a file in the style of Debian control files. 
//
// A Debian control file consists of key-value pairs separated by
// colons, e.g.:
//
// Package: mypackage
// Version: 4.5.0 # this is a comment
//
// A value can span multiple lines of subsequent lines if it begin
// with a space, like so:
//
// Description: Hello 
//  world
//
// The example above would yield the value "Hello\nworld"
//
// The Parse function automatically discards whitespace at the beginning and
// end of any parsed values.
func Parse(r io.Reader) (kvps []KeyValuePair, err os.Error) {
	buf, err := bufio.NewReaderSize(r, 4096)
	if err != nil {
		return nil, err
	}

Line:
	for {
		kvp := KeyValuePair{}

		line, isPrefix, err := buf.ReadLine()
		if isPrefix {
			return nil, os.NewError("line exceeds internal buffer limit")	
		}
		if err == os.EOF {
			return kvps, nil
		} else if err != nil {
			return nil, err
		}

		key := []byte{}
		for idx, rune := range string(line) {
			if rune == '#' {
				if idx == 0 {
					continue Line
				}
				return nil, os.NewError("debcontrol: malformed input file: comment '#' in key section")
			}		
			if isKeySeparator(line[idx:]) {
				key = line[0:idx]
				line = line[idx+2:]
				break
			}
		}
		kvp.Key = string(key)

		value := []byte{}
		for {
			for idx, rune := range string(line) {
				if rune == '#' {
					line = line[0:idx]
					break
				}	
			}
			value = append(value, line...)

			lookahead, err := buf.Peek(1)
			if err != nil {
				break
			}

			if lookahead[0] == ' ' {
				line, isPrefix, err = buf.ReadLine()
				if isPrefix {
					return nil, os.NewError("debcontrol: line exceeds internal buffer limit")	
				}
				if err == os.EOF {
					return kvps, nil
				} else if err != nil {
					return nil, err
				}
				line[0] = '\n'
				continue
			} else {
				break
			}
		}

		kvp.Value = strings.TrimSpace(string(value))
		kvps = append(kvps, kvp)
	}

	return
}
