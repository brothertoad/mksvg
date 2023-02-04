package main

import (
  "image"
  _ "image/jpeg"
  "io/ioutil"
  "github.com/pelletier/go-toml"
  "github.com/brothertoad/btu"
)

type GlobalInfo struct {
  Image string
  StrokeColor string
  StrokeWidth int
  Title string
  PrintName string
  Width int
  Height int
}

type InputObject struct {
  Curves []string
  Beziers []string
  Lines []string
  Rects []string
}

type curveInfo struct {
  points []image.Point
}

type bezierInfo struct {
  points [4]image.Point
}

type lineInfo struct {
  points []image.Point
}

// Note that some fields in this object are read directly from the input file,
// whereas others are computed.
var mask struct {
  Global GlobalInfo
  Points map[string]image.Point
  InputObjects map[string]InputObject
  curves []curveInfo
  beziers []bezierInfo
  lines []lineInfo
  rects []image.Rectangle
}

func loadMask(path string) {
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  err = toml.Unmarshal(b, &mask)
  btu.CheckError(err)
  if mask.Global.PrintName == "" {
    mask.Global.PrintName = mask.Global.Title
  }
  // If there are not objects, there is nothing to do.
  if mask.InputObjects == nil {
    btu.Fatal("No objects in mask file.\n")
  }
  // If there are no points defined, make Points an empty slice to ease error checking.
  if mask.Points == nil {
    mask.Points = make(map[string]image.Point)
  }
  // If a list of items is nil, make it an empty slice to ease error checking.
  for _, obj := range(mask.InputObjects) {
    if obj.Curves == nil {
      obj.Curves = make([]string, 0)
    }
    if obj.Beziers == nil {
      obj.Beziers = make([]string, 0)
    }
    if obj.Lines == nil {
      obj.Lines = make([]string, 0)
    }
    if obj.Rects == nil {
      obj.Rects = make([]string, 0)
    }
  }
}
