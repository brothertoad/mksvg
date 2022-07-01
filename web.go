package main

import (
  "fmt"
  "image"
  "path"
  "strconv"
  "strings"
  "unicode"
  "unicode/utf8"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/bezier"
  "github.com/brothertoad/btu"
)

var webCommand = cli.Command {
  Name: "web",
  Usage: "create SVG file for the web",
  Action: doWeb,
}

const radius = 2

func doWeb(c *cli.Context) error {
  openSvg(path.Join(config.OutputDir, "mask.svg"))
  for _, obj := range(mask.Objects) {
    for _, curve := range(obj.Curves) {
      s := substitutePoints(curve)
      // At this point we have a list of space-separated points.
      points := stringToPoints(s)
      beziers := bezier.SvgControlPointsI(points)
      for _, b := range(beziers) {
        writeSvg(`<path d="` + b + `"/>` + "\n")
      }
      // Write dots to represent the knots.  Might want to make this
      // controllable by command line flag.
      for _, p := range(points) {
        s := fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d"/>%s`, p.X, p.Y, radius, "\n")
        writeSvg(s)
      }
    }
    for _, line := range(obj.Lines) {
      // For each line, just find the point names and change them to the actual coordinates.
      writeSvg(`<polyline points="` + substitutePoints(line) + `"/>` + "\n")
    }
  }
  closeSvg()
  return nil
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
