package main

import (
  "image"
)

// This file containts types, constants, vars, etc.
// that are used across the app.

// This structure holds values from the configuration file
// and values from command line arguments.
var config struct {
  OutputDir string
  PointRadius int
  MarginEdge int
  StrokeColor string
  StrokeWidth int
  inputPath string
  outputPath string
  printPoints bool
  printBorder bool
}

type GlobalInfo struct {
  Image string
  StrokeColor string
  StrokeWidth int
  Title string
  PrintName string
  Scale float64
  Width int
  Height int
}

type pointCollection struct {
    points []image.Point
  }

type RenderObject struct {
  Object    string
  Comment   string
  Hide      bool
  Translate image.Point
  Scale     float64
  Flip      string
}

// The fields prefixed with "raw" are computed from the
// corresponding non-raw fields.
type Object struct {
  Curves []string
  Beziers []string
  Lines []string
  Rects []string
  rawCurves []pointCollection
  rawBeziers []pointCollection
  rawLines []pointCollection
  rawRects []image.Rectangle
  center image.Point
}

// Note that some fields in this object are read directly from the input file,
// whereas others are computed.
var mask struct {
  Global GlobalInfo
  Points map[string]image.Point
  Objects map[string]Object
  Renders []RenderObject
}
