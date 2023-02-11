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

type pointSlice []image.Point
type rectangle struct {
  width, height int
}

// The fields prefixed with "raw" are computed from the
// corresponding non-raw fields.
type Object struct {
  Curves []string
  Beziers []string
  Lines []string
  Rects []string
  rawCurves []pointSlice
  rawBeziers []pointSlice
  rawLines []pointSlice
  rawRects []rectangle
}

// Note that some fields in this object are read directly from the input file,
// whereas others are computed.
var mask struct {
  Global GlobalInfo
  Points map[string]image.Point
  Objects map[string]Object
}

func loadMask(path string) {
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  err = toml.Unmarshal(b, &mask)
  btu.CheckError(err)
  if mask.Global.PrintName == "" {
    mask.Global.PrintName = mask.Global.Title
  }
  // If there are not objects, there is nothing to do.
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
    fmt.Printf("    Curves: %+v\n", obj.Curves)
    fmt.Printf("    Beziers: %+v\n", obj.Beziers)
    fmt.Printf("    Lines: %+v\n", obj.Lines)
    printPointSlices("raw curves", obj.rawCurves)
    printPointSlices("raw beziers", obj.rawBeziers)
    printPointSlices("raw lines", obj.rawLines)
  }
}

func printPointSlices(label string, slices []pointSlice) {
  fmt.Printf("    %s:\n", label)
  for _, slice := range(slices) {
    fmt.Print("     ")
    for _, p := range(slice) {
      fmt.Printf(" %d,%d", p.X, p.Y)
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
    //  Each rectangle just has a width and a height, so we will parse
    // the rectangles as if they were points.  Then we check that each
    // rectangle consists of a single point.
    rectsAsPoints := parsePointLists(obj.Rects)
    obj.rawRects = make([]rectangle, len(rectsAsPoints))
    for j, r := range(rectsAsPoints) {
      if len(r) != 1 {
        btu.Fatal("Found a rectangle with more than two coordinates\n")
      }
      obj.rawRects[j].width = r[0].X
      obj.rawRects[j].height = r[0].Y
    }
    // OK, work around the fact that obj is a *copy* of the entry in
    // mask.Objects by copying the result back.
    mask.Objects[name] = obj
  }
}

// Parse a slice of strings, where each string is a list of points
// (either by coordinates or by name).
func parsePointLists(lists []string) []pointSlice {
  slices := make([]pointSlice, 0)
  if lists == nil {
    return slices
  }
  for _, list := range(lists) {
    words := strings.Split(list, " ")
    if len(words) == 0 {
      btu.Warn("Found empty point list\n")
      continue
    }
    slice := make([]image.Point, len(words))
    for j, word := range(words) {
      // Each word is either a pair of coordinates (separated by a comma, with no
      // extra whitespace) or a point name.  Look at the first rune to determine
      // which of the two it is.
      rune, _ := utf8.DecodeRuneInString(word)
      if unicode.IsDigit(rune) {
        slice[j] = parseCoordinates(word)
      } else {
        slice[j] = mask.Points[word]
      }
    }
    slices = append(slices, slice)
  }
  return slices
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
