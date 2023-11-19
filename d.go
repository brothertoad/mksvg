package main

import (
  "fmt"
  "image"
)

type segmentInfo struct {
  start, end image.Point
  d string
}

func createD(obj Object) string {
  segments := make([]segmentInfo, 0)
  for _, line := range(obj.rawLines) {
    segments = append(segments, createLineSegment(line))
  }
  return ""
}

func createLineSegment(line pointCollection) segmentInfo {
  var segment segmentInfo
  p := line.points
  segment.start = p[0]
  segment.end = p[len(p)-1]
  d := "L "
  for j := 1; j < len(p); j++ {
    d = fmt.Sprintf("%s %d,%d", d, p[j].X, p[j].Y)
  }
  segment.d = d
  return segment
}
