package database

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DownloadCurrentImage(url string) error {
	currentImage, err := CurrentImageFilepath()
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Clean(currentImage))
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	validUrl, err := validateURL(url)
	if err != nil {
		return err
	}

	resp, err := http.Get(validUrl) // #nosec G107
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func CurrentImageFilepath() (string, error) {
	cacheDir, err := cacheDirectory()
	if err != nil {
		return "", err
	}
	imageFilepath := filepath.Join(cacheDir, "assets", "current.jpg")

	return imageFilepath, nil
}

func WidgetImageFilepath() (string, error) {
	cacheDir, err := cacheDirectory()
	if err != nil {
		return "", err
	}
	imageFilepath := filepath.Join(cacheDir, "output.png")

	return imageFilepath, nil
}

func validateURL(rawURL string) (string, error) {
	allowedDomains := []string{"api.spotify.com", "i.scdn.co", "open.spotify"}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("invalid URL")
	}

	// Ensure the URL scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", errors.New("invalid URL scheme")
	}

	// Check if the URL host is in the list of allowed domains
	for _, domain := range allowedDomains {
		if strings.HasSuffix(parsedURL.Host, domain) {
			return parsedURL.String(), nil
		}
	}

	return "", errors.New("URL is not allowed")
}
