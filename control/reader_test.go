// Copyright (c) 2011 Mikkel Krautz
// The use of this source code is goverened by a BSD-style
// license that can be found in the LICENSE-file.

package control

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

type readerTestCase struct {
	Name     string
	Input    string
	Expected []KeyValuePair
}

var readerTests []readerTestCase = []readerTestCase{
	{
		"single line",

`A: 0
B: 1
C: 2`,

		[]KeyValuePair{ { "A", "0" }, { "B", "1" }, { "C", "2" }, },
	},
	{
		"multi-line",
	
`Test: This text
 spans multiple lines`,

 		[]KeyValuePair{ { "Test", `This text
spans multiple lines` }, },
	},
	{
		"trimming",
		`A: are we getting trimmed?      `,
		[]KeyValuePair{ { "A", "are we getting trimmed?" }, },
	},
	{
		"trimming-comment",
		"A: are we getting trimmed? # why yes, yes of course!",
		[]KeyValuePair{ { "A", "are we getting trimmed?" }, },
	},
}


func cmpKeyValuePair(a []KeyValuePair, b []KeyValuePair) (bool, os.Error) {
	if len(a) != len(b) {
		return false, fmt.Errorf("length mismatch (a=%v, b=%v)", len(a), len(b))
	}
	for i := 0; i < len(a); i++ {
		if a[i].Key != b[i].Key {
			return false, fmt.Errorf("key: %v != %v", a[i].Key, b[i].Key)
		}
		if a[i].Value != b[i].Value {
			return false, fmt.Errorf("value: %v != %v", a[i].Value, b[i].Value)
		}
	}
	return true, nil
}

func TestEverything(t *testing.T) {
	for _, testCase := range readerTests {
		kvp, err := Parse(bytes.NewBufferString(testCase.Input))	
		if err != nil {
			t.Fatal("%v failed: %v", testCase.Name, err)
		}
		if match, _ := cmpKeyValuePair(kvp, testCase.Expected); !match {
			t.Fatalf("%v failed: kvp mismatch\nexpected=%v\ngot=%v", testCase.Name, testCase.Expected, kvp)
		}
	}
}

// Test that we handle a comment in the middle of a key
func TestKeycomment(t *testing.T) {
	_, err := Parse(bytes.NewBufferString("KeyComment#: Blah"))
	if err == nil {
		t.Fatal("expected malformed file error")
	}
}

var nautilusDropboxExpected []KeyValuePair = []KeyValuePair{
	{"Package", "nautilus-dropbox"},
	{"Version", "0.6.9"},
	{"Architecture", "amd64"},
	{"Maintainer", "Rian Hunter <rian@dropbox.com>"},
	{"Installed-Size", "460"},
	{"Depends", "libatk1.0-0 (>= 1.20.0), libc6 (>= 2.4), libcairo2 (>= 1.6.0), libglib2.0-0 (>= 2.16.0), libgtk2.0-0 (>= 2.12.0), libnautilus-extension1 (>= 1:2.22.2), libpango1.0-0 (>= 1.20.1), python (>= 2.5), python-gtk2 (>= 2.12)"},
	{"Suggests", "nautilus (>= 2.16.0)"},
	{"Section", "gnome"},
	{"Priority", "optional"},
	{"Description", `Dropbox integration for Nautilus
Nautilus Dropbox is an extension that integrates
the Dropbox web service with your GNOME Desktop.
.
Check us out at http://www.dropbox.com/`},
}

// Test that we can read a real control file (nautilus-dropbox)
func TestReadRealControlFile(t *testing.T) {
	f, err := os.Open("testdata/control-nautilus-dropbox")	
	if err != nil {
		t.Fatal(err)
	}

	kvp, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	match, err := cmpKeyValuePair(kvp, nautilusDropboxExpected)
	if !match {
		t.Fatalf("TestReadrealControlFile failed: %v", err)
	}
}
