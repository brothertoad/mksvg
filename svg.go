package main

import (
  "fmt"
  "image"
  "log"
  "math"
  "os"
  "github.com/brothertoad/bezier"
  "github.com/brothertoad/btu"
)

var svgPrefix =`<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="%dmm" height="%dmm" viewBox="0 0 %d %d" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink= "http://www.w3.org/1999/xlink">
  <style>
  :root {
    fill: none;
    stroke: %s;
    stroke-width: %d;
  }
  .dot {
    fill: %s;
  }
  </style>
`
var svgSuffix = `</svg>
`

var svgFile *os.File

func openSvg(path string) {
  svgFile = btu.CreateFile(path)
  w := mask.Global.Width
  h := mask.Global.Height
  if mask.Global.Scale != 0.0 {
    w = int(math.Round(float64(w) * mask.Global.Scale))
    h = int(math.Round(float64(h) * mask.Global.Scale))
  }
  writeSvgF(svgPrefix, w, h, 10 * w, 10 * h, mask.Global.StrokeColor, mask.Global.StrokeWidth, mask.Global.StrokeColor)
  if mask.Global.Scale != 0.0 {
    writeSvgF(`<g transform="scale(%.3f)">`, mask.Global.Scale)
    writeSvg("")
  }
}

func writeCurveToSvg(curve pointCollection, center image.Point, render RenderObject) {
  beziers := bezier.GetControlPointsI(curve.points)
  xform := createTransformString(render)
  for _, bezier := range(beziers) {
    writeSvgF(`<path vector-effect="non-scaling-stroke" d="M %d %d `, bezier.P0.X - center.X, bezier.P0.Y - center.Y)
    writeSvgF(` C %d %d,`, bezier.P1.X - center.X, bezier.P1.Y- center.Y)
    writeSvgF(` %d %d,`, bezier.P2.X - center.X, bezier.P2.Y - center.Y)
    writeSvgF(` %d %d"`, bezier.P3.X - center.X, bezier.P3.Y - center.Y)
    writeSvgF(` %s/>`, xform)
    writeSvg("")  // to get a newline
  }
  writePointsToSvg(curve.points, center, xform)
}

func writeBezierToSvg(bezier pointCollection, center image.Point, render RenderObject) {
  xform := createTransformString(render)
  writeSvgF(`<path vector-effect="non-scaling-stroke" d="M %d %d `, bezier.points[0].X - center.X, bezier.points[0].Y - center.Y)
  writeSvgF(` C %d %d,`, bezier.points[1].X - center.X, bezier.points[1].Y - center.Y)
  writeSvgF(` %d %d,`, bezier.points[2].X - center.X, bezier.points[2].Y - center.Y)
  writeSvgF(` %d %d"`, bezier.points[3].X - center.X, bezier.points[3].Y - center.Y)
  writeSvgF(` %s/>`, xform)
  writeSvg("")  // to get a newline
  writePointsToSvg(bezier.points, center, xform)
}

func writeLineToSvg(line pointCollection, center image.Point, render RenderObject) {
  xform := createTransformString(render)
  writeSvgF(`<polyline vector-effect="non-scaling-stroke" points="`)
  for j, p := range(line.points) {
    if j != 0 {
      writeSvgF(" ")
    }
    writeSvgF("%d,%d", p.X - center.X, p.Y - center.Y)
  }
  writeSvgF(`" %s/>`, xform)
  writeSvg("")
  writePointsToSvg(line.points, center, xform)
}

func writeRectangleToSvg(rect image.Rectangle, center image.Point, render RenderObject) {
  writeSvgF(`<rect vector-effect="non-scaling-stroke" x="%d" y="%d"`, rect.Min.X - center.X, rect.Min.Y - center.Y)
  writeSvgF(`width="%d" height="%d" `, rect.Max.X - rect.Min.X - center.X, rect.Max.Y - rect.Min.Y - center.Y)
  writeSvgF(`%s/>`, createTransformString(render))
  writeSvg("")  // to get a newline
}

func writePlainRectangleToSvg(x, y, width, height int) {
  writeSvgF(`<rect vector-effect="non-scaling-stroke" x="%d" y="%d" width="%d" height="%d"/>`, x, y, width, height)
  writeSvg("")
}

func writePointsToSvg(points []image.Point, center image.Point, xform string) {
  if !config.printPoints {
    return
  }
  for _, p := range(points) {
    writeSvgF(`<circle vector-effect="non-scaling-stroke" class="dot" cx="%d" cy="%d" r="%d" %s/>`, p.X - center.X, p.Y - center.Y, config.PointRadius, xform)
    writeSvg("")
  }
}

func createTransformString(render RenderObject) string {
  return fmt.Sprintf(`transform="translate(%d,%d)%s"`, render.Translate.X, render.Translate.Y, createScaleString(render))
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
    return fmt.Sprintf(" scale(%.3f)", scale)
  case "hflip":
    return fmt.Sprintf(" scale(%.3f,%.3f)", - scale, scale)
  case "vflip":
    return fmt.Sprintf(" scale(%.3f,%.3f)", scale, - scale)
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
  if mask.Global.Scale != 0.0 {
    writeSvg("</g>")
  }
  svgFile.WriteString(svgSuffix)
  err := svgFile.Close()
  btu.CheckError(err)
}
