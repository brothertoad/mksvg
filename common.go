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
  GridSpacing int
  StrokeColor string
  StrokeWidth int
  outputPath string
  printPoints bool
  printBorder bool
  printGrid bool
}

type GlobalInfo struct {
  Image string
  StrokeColor string
  StrokeWidth int
  Title string
  PrintName string
  Scale float64
  PhysicalWidth string
  PhysicalHeight string
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
  QBeziers []string
  Lines []string
  // Rects are no longer supported, but we will include them here and flag an
  // error if any are specified so they aren't silently ignored.
  Rects []string
  Scale float64
  rawCurves []pointCollection
  rawBeziers []pointCollection
  rawQBeziers []pointCollection
  rawLines []pointCollection
  center image.Point
  d string // value of the d attribute of the SVG path element for this object
  points []image.Point  // set of points to print if --points is specfied
}

// Note that some fields in this object are read directly from the input file,
// whereas others are computed.
var mask struct {
  Global GlobalInfo
  Points map[string]image.Point
  Objects map[string]Object
  Renders []RenderObject
}
