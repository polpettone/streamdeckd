package _interface

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_UnmarshalRow(t *testing.T) {

	input :=
		`
line:
  - key_handler: IconState
    icon_handler: IconState
    icon_handler_fields: {
      "text_1" : "A",
      "icon_1" :"picture.jpg",
    }

  -

  - command: terminator -e "sudo systemctl restart iwd"
    icon: picture.png

  - command: deploy.sh
    icon: picture.png
`

	line, err := UnmarshalRow(input)

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(line)
	if len(line.Keys) != 4 {
		t.Errorf("wanted %d got %d", 4, len(line.Keys))
	}
}

func TestReadDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dir) // Bereinige das temporäre Verzeichnis am Ende des Tests

	os.Mkdir(filepath.Join(dir, "page-1"), 0755)
	os.Mkdir(filepath.Join(dir, "page-2"), 0755)
	os.Mkdir(filepath.Join(dir, "page-3"), 0755)
	os.Mkdir(filepath.Join(dir, "invalid-3"), 0755)
	os.Mkdir(filepath.Join(dir, "page-3-invalid"), 0755)

	ioutil.WriteFile(filepath.Join(dir, "navi.yaml"), []byte("test"), 0644)

	pages, err := DetectPages(dir)
	if err != nil {
		t.Errorf("ReadDir() returned an error: %v", err)
	}

	if len(pages) != 3 {
		t.Errorf("wanted %d, got %d", 3, len(pages))
	}

	fmt.Printf("%v", pages)

}
