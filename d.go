package main

import (
  "fmt"
  "image"
  "github.com/brothertoad/btu"
  "github.com/brothertoad/bezier"
)

type segmentInfo struct {
  start, end image.Point
  d string
}

func createD(name string, obj Object) string {
  segments := make([]segmentInfo, 0)
  for _, curve := range(obj.rawCurves) {
    segments = append(segments, createCurveSegment(curve, obj.center))
  }
  for _, bezier := range(obj.rawBeziers) {
    segments = append(segments, createBezierSegment(bezier, obj.center))
  }
  for _, qbezier := range(obj.rawQBeziers) {
    segments = append(segments, createQBezierSegment(qbezier, obj.center))
  }
  for _, line := range(obj.rawLines) {
    segments = append(segments, createLineSegment(line, obj.center))
  }
  if len(segments) > 0 {
    // Sort the segments so that the starting point of each segment is the
    // ending point of the previous segment.
    segments = sortSegments(name, segments)
  }
  // We will return a move to the start of the first segment, followed by
  // the d of each segment.
  d := fmt.Sprintf("M %d,%d", segments[0].start.X - obj.center.X, segments[0].start.Y - obj.center.Y)
  for _, segment := range(segments) {
    d = fmt.Sprintf("%s%s", d, segment.d)
  }
  return d
}

func sortSegments(name string, unsorted []segmentInfo) []segmentInfo {
  sorted := make([]segmentInfo, 0, len(unsorted))
  // We will start, arbitrarily, with the first unsorted segment.
  sorted = append(sorted, unsorted[0])
  // Create a pseudo-set of the remaining unsorted segments.
  // (See https://stackoverflow.com/questions/34018908/golang-why-dont-we-have-a-set-datastructure
  // for the paradigm.)
  remaining := make(map[int]bool, len(unsorted) - 1)
  for j := 1; j < len(unsorted); j++ {
    remaining[j] = true
  }
  nextStart := sorted[0].end
  for len(remaining) > 0 {
    // Find the entry in remaining with a start point that matches the previous end point.
    found := -1
    for k, _ := range(remaining) {
      if unsorted[k].start == nextStart {
        sorted = append(sorted, unsorted[k])
        nextStart = unsorted[k].end
        found = k
        break
      }
    }
    // If found is still less than zero, we didn't find one - this is a fatal error.
    if found < 0 {
      btu.Fatal("Can't find matching segment in object %s.", name)
    }
    // Remove the one we found.
    delete(remaining, found)
  }
  return sorted
}

func createCurveSegment(curve pointCollection, center image.Point) segmentInfo {
  var segment segmentInfo
  beziers := bezier.GetControlPointsI(curve.points)
  segment.start = beziers[0].P0
  segment.end = beziers[len(beziers)-1].P3
  d := ""
  for _, bezier := range(beziers) {
    d = fmt.Sprintf("%s C %d %d, %d %d, %d %d", d, bezier.P1.X - center.X, bezier.P1.Y- center.Y,
      bezier.P2.X - center.X, bezier.P2.Y - center.Y, bezier.P3.X - center.X, bezier.P3.Y - center.Y)
  }
  segment.d = d
  return segment
}

func createBezierSegment(bezier pointCollection, center image.Point) segmentInfo {
  var segment segmentInfo
  p := bezier.points
  segment.start = p[0]
  segment.end = p[3]
  segment.d = fmt.Sprintf(" C %d %d, %d %d, %d %d", p[1].X - center.X, p[1].Y - center.Y,
    p[2].X - center.X, p[2].Y - center.Y, p[3].X - center.X, p[3].Y - center.Y)
  return segment
}

func createQBezierSegment(qbezier pointCollection, center image.Point) segmentInfo {
  var segment segmentInfo
  p := qbezier.points
  segment.start = p[0]
  segment.end = p[2]
  segment.d = fmt.Sprintf(" Q %d %d, %d %d", p[1].X - center.X, p[1].Y - center.Y,
    p[2].X - center.X, p[2].Y - center.Y)
  return segment
}

func createLineSegment(line pointCollection, center image.Point) segmentInfo {
  var segment segmentInfo
  p := line.points
  segment.start = p[0]
  segment.end = p[len(p)-1]
  d := ""
  for j := 1; j < len(p); j++ {
    d = fmt.Sprintf("%s L %d,%d", d, p[j].X - center.X, p[j].Y - center.Y)
  }
  segment.d = d
  return segment
}
