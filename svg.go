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

func writeCurveToSvg(curve pointCollection, center, translate image.Point) {
  xoffset := translate.X - center.X
  yoffset := translate.Y - center.Y
  beziers := bezier.GetControlPointsI(curve.points)
  for _, bezier := range(beziers) {
    writeSvgF(`<path d="M %d %d `, bezier.P0.X + xoffset, bezier.P0.Y + yoffset)
    writeSvgF(` C %d %d,`, bezier.P1.X + xoffset, bezier.P1.Y + yoffset)
    writeSvgF(` %d %d,`, bezier.P2.X + xoffset, bezier.P2.Y + yoffset)
    writeSvgF(` %d %d" fill="none"/>`, bezier.P3.X + xoffset, bezier.P3.Y + yoffset)
    writeSvg("")  // to get a newline
  }
}

func writeBezierToSvg(bezier pointCollection, center, translate image.Point) {
  xoffset := translate.X - center.X
  yoffset := translate.Y - center.Y
  writeSvgF(`<path d="M %d %d `, bezier.points[0].X + xoffset, bezier.points[0].Y + yoffset)
  writeSvgF(` C %d %d,`, bezier.points[1].X + xoffset, bezier.points[1].Y + yoffset)
  writeSvgF(` %d %d,`, bezier.points[2].X + xoffset, bezier.points[2].Y + yoffset)
  writeSvgF(` %d %d" fill="none"/>`, bezier.points[3].X + xoffset, bezier.points[3].Y + yoffset)
  writeSvg("")  // to get a newline
}

func writeLineToSvg(line pointCollection, center, translate image.Point) {
  xoffset := translate.X - center.X
  yoffset := translate.Y - center.Y
  writeSvgF(`<polyline points="`)
  for j, p := range(line.points) {
    if j != 0 {
      writeSvgF(" ")
    }
    writeSvgF("%d,%d", p.X + xoffset, p.Y + yoffset)
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
