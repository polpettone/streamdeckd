package _interface

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/polpettone/streamdeckd/cmd/models"
	"github.com/unix-streamdeck/api"
	"gopkg.in/yaml.v2"
)

const PAGE_NAME_PATTERN = "page-"

func UnmarshalRow(raw string) (*models.PageRow, error) {
	var rows models.PageRow
	err := yaml.Unmarshal([]byte(raw), &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}

func SetupConfigurationFromDir(dirPath string) (*api.Config, error) {
	config := &api.Config{}

	pages, err := DetectPages(dirPath)
	if err != nil {
		return nil, err
	}

	rawContents, err := ReadPages(dirPath, pages)
	if err != nil {
		return nil, err
	}

	max := pages[0]
	for _, v := range pages {
		if v > max {
			max = v
		}
	}

	configPages := []api.Page{}
	for n := 0; n <= max; n++ {
		configPages = append(configPages, api.Page{})
	}

	for _, raw := range rawContents {
		index := raw.PageNumber
		pageRaw, err := UnmarshalRow(raw.Content)
		if err != nil {
			return nil, err
		}
		configPages[index] = append(configPages[index], pageRaw.Keys...)
	}

	deck := api.Deck{
		Serial: "CL33L2A02177",
		Pages:  configPages,
	}

	config.Decks = []api.Deck{deck}

	return config, nil
}

func ReadPages(dirPath string, pages []int) ([]PageRawContent, error) {
	log.Printf("Read Pages: %s", dirPath)
	pageRawContents := []PageRawContent{}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i] < pages[j]
	})

	for _, page := range pages {
		pageDirName := fmt.Sprintf("%s/%s%d", dirPath, PAGE_NAME_PATTERN, page)
		entries, err := os.ReadDir(pageDirName)
		if err != nil {
			return nil, err
		}

		sort.Slice(entries, func(i, j int) bool {
			if strings.Compare(entries[i].Name(), entries[j].Name()) > 0 {
				return false
			}
			return true
		})

		for i, entry := range entries {
			if !entry.IsDir() {
				content, err := os.ReadFile(filepath.Join(pageDirName, entry.Name()))
				pageRawContents = append(
					pageRawContents,
					PageRawContent{
						PageNumber: page,
						RowNumber:  i,
						Content:    string(content)})

				if err != nil {
					return nil, err
				}
			}
		}
	}
	return pageRawContents, nil
}

type PageRawContent struct {
	RowNumber  int
	Content    string
	PageNumber int
}

func DetectPages(dir string) ([]int, error) {
	log.Printf("Detect pages: %s", dir)
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	var numbers []int

	for _, entry := range entries {
		entryType := "File"
		if entry.IsDir() {
			number, err := extractPageNumber(entry.Name())
			if err == nil {
				numbers = append(numbers, number)
			}
		}
		fmt.Printf("%s: %s\n", entry.Name(), entryType)
	}
	return numbers, nil
}

func extractPageNumber(s string) (int, error) {

	if !strings.HasPrefix(s, PAGE_NAME_PATTERN) {
		return 0, errors.New("invalid prefix")
	}

	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return 0, errors.New("invalid input format")
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("no number after -")
	}
	return num, nil
}
