package main

import (
  "fmt"
  "image"
  "os"
  "github.com/brothertoad/bezier"
  "github.com/brothertoad/btu"
)

var xxxsvgPrefix =`<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="WIDTHpx" height="HEIGHTpx" viewBox="0 0 WIDTH HEIGHT" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink= "http://www.w3.org/1999/xlink">
  <style>
  * {
    fill: none;
    stroke: STROKE-COLOR;
    stroke-width: STROKE-WIDTH;
  }
  </style>
`
var xxxsvgSuffix = `</svg>
`

var svgPrefix =`<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="%dpx" height="%dpx" viewBox="0 0 %d %d" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink= "http://www.w3.org/1999/xlink">
  <style>
  * {
    fill: none;
    stroke: %s;
    stroke-width: %d;
  }
  </style>
`
var svgSuffix = `</svg>
`

var svgFile *os.File

func openSvg(path string) {
  svgFile = btu.CreateFile(path)
  writeSvgF(svgPrefix, mask.Global.Width, mask.Global.Height, mask.Global.Width, mask.Global.Height,
    mask.Global.StrokeColor, mask.Global.StrokeWidth)
}

func writeCurveToSvg(curve pointCollection, offset image.Point) {
  beziers := bezier.GetControlPointsI(curve.points)
  for _, bezier := range(beziers) {
    writeSvgF(`<path d="M %d %d `, bezier.P0.X + offset.X, bezier.P0.Y + offset.Y)
    writeSvgF(` C %d %d,`, bezier.P1.X + offset.X, bezier.P1.Y + offset.Y)
    writeSvgF(` %d %d,`, bezier.P2.X + offset.X, bezier.P2.Y + offset.Y)
    writeSvgF(` %d %d" fill="none"/>`, bezier.P3.X + offset.X, bezier.P3.Y + offset.Y)
    writeSvg("")  // to get a newline
  }
}

func writeBezierToSvg(bezier pointCollection, offset image.Point) {
  writeSvgF(`<path d="M %d %d `, bezier.points[0].X + offset.X, bezier.points[0].Y + offset.Y)
  writeSvgF(` C %d %d,`, bezier.points[1].X + offset.X, bezier.points[1].Y + offset.Y)
  writeSvgF(` %d %d,`, bezier.points[2].X + offset.X, bezier.points[2].Y + offset.Y)
  writeSvgF(` %d %d" fill="none"/>`, bezier.points[3].X + offset.X, bezier.points[3].Y + offset.Y)
  writeSvg("")  // to get a newline
}

func writeLineToSvg(line pointCollection, offset image.Point) {
  writeSvgF(`<polyline points="`)
  for j, p := range(line.points) {
    if j != 0 {
      writeSvgF(" ")
    }
    writeSvgF("%d,%d", p.X + offset.X, p.Y + offset.Y)
  }
  writeSvg(`" fill="none"/>`)
}

func writeSvg(s string) {
  svgFile.WriteString(s + "\n")
}

func writeSvgF(msg string, a ...any) {
  fmt.Fprintf(svgFile, msg, a...)
}

func closeSvg() {
  svgFile.WriteString(svgSuffix)
  err := svgFile.Close()
  btu.CheckError(err)
}
