package main

import (
  "image"
  "path"
  "strconv"
  "strings"
  "unicode"
  "unicode/utf8"
  "github.com/brothertoad/btu"
)

// const margin = 5

func render() {
  openSvg(path.Join(config.OutputDir, "mask.svg"))
  for _, render := range(mask.Renders) {
    if render.Hide {
      continue
    }
    if render.Comment != "" {
      writeSvgF(`<!-- %s -->`, render.Comment)
      writeSvg("")
    }
    obj := mask.Objects[render.Object]
    for _, curve := range(obj.rawCurves) {
      writeCurveToSvg(curve, obj.center, render)
    }
    for _, bezier := range(obj.rawBeziers) {
      writeBezierToSvg(bezier, obj.center, render)
    }
    for _, line := range(obj.rawLines) {
      writeLineToSvg(line, obj.center, render)
    }
    for _, rect := range(obj.rawRects) {
      writeRectangleToSvg(rect, obj.center, render)
    }
    writeSvg("")
  }
  if config.printBorder {
    w := mask.Global.Width * 10 - 2 * config.MarginEdge
    h := mask.Global.Height * 10 - 2 * config.MarginEdge
    writePlainRectangleToSvg(config.MarginEdge, config.MarginEdge, w, h)
  }
  closeSvg()
}

func substitutePoints(s string) string {
  var sb strings.Builder
  words := strings.Fields(s)
  for j, word := range(words) {
    if j != 0 {
      sb.WriteRune(' ')
    }
    rune, _ := utf8.DecodeRuneInString(word)
    if unicode.IsDigit(rune) {
      sb.WriteString(word)
    } else {
      p := mask.Points[word]
      sb.WriteString(ptos(p, ","))
    }
  }
  return sb.String()
}

// Convert a list of space-separated points to a slice of image.Points.
func stringToPoints(s string) []image.Point {
  words := strings.Fields(s)
  points := make([]image.Point, len(words), len(words))
  for j, word := range(words) {
    coords := strings.Split(word, ",")
    var err error
    points[j].X, err = strconv.Atoi(coords[0])
    btu.CheckError(err)
    points[j].Y, err = strconv.Atoi(coords[1])
    btu.CheckError(err)
  }
  return points
}

// Show an image.Point as a string.
func ptos(p image.Point, separator string) string {
  return strconv.Itoa(p.X) + separator + strconv.Itoa(p.Y)
}
