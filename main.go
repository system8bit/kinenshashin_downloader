package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Setting struct {
	SessionID string  `yaml:"session_id"`
	PhotoID   string  `yaml:"photo_id"`
	MaxX      float64 `yaml:"max_x"`
	MaxY      float64 `yaml:"max_y"`
	StepX     float64 `yaml:"step_x"`
	StepY     float64 `yaml:"step_y"`
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func main() {
	yamlFile := "./setting.yaml"

	flag.Parse()
	flag.String("yaml", yamlFile, "yaml file")

	bytes, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	setting := Setting{}
	err = yaml.Unmarshal(bytes, &setting)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = downloadThumbnailAll(setting)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = combineAllImage(setting, "out.png")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func downloadThumbnailAll(setting Setting) error {
	stepX := setting.StepX
	stepY := setting.StepY
	jar := NewJar()
	httpClient := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       0,
	}
	var baseURL = fmt.Sprintf("https://pan.kinenshashin.net/list/zoom/%s?preview_width=149&preview_height=224&preview_x=%%d&preview_y=%%d", setting.PhotoID)
	for i := 0.0; i < setting.MaxX; i += stepX {
		for j := 0.0; j < setting.MaxY; j += stepY {
			err := downloadThumbnailImage(httpClient, baseURL, i, j, setting)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func combineAllImage(setting Setting, output string) error {
	stepX := setting.StepX
	stepY := setting.StepY
	filenameMatrix := make([][]string, 0)
	for i := 0.0; i < 144; i += stepX {
		filenames := make([]string, 0)
		for j := 0; j < 213; j += int(stepY) {
			fileName := fmt.Sprintf("thumbnails/save_%03d_%03d.jpg", int(i), j)
			filenames = append(filenames, fileName)
		}
		filenameMatrix = append(filenameMatrix, filenames)
	}

	if err := concatImageFiles(output, filenameMatrix); err != nil {
		return err
	}

	return nil
}
