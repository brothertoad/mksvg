package main

import (
  "image"
  "github.com/brothertoad/btu"
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

func pointSetFromPath(tokens []string) []image.Point {
  points := make([]image.Point, 0)
  for j := 0; j < len(tokens); {
    cmd := tokens[j]
    j++
    switch cmd {
    case "M", "L":
      ensureEnoughPoints(cmd, 2, j, len(tokens))
    case "m", "l":
      ensureEnoughPoints(cmd, 2, j, len(tokens))
    case "V", "H":
      ensureEnoughPoints(cmd, 1, j, len(tokens))
    case "v", "h":
      ensureEnoughPoints(cmd, 1, j, len(tokens))
    case "C":
      ensureEnoughPoints(cmd, 3, j, len(tokens))
    case "c":
      ensureEnoughPoints(cmd, 3, j, len(tokens))
    case "Q":
      ensureEnoughPoints(cmd, 2, j, len(tokens))
    case "q":
      ensureEnoughPoints(cmd, 2, j, len(tokens))
    case "Z", "z":
    default:
      btu.Fatal("Unknown command in path: %s\n", cmd)
    }
  }
  return points
}

func ensureEnoughPoints(cmd string, req, offset, total int) {
  if (offset + req) > total {
    btu.Fatal("Not enough points for %s command, need %d, have %d\n", cmd, req, total - offset)
  }
}
