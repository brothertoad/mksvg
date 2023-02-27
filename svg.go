package main

import (
  "fmt"
  "image"
  "log"
  "os"
  "github.com/brothertoad/bezier"
  "github.com/brothertoad/btu"
)

var svgPrefix =`<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="%dmm" height="%dmm" viewBox="0 0 %d %d" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink= "http://www.w3.org/1999/xlink">
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
  writeSvgF(svgPrefix, mask.Global.Width, mask.Global.Height, 10 * mask.Global.Width, 10 * mask.Global.Height,
    mask.Global.StrokeColor, mask.Global.StrokeWidth)
}

func writeCurveToSvg(curve pointCollection, center image.Point, render RenderObject) {
  beziers := bezier.GetControlPointsI(curve.points)
  for _, bezier := range(beziers) {
    writeSvgF(`<path d="M %d %d `, bezier.P0.X - center.X, bezier.P0.Y - center.Y)
    writeSvgF(` C %d %d,`, bezier.P1.X - center.X, bezier.P1.Y- center.Y)
    writeSvgF(` %d %d,`, bezier.P2.X - center.X, bezier.P2.Y - center.Y)
    writeSvgF(` %d %d"`, bezier.P3.X - center.X, bezier.P3.Y - center.Y)
    writeSvgF(` transform="translate(%d,%d) %s"/>`, render.Translate.X, render.Translate.Y, createScaleString(render))
    writeSvg("")  // to get a newline
  }
}

func writeBezierToSvg(bezier pointCollection, center image.Point, render RenderObject) {
  writeSvgF(`<path d="M %d %d `, bezier.points[0].X - center.X, bezier.points[0].Y - center.Y)
  writeSvgF(` C %d %d,`, bezier.points[1].X - center.X, bezier.points[1].Y - center.Y)
  writeSvgF(` %d %d,`, bezier.points[2].X - center.X, bezier.points[2].Y - center.Y)
  writeSvgF(` %d %d"`, bezier.points[3].X - center.X, bezier.points[3].Y - center.Y)
  writeSvgF(` transform="translate(%d,%d) %s"/>`, render.Translate.X, render.Translate.Y, createScaleString(render))
  writeSvg("")  // to get a newline
}

func writeLineToSvg(line pointCollection, center image.Point, render RenderObject) {
  writeSvgF(`<polyline points="`)
  for j, p := range(line.points) {
    if j != 0 {
      writeSvgF(" ")
    }
    writeSvgF("%d,%d", p.X - center.X, p.Y - center.Y)
  }
  writeSvgF(`" transform="translate(%d,%d) %s"/>`, render.Translate.X, render.Translate.Y, createScaleString(render))
  writeSvg("")
}

func writeRectangleToSvg(rect image.Rectangle, center image.Point, render RenderObject) {
  writeSvgF(`<rect x="%d" y="%d"`, rect.Min.X - center.X, rect.Min.Y - center.Y)
  writeSvgF(`width="%d" height="%d" `, rect.Max.X - rect.Min.X - center.X, rect.Max.Y - rect.Min.Y - center.Y)
  writeSvgF(`transform="translate(%d,%d) %s"/>`, render.Translate.X, render.Translate.Y, createScaleString(render))
  writeSvg("")  // to get a newline
}

func writePlainRectangleToSvg(x, y, width, height int) {
  writeSvgF(`<rect x="%d" y="%d" width="%d" height="%d"/>`, x, y, width, height)
  writeSvg("")
}

func createScaleString(render RenderObject) string {
  if render.Scale == 0.0  && render.Flip == "" {
    return ""
  }
  scale := render.Scale
  if scale == 0.0 {
    scale = 1.0
  }
  switch render.Flip {
  case "":
    return fmt.Sprintf("scale(%.3f)", scale)
  case "hflip":
    return fmt.Sprintf("scale(%.3f,%.3f)", - scale, scale)
  case "vflip":
    return fmt.Sprintf("scale(%.3f,%.3f)", scale, - scale)
  }
  log.Fatalf("Invalid flip value: %s\n", render.Flip)
  return ""
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
