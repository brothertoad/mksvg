package main

import (
  "image"
)

func createPointSet(name string, obj Object) []image.Point {
  pointMap := make(map[image.Point]bool)
  addPoints(pointMap, getAllPoints(obj.rawCurves))
  addPoints(pointMap, getAllPoints(obj.rawBeziers))
  addPoints(pointMap, getAllPoints(obj.rawQBeziers))
  addPoints(pointMap, getAllPoints(obj.rawLines))
  // Create a slice of the keys in pointMap.  Link to logic:
  // https://stackoverflow.com/questions/21362950/getting-a-slice-of-keys-from-a-map
  points := make([]image.Point, len(pointMap))
  j := 0
  for k := range(pointMap) {
    points[j] = k
    j++
  }
  return points
}

func addPoints(sum map[image.Point]bool, points[]image.Point) {
  for _, p := range(points) {
    sum[p] = true
  }
}

func getAllPoints(colls []pointCollection) []image.Point {
  allPoints := make([]image.Point, 0)
  for _, pc := range(colls) {
    for _, p := range(pc.points) {
      allPoints = append(allPoints, p)
    }
  }
  return allPoints
}
