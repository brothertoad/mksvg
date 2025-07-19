package main

import (
  "fmt"
  "image"
  _ "image/jpeg"
  "io/ioutil"
  "strings"
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
    if !parsePaths(name, &obj) {
      obj.rawCurves = parsePointLists(obj.Curves)
      obj.rawBeziers = parsePointLists(obj.Beziers)
      obj.rawQBeziers = parsePointLists(obj.QBeziers)
      obj.rawLines = parsePointLists(obj.Lines)
      obj.rawSegments = parsePointLists(obj.Segments)
      obj.center, obj.bbox = getObjectCenter(obj)
      verifyPointSliceLengths("bezier", 4, obj.rawBeziers)
      verifyPointSliceLengths("qbezier", 3, obj.rawQBeziers)
      obj.d = createD(name, obj)
      obj.points = createPointSet(name, obj)
    }
    // OK, work around the fact that obj is a *copy* of the entry in
    // mask.Objects by copying the result back.
    mask.Objects[name] = obj
  }
}

func parsePaths(name string, obj *Object) bool {
  // If the object has both a path and paths, then stop.  This should never happen.
  if obj.Path != "" && len(obj.Paths) > 0 {
    btu.Fatal("Object %s has both a path and paths specified.\n", name)
  }
  if obj.Path == "" && len(obj.Paths) == 0 {
    return false
  }
  if len(obj.Paths) > 0 {
    obj.Path = strings.Join(obj.Paths, " ")
  }
  rawTokens := strings.Split(strings.ReplaceAll(obj.Path, "\n", " "), " ")
  // note that some of the rawTokens will be empty strings - we will remove those
  tokens := make([]string, 0, len(rawTokens))
  for _, rt := range rawTokens {
    if rt != "" {
      tokens = append(tokens, rt)
    }
  }
  obj.d = strings.Join(tokens, " ")
  obj.points = pointSetFromPath(tokens)
  obj.center, obj.bbox = getPathCenter(obj)
  return true
}

// This finds the center of the bounding rectangle of the object.
// This should probably be a method.
func getPathCenter(obj *Object) (image.Point, image.Rectangle) {
  xmin := 2000000000
  xmax := -2000000000
  ymin := 2000000000
  ymax := -2000000000
  for _, p := range obj.points {
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
  var c image.Point
  c.X = (xmin + xmax) / 2
  c.Y = (ymin + ymax) / 2
  return c, image.Rect(xmin - c.X, ymin - c.Y, xmax - c.X, ymax - c.Y)
}

// Parse a slice of strings, where each string is a list of points
// (either by coordinates or by name).
func parsePointLists(lists []string) []pointCollection {
  collections := make([]pointCollection, 0)
  if lists == nil || len(lists) == 0 {
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
      // extra whitespace) or a point name.  If it contains a comma, it is a pair
      // of coordinates.
      if strings.Contains(word, ",") {
        pc.points[j] = parseCoordinates(word)
      } else {
        if p, exists := mask.Points[word]; exists {
          pc.points[j] = p
        } else {
          btu.Fatal("No point named %s\n", word)
        }
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

// Verifies that each entry in the slice has the specified length.
// An invalid length is a fatal error.
func verifyPointSliceLengths(groupName string, length int, pcs []pointCollection) {
  for _, pc := range pcs {
    if len(pc.points) != length {
      btu.Fatal("Found a %s with an invalid number of points (%d)\n", groupName, len(pc.points))
    }
  }
}

// This finds the center of the bounding rectangle of the object.
// This should probably be a method.
func getObjectCenter(obj Object) (image.Point, image.Rectangle) {
  xmin := 2000000000
  xmax := -2000000000
  ymin := 2000000000
  ymax := -2000000000
  xmin, ymin, xmax, ymax = updateLimits(obj.rawCurves, xmin, ymin, xmax, ymax)
  xmin, ymin, xmax, ymax = updateLimits(obj.rawBeziers, xmin, ymin, xmax, ymax)
  xmin, ymin, xmax, ymax = updateLimits(obj.rawQBeziers, xmin, ymin, xmax, ymax)
  xmin, ymin, xmax, ymax = updateLimits(obj.rawLines, xmin, ymin, xmax, ymax)
  xmin, ymin, xmax, ymax = updateLimits(obj.rawSegments, xmin, ymin, xmax, ymax)
  var c image.Point
  c.X = (xmin + xmax) / 2
  c.Y = (ymin + ymax) / 2
  return c, image.Rect(xmin - c.X, ymin - c.Y, xmax - c.X, ymax - c.Y)
}

func updateLimits(pcs []pointCollection, xmin int, ymin int, xmax int, ymax int) (int, int, int, int) {
  for _, pc := range pcs {
    for _, p := range pc.points {
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
  return xmin, ymin, xmax, ymax
}
