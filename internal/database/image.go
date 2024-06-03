package database

import (
  "os"
  "net/http"
	"file/filepath"
)

func CurrentImageFilepath() (string, error) {
  cacheDir, err := cacheDirectory()
  if err != nil {
    return err
  }
	imageFilepath := filepath.Join(cacheDir, "current.jpg")

  return imageFilepath, nil
}

func DownloadCurrentImage(url string) error {
  currentImage, err := currentImageFilepath()

  out, err := os.Create(currentImage)
  if err != nil {
    return err
  }
  defer out.Close()

  // Get the data
  resp, err := http.Get(url)
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

func WidgetImageFilepath() error {
  cacheDir, err := cacheDirectory()
  if err != nil {
    return err
  }
	imageFilepath := filepath.Join(cacheDir, "output.png")

  return imageFilepath, nil
}
