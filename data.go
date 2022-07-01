package main

import (
  "image"
  "image/jpeg"
  "io/ioutil"
  "os"
  "github.com/pelletier/go-toml"
  "github.com/brothertoad/btu"
)

type GlobalStruct struct {
  Image string
  StrokeColor string
  StrokeWidth int
  Title string
  PrintName string
}

type ObjectStruct struct {
  Curves []string
  Lines []string
}

var mask struct {
  Global GlobalStruct
  Points map[string]image.Point
  Objects map[string]ObjectStruct
  image image.Image
  width int
  height int
}

func loadMask(path string) {
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  err = toml.Unmarshal(b, &mask)
    btu.CheckError(err)
  if mask.Global.PrintName == "" {
    mask.Global.PrintName = mask.Global.Title
  }
  loadImage(mask.Global.Image)
}

func loadImage(path string) {
  reader, err := os.Open(path)
  btu.CheckError(err)
  defer reader.Close();
  mask.image, err = jpeg.Decode(reader)
    btu.CheckError(err)
  mask.width = mask.image.Bounds().Max.X - mask.image.Bounds().Min.X
  mask.height = mask.image.Bounds().Max.Y - mask.image.Bounds().Min.Y
}
