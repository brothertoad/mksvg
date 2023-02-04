package main

import (
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
  Flags: []cli.Flag {
    &cli.BoolFlag{Name: "points", Usage:"show points on curves"},
  },
  Action: doWeb,
}

const radius = 2

func doWeb(c *cli.Context) error {
  openSvg(path.Join(config.OutputDir, "mask.svg"))
  for _, obj := range(mask.InputObjects) {
    for _, curve := range(obj.Curves) {
      s := substitutePoints(curve)
      // At this point we have a list of space-separated points.
      points := stringToPoints(s)
      beziers := bezier.SvgControlPointsI(points)
      for _, b := range(beziers) {
        writeSvg(`<path d="` + b + `"/>` + "\n")
      }
      if c.Bool("points") {
        for _, p := range(points) {
          writeSvgF(`<circle cx="%d" cy="%d" r="%d"/>%s`, p.X, p.Y, radius, "\n")
        }
      }
    }
    for _, bezier := range(obj.Beziers) {
      points := stringToPoints(substitutePoints(bezier))
      if len(points) != 4 {
        btu.Fatal("Wrong number of points (%d) in bezier '%s'\n", len(points), bezier)
      }
      writeSvgF(`<path d="M %d,%d C`, points[0].X, points[0].Y)
      for _, p := range(points[1:]) {
        writeSvgF(" %d,%d", p.X, p.Y)
      }
      writeSvg(`"/>` + "\n")
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
