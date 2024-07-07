package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadThumbnailImage(httpClient http.Client, baseURL string, i, j float64, setting Setting) error {
	u := fmt.Sprintf(baseURL, int(i), int(j))
	httpRequest, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  "pan_session",
		Value: setting.SessionID,
	}

	httpClient.Jar.SetCookies(httpRequest.URL, []*http.Cookie{cookie})

	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fileName := fmt.Sprintf("thumbnails/save_%03d_%03d.jpg", int(i), j)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
