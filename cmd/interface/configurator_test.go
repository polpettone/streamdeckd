package _interface

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const row0 = `
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

const row1 = `
line:
  - command: terminator -e "sudo systemctl restart iwd"
    icon: picture.png

  - command: deploy.sh
    icon: picture.png
`

func Test_UnmarshalRow(t *testing.T) {

	line, err := UnmarshalRow(row0)

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(line)
	if len(line.Keys) != 4 {
		t.Errorf("wanted %d got %d", 4, len(line.Keys))
	}
}

func TestDetectPages(t *testing.T) {
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dir)

	os.Mkdir(filepath.Join(dir, "page-1"), 0755)
	os.Mkdir(filepath.Join(dir, "page-2"), 0755)
	os.Mkdir(filepath.Join(dir, "page-3"), 0755)
	os.Mkdir(filepath.Join(dir, "page-12"), 0755)
	os.Mkdir(filepath.Join(dir, "invalid-3"), 0755)
	os.Mkdir(filepath.Join(dir, "page-3-invalid"), 0755)

	ioutil.WriteFile(filepath.Join(dir, "navi.yaml"), []byte("test"), 0644)

	pages, err := DetectPages(dir)
	if err != nil {
		t.Errorf("ReadDir() returned an error: %v", err)
	}

	if len(pages) != 4 {
		t.Errorf("wanted %d, got %d", 4, len(pages))
	}

	fmt.Printf("%v", pages)
}

func TestReadPages(t *testing.T) {
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dir)

	os.Mkdir(filepath.Join(dir, "page-1"), 0755)
	os.Mkdir(filepath.Join(dir, "page-2"), 0755)

	ioutil.WriteFile(filepath.Join(dir, "page-1", "row0.yaml"), []byte(row0), 0644)
	ioutil.WriteFile(filepath.Join(dir, "page-1", "row1.yaml"), []byte(row1), 0644)

	ioutil.WriteFile(filepath.Join(dir, "page-2", "row0.yaml"), []byte(row0), 0644)

	contents, err := ReadPages(dir, []int{1, 2})

	if len(contents) != 3 {
		t.Errorf("wanted %d got %d", 3, len(contents))
	}

}
