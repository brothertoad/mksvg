package main

import (
  "fmt"
  "image"
  _ "image/jpeg"
  _ "image/png"
  "os"
  "path/filepath"
  "github.com/brothertoad/btu"
)

func initFromImage(imagePath string) {
  w, h := getImageDimensions(imagePath)
  // Create new mask.html and mask.css files, and copy the image
  // file to the output directory.
  html := fmt.Sprintf(htmlTemplate, mask.Global.Title)
  path := filepath.Join(config.OutputDir, "mask.html")
  err := os.WriteFile(path, []byte(html), 0644)
  btu.CheckError2(err, "unable to write html file '%s'", path)
  css := fmt.Sprintf(cssTemplate, w, h)
  path = filepath.Join(config.OutputDir, "mask.css")
  err = os.WriteFile(path, []byte(css), 0644)
  btu.CheckError2(err, "Unable to write css file '%s'", path)
  path = filepath.Join(config.OutputDir, "mask.jpg")
  btu.CopyFile(imagePath, path)
}

// This was copied from https://gist.github.com/sergiotapia/7882944
func getImageDimensions(imagePath string) (int, int) {
    file := btu.OpenFile(imagePath)
    image, _, err := image.DecodeConfig(file)
    btu.CheckError2(err, "Can't get dimensions from file '%s'", imagePath)
    return image.Width, image.Height
}
