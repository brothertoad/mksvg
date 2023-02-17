package main

import (
  "fmt"
  "image"
  _ "image/jpeg"
  "io/ioutil"
  "strconv"
  "strings"
  "unicode"
  "unicode/utf8"
  "github.com/pelletier/go-toml"
  "github.com/brothertoad/btu"
)

type GlobalInfo struct {
  Image string
  StrokeColor string
  StrokeWidth int
  Title string
  PrintName string
  Width int
  Height int
}

type pointCollection struct {
    points []image.Point
    // center image.Point
  }

type RenderObject struct {
  Object    string
  Translate image.Point
  Scale     float64
  Flip      string
}

// The fields prefixed with "raw" are computed from the
// corresponding non-raw fields.
type Object struct {
  Curves []string
  Beziers []string
  Lines []string
  Rects []string
  rawCurves []pointCollection
  rawBeziers []pointCollection
  rawLines []pointCollection
  rawRects []image.Rectangle
  center image.Point
}

// Note that some fields in this object are read directly from the input file,
// whereas others are computed.
var mask struct {
  Global GlobalInfo
  Points map[string]image.Point
  Objects map[string]Object
  Renders []RenderObject
}

func parseMask(path string) {
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  err = toml.Unmarshal(b, &mask)
  btu.CheckError(err)

  // Get values from config, if necessary.
  if mask.Global.StrokeColor == "" {
    mask.Global.StrokeColor = config.StrokeColor
  }
  if mask.Global.StrokeWidth == 0 {
    mask.Global.StrokeWidth = config.StrokeWidth
  }

  if mask.Global.PrintName == "" {
    mask.Global.PrintName = mask.Global.Title
  }
  // If there are no objects, there is nothing to do.
  if mask.Objects == nil {
    btu.Fatal("No objects in mask file.\n")
  }
  // If there are no points defined, make Points an empty slice to ease error checking.
  if mask.Points == nil {
    mask.Points = make(map[string]image.Point)
  }
  // If a list of items is nil, make it an empty slice to ease error checking.
  for _, obj := range(mask.Objects) {
    if obj.Rects == nil {
      obj.Rects = make([]string, 0)
    }
  }
  parseObjects()
  // dumpObjects()
}

func dumpObjects() {
  // Let's dump the objects.
  for name, obj := range(mask.Objects) {
    fmt.Printf("Dump of object %s\n", name)
    printPointSlices("curves", obj.rawCurves)
    printPointSlices("beziers", obj.rawBeziers)
    printPointSlices("lines", obj.rawLines)
  }
}

func printPointSlices(label string, collections []pointCollection) {
  fmt.Printf("    %s:\n", label)
  for _, pc := range(collections) {
    fmt.Print("     ")
    for _, p := range(pc.points) {
      fmt.Printf(" (%d,%d)", p.X, p.Y)
    }
    fmt.Printf("\n")
  }
}

// Parse the objects in the mask to actual numeric values.
func parseObjects() {
  for name, obj := range(mask.Objects) {
    obj.rawCurves = parsePointLists(obj.Curves)
    obj.rawBeziers = parsePointLists(obj.Beziers)
    obj.rawLines = parsePointLists(obj.Lines)
    //  Each rectangle consists of two points, so we will parse
    // the rectangles as if they were points.  Then we check that each
    // rectangle consists of exactly two points.
    rectsAsPoints := parsePointLists(obj.Rects)
    obj.rawRects = make([]image.Rectangle, len(rectsAsPoints))
    for j, r := range(rectsAsPoints) {
      if len(r.points) != 2 {
        btu.Fatal("Found a rectangle with more than two points\n")
      }
      obj.rawRects[j] = image.Rectangle{r.points[0], r.points[1]}
    }
    obj.center = getObjectCenter(obj)
    // OK, work around the fact that obj is a *copy* of the entry in
    // mask.Objects by copying the result back.
    mask.Objects[name] = obj
  }
}

// Parse a slice of strings, where each string is a list of points
// (either by coordinates or by name).
func parsePointLists(lists []string) []pointCollection {
  collections := make([]pointCollection, 0)
  if lists == nil {
    return collections
  }
  for _, list := range(lists) {
    words := strings.Split(list, " ")
    if len(words) == 0 {
      btu.Warn("Found empty point list\n")
      continue
    }
    var pc pointCollection
    pc.points = make([]image.Point, len(words))
    for j, word := range(words) {
      // Each word is either a pair of coordinates (separated by a comma, with no
      // extra whitespace) or a point name.  Look at the first rune to determine
      // which of the two it is.
      rune, _ := utf8.DecodeRuneInString(word)
      if unicode.IsDigit(rune) {
        pc.points[j] = parseCoordinates(word)
      } else {
        pc.points[j] = mask.Points[word]
      }
    }
    collections = append(collections, pc)
  }
  return collections
}

func parseCoordinates(s string) image.Point {
  var point image.Point
  coords := strings.Split(s, ",")
  if len(coords) != 2 {
    btu.Fatal("Expected two coordinates in point\n")
  }
  var err error
  point.X, err = strconv.Atoi(coords[0])
  btu.CheckError(err)
  point.Y, err = strconv.Atoi(coords[1])
  btu.CheckError(err)
  return point
}

// This should probably be a method.
func getObjectCenter(obj Object) image.Point{
  sumx := 0
  sumy := 0
  n := 0
  for _, curve := range(obj.rawCurves) {
    for _, p := range(curve.points) {
      sumx += p.X
      sumy += p.Y
      n++
    }
  }
  for _, bezier := range(obj.rawBeziers) {
    for _, p := range(bezier.points) {
      sumx += p.X
      sumy += p.Y
      n++
    }
  }
  for _, line := range(obj.rawLines) {
    for _, p := range(line.points) {
      sumx += p.X
      sumy += p.Y
      n++
    }
  }
  var c image.Point
  c.X = sumx / n
  c.Y = sumy / n
  return c
}
