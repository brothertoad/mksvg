package main

import (
  "fmt"
  "image"
  _ "image/jpeg"
  "io/ioutil"
  "strings"
  "unicode"
  "unicode/utf8"
  "github.com/pelletier/go-toml"
  "github.com/brothertoad/btu"
)

func parseMask(paths []string) {
  for _, path := range(paths) {
    b, err := ioutil.ReadFile(path)
    btu.CheckError2(err, "Unable to read TOML file '%s'", path)
    err = toml.Unmarshal(b, &mask)
    btu.CheckError2(err, "Unable to unmarshal TOML file '%s'", path)
  }

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
    // If this object has rectangles, stop here - the close path logic does not
    // support rectangles.
    if obj.Rects != nil && len(obj.Rects) > 0 {
      btu.Fatal("%s has rectangles, which are no longer supported.\n", name)
    }
    obj.rawCurves = parsePointLists(obj.Curves)
    obj.rawBeziers = parsePointLists(obj.Beziers)
    obj.rawQBeziers = parsePointLists(obj.QBeziers)
    obj.rawLines = parsePointLists(obj.Lines)
    obj.center, obj.bbox = getObjectCenter(obj)
    obj.d = createD(name, obj)
    obj.points = createPointSet(name, obj)
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
    btu.Fatal("Expected two coordinates in point '%s'\n", s)
  }
  point.X = btu.Atoi2(coords[0], "Can't convert coordinate value '%s' to a number (full coordinates '%s')", coords[0], s)
  point.Y = btu.Atoi2(coords[1], "Can't convert coordinate value '%s' to a number (full coordinates '%s')", coords[1], s)
  return point
}

// This finds the center of the bounding rectangle of the object.
// This should probably be a method.
func getObjectCenter(obj Object) (image.Point, image.Rectangle) {
  xmin := 2000000000
  xmax := -2000000000
  ymin := 2000000000
  ymax := -2000000000
  for _, curve := range(obj.rawCurves) {
    for _, p := range(curve.points) {
      if p.X < xmin {
        xmin = p.X
      }
      if p.Y < ymin {
        ymin = p.Y
      }
      if p.X > xmax {
        xmax = p.X
      }
      if p.Y > ymax {
        ymax = p.Y
      }
    }
  }
  for _, bezier := range(obj.rawBeziers) {
    for _, p := range(bezier.points) {
      if p.X < xmin {
        xmin = p.X
      }
      if p.Y < ymin {
        ymin = p.Y
      }
      if p.X > xmax {
        xmax = p.X
      }
      if p.Y > ymax {
        ymax = p.Y
      }
    }
  }
  for _, qbezier := range(obj.rawQBeziers) {
    for _, p := range(qbezier.points) {
      if p.X < xmin {
        xmin = p.X
      }
      if p.Y < ymin {
        ymin = p.Y
      }
      if p.X > xmax {
        xmax = p.X
      }
      if p.Y > ymax {
        ymax = p.Y
      }
    }
  }
  for _, line := range(obj.rawLines) {
    for _, p := range(line.points) {
      if p.X < xmin {
        xmin = p.X
      }
      if p.Y < ymin {
        ymin = p.Y
      }
      if p.X > xmax {
        xmax = p.X
      }
      if p.Y > ymax {
        ymax = p.Y
      }
    }
  }
  var c image.Point
  c.X = (xmin + xmax) / 2
  c.Y = (ymin + ymax) / 2
  return c, image.Rect(xmin - c.X, ymin - c.Y, xmax - c.X, ymax - c.Y)
}
