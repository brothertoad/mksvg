package main

import (
  "fmt"
  "image"
  "github.com/brothertoad/bezier"
)

type segmentInfo struct {
  start, end image.Point
  d string
}

func createD(obj Object) string {
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
  // Sort the segments so that the starting point of each segment is the
  // ending point of the previous segment.
  return ""
}

func createCurveSegment(curve pointCollection, center image.Point) segmentInfo {
  var segment segmentInfo
  beziers := bezier.GetControlPointsI(curve.points)
  segment.start = beziers[0].P0
  segment.end = beziers[len(beziers)-1].P3
  d := ""
  for _, bezier := range(beziers) {
    d = fmt.Sprintf("%s C %d %d, %d %d, %d %d", bezier.P1.X - center.X, bezier.P1.Y- center.Y,
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
  segment.d = fmt.Sprintf("C %d %d, %d %d, %d %d", p[1].X - center.X, p[1].Y - center.Y,
    p[2].X - center.X, p[2].Y - center.Y, p[3].X - center.X, p[3].Y - center.Y)
  return segment
}

func createQBezierSegment(qbezier pointCollection, center image.Point) segmentInfo {
  var segment segmentInfo
  p := qbezier.points
  segment.start = p[0]
  segment.end = p[2]
  segment.d = fmt.Sprintf("C %d %d, %d %d, %d %d", p[1].X - center.X, p[1].Y - center.Y,
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
    d = fmt.Sprintf("%s L %d,%d", d, p[j].X - center.X, p[j].Y, - center.Y)
  }
  segment.d = d
  return segment
}
