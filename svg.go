package main

import (
  "fmt"
  "image"
  "log"
  "math"
  "os"
  "github.com/brothertoad/btu"
)

var svgPrefix =`<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="%s" height="%s" viewBox="0 0 %d %d" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink= "http://www.w3.org/1999/xlink">
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
  pw := mask.Global.PhysicalWidth
  ph := mask.Global.PhysicalHeight
  w := mask.Global.Width
  h := mask.Global.Height
  if mask.Global.Scale != 0.0 {
    pw = scalePhysicalDimension(pw, mask.Global.Scale)
    ph = scalePhysicalDimension(ph, mask.Global.Scale)
    w = int(math.Round(float64(w) * mask.Global.Scale))
    h = int(math.Round(float64(h) * mask.Global.Scale))
  }
  writeSvgF(svgPrefix, pw, ph, w, h, mask.Global.StrokeColor, mask.Global.StrokeWidth, mask.Global.StrokeColor)
  if mask.Global.Scale != 0.0 {
    writeSvgF(`<g transform="scale(%.3f)">`, mask.Global.Scale)
    writeSvg("")
  }
}

// We assume the physical dimension is all ASCII (i.e., no Unicode).
func scalePhysicalDimension(dimension string, scale float64) string {
  // Get the numeric part of the dimension.
  var n int // index of units
  for j, c := range(dimension) {
     if c >= '0' && c <= '9' {
       continue
     }
     n = j
     break
  }
  i := btu.Atoi(dimension[0:n])
  units := dimension[n:]
  i = int(math.Round(float64(i) * scale))
  return fmt.Sprintf("%d%s", i, units)
}

func writePathToSvg(d string, xform string) {
  writeSvgF(`<path vector-effect="non-scaling-stroke" d="%s" %s/>`, d, xform)
  writeSvg("")  // to get a newline
}

func writePlainRectangleToSvg(x, y, width, height int) {
  writeSvgF(`<rect vector-effect="non-scaling-stroke" x="%d" y="%d" width="%d" height="%d"/>`, x, y, width, height)
  writeSvg("")
}

func writePointsToSvg(points []image.Point, center image.Point, radius int, xform string) {
  if !config.printPoints {
    return
  }
  for _, p := range(points) {
    writeSvgF(`<circle vector-effect="non-scaling-stroke" class="dot" cx="%d" cy="%d" r="%d" %s/>`, p.X - center.X, p.Y - center.Y, radius, xform)
    writeSvg("")
  }
}

func writeGridToSvg(x, y, width, height, spacing int) {
  rows := height / spacing;
  cols := width / spacing;
  for r := 0; r < rows; r++ {
    writeSvgF(`<line vector-effect="non-scaling-stroke" x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y + r * spacing, x + width, y + r * spacing)
    writeSvg("")
  }
  for c := 0; c < cols; c++ {
    writeSvgF(`<line vector-effect="non-scaling-stroke" x1="%d" y1="%d" x2="%d" y2="%d"/>`, x + c * spacing, y, x + c * spacing, y + height)
    writeSvg("")
  }
  writeSvg("")
}

func createTransformString(render RenderObject, object Object) string {
  return fmt.Sprintf(`transform="translate(%d,%d)%s"`, render.Translate.X, render.Translate.Y, createScaleString(render, object))
}

func createScaleString(render RenderObject, object Object) string {
  if render.Scale == 0.0 && object.Scale == 0.0 && render.Flip == "" {
    return ""
  }
  scale := render.Scale
  if scale == 0.0 {
    scale = 1.0
  }
  if object.Scale != 0.0 {
    scale *= object.Scale
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
