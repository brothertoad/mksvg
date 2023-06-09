package main

import (
  "fmt"
  "image"
  _ "image/jpeg"
  _ "image/png"
  "log"
  "os"
  "path/filepath"
  "github.com/brothertoad/btu"
)

func initFromImage(imagePath string) {
  w, h := getImageDimensions(imagePath)
  // Create new mask.html and mask.css files, and copy the image
  // file to the output directory.
  html := fmt.Sprintf(htmlTemplate, mask.Global.Title, w, h)
  path := filepath.Join(config.OutputDir, "mask.html")
  err := os.WriteFile(path, []byte(html), 0644)
  btu.CheckError(err)
  css := fmt.Sprintf(cssTemplate, w, h)
  path = filepath.Join(config.OutputDir, "mask.css")
  err = os.WriteFile(path, []byte(css), 0644)
  btu.CheckError(err)
  path = filepath.Join(config.OutputDir, "mask.jpg")
  btu.CopyFile(imagePath, path)
}

// This was copied from https://gist.github.com/sergiotapia/7882944
func getImageDimensions(imagePath string) (int, int) {
    file := btu.OpenFile(imagePath)
    image, _, err := image.DecodeConfig(file)
    if err != nil {
        log.Fatalln(err)
    }
    return image.Width, image.Height
}
