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
  scale := finalScale()
  pw = scalePhysicalDimension("physical width", pw, scale)
  ph = scalePhysicalDimension("physical height", ph, scale)
  w = int(math.Round(float64(w) * scale))
  h = int(math.Round(float64(h) * scale))
  writeSvgF(svgPrefix, pw, ph, w, h, mask.Global.StrokeColor, mask.Global.StrokeWidth, mask.Global.StrokeColor)
  writeSvgF(`<g transform="scale(%.3f)">`, scale)
  writeSvg("")
}

// The final scale is the product of the scale in the global section of the mask file and the
// scale passed as a command line flag.
func finalScale() float64 {
  if mask.Global.Scale == 0.0 && config.scale == 0.0 {
    return 1.0
  }
  scale := mask.Global.Scale
  if scale == 0.0 {
    scale = 1.0
  }
  if config.scale != 0.0 {
    scale = scale * config.scale
  }
  return scale
}

// We assume the physical dimension is all ASCII (i.e., no Unicode).
func scalePhysicalDimension(name, dimension string, scale float64) string {
  // Get the numeric part of the dimension.
  var n int // index of units
  for j, c := range(dimension) {
     if c >= '0' && c <= '9' {
       continue
     }
     n = j
     break
  }
  i := btu.Atoi2(dimension[0:n], "Can't convert '%s' (numeric value '%s', full value '%s') to a number", name, dimension[0:n], dimension)
  units := dimension[n:]
  i = int(math.Round(float64(i) * scale))
  return fmt.Sprintf("%d%s", i, units)
}

func writePathToSvg(d string, xform string) {
  writeSvgF(`<path vector-effect="non-scaling-stroke" d="%s" %s/>`, d, xform)
  writeSvg("")  // to get a newline
}

func writeRectangleToSvg(r image.Rectangle, xform string) {
  writeSvgF(`<rect vector-effect="non-scaling-stroke" x="%d" y="%d" width="%d" height="%d" %s/>`, r.Min.X, r.Min.Y, r.Max.X - r.Min.X, r.Max.Y - r.Min.Y, xform)
  writeSvg("")
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
  writeSvg("</g>")
  svgFile.WriteString(svgSuffix)
  err := svgFile.Close()
  btu.CheckError2(err, "Unable to close SVG file.")
}
