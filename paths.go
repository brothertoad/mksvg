package main

import (
  "fmt"
  "image"
  "strings"
  "unicode"
  "github.com/brothertoad/btu"
)

// Create a slice of components from the path tokens.  While doing so,
// convert all relative commands to absolute, and convert V/H/v/h to L
// so that each component has points (rather than a single parameter).
func parseComponentsFromPath(tokens []string) []pathComponent {
  x := 0
  y := 0
  haveZ := false
  components := make([]pathComponent, 0, len(tokens))
  var pc pathComponent
  for j := 0; j < len(tokens); {
    if haveZ {
      btu.Fatal("No commands can appear after Z/z.\n")
    }
    cmd := tokens[j]
    j++
    switch cmd {
    case "M", "L", "m", "l", "T", "t":
      x, y, pc = cmdToComponent(x, y, cmd, 2, tokens[j:])
      components = append(components, pc)
      j += 2
    case "V", "H", "v", "h":
      x, y, pc = createSingleParameterComponent(x, y, cmd, tokens[j:])
      components = append(components, pc)
      j += 1
    case "C", "c":
      x, y, pc = cmdToComponent(x, y, cmd, 6, tokens[j:])
      components = append(components, pc)
      j += 6
    case "Q", "q", "S", "s":
      x, y, pc = cmdToComponent(x, y, cmd, 4, tokens[j:])
      components = append(components, pc)
      j += 4
    case "A", "a":
      btu.Fatal("arcs in paths are not supported.\n")
    case "Z", "z":
      components = append(components, pathComponent{"Z", make([]image.Point, 0)})
      haveZ = true
    default:
      btu.Fatal("Unknown command in path: %s\n", cmd)
    }
  }
  return components
}

func cmdToComponent(x, y int, cmd string, numValues int, tokens []string) (int, int, pathComponent) {
    // ensure we have enough values
  if len(tokens) < numValues {
    btu.Fatal("Not enough values for %s command, need %d, have %d\n", cmd, numValues, len(tokens))
  }
  r := []rune(cmd)[0]
  relative := unicode.IsLower(r)
  numPoints := numValues / 2  // since each value is a coordinate, there are two per point
  p := make([]image.Point, numPoints)
  for j := 0; j < numPoints; j++ {
    p[j].X = parsePathNumber(tokens[2*j])
    p[j].Y = parsePathNumber(tokens[2*j + 1])
    // if this command is relative, make the coordinates absolute
    if relative {
      p[j].X += x
      p[j].Y += y
    }
  }
  x = p[numPoints-1].X
  y = p[numPoints-1].Y
  var pc pathComponent
  pc.cmd = strings.ToUpper(cmd)
  pc.points = p
  return x, y, pc
}

func createSingleParameterComponent(x, y int, cmd string, tokens []string) (int, int, pathComponent) {
  if len(tokens) < 1 {
    btu.Fatal("Need a value for %s command.\n", cmd)
  }
  p := make([]image.Point, 1)
  target := parsePathNumber(tokens[0])
  if cmd == "V" {
    p[0].X = x
    p[0].Y = target
  } else if cmd == "H" {
    p[0].X = target
    p[0].Y = y
  } else if cmd == "v" {
    p[0].X = x
    p[0].Y = target + y
  } else if cmd == "h" {
    p[0].X = target + x
    p[0].Y = y
  }
  var pc pathComponent
  pc.cmd = "L"
  pc.points = p
  return p[0].X, p[0].Y, pc
}

func isRelative(cmd string) bool {
  r := []rune(cmd)[0]
  return unicode.IsLower(r)
}

// This finds the center of the bounding rectangle of the components.
func centerAndBboxFromComponents(components []pathComponent) (image.Point, image.Rectangle) {
  xmin := 2000000000
  xmax := -2000000000
  ymin := 2000000000
  ymax := -2000000000
  for _, component := range components {
    for _, p := range component.points {
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

func dAndPointsFromComponents(components []pathComponent, center image.Point) (string, []image.Point) {
  // Use a map to avoid duplicate points.
  pointMap := make(map[image.Point]bool)
  var sb strings.Builder
  for _, component := range components {
    fmt.Fprintf(&sb, "%s", component.cmd)
    pss := make([]string, 0, len(component.points))
    for _, p := range component.points {
      pss = append(pss, fmt.Sprintf(" %d %d", p.X - center.X, p.Y - center.Y))
      pointMap[p] = true
    }
    fmt.Fprintf(&sb, "%s ", strings.Join(pss, ","))
  }
  points := make([]image.Point, 0, len(pointMap))
  for p, _ := range pointMap {
    points = append(points, p)
  }
  return sb.String(), points
}

func parsePathNumber(s string) int {
  // parse until the end of the string or we find a non-digit
  n := 0
  negate := false
  for j, ch := range s {
    // If there is a leading minus sign, the number is negative.
    if j == 0 && ch == '-' {
      negate = true
      continue
    }
    if !unicode.IsDigit(ch) {
      // anything other than a comma is a fatal error
      if ch != ',' {
        btu.Fatal("Found a non-digit that is not a comma in %s\n", s)
      }
      break
    }
    n = (n * 10) + int(ch - '0')
  }
  if negate {
    return -n
  }
  return n
}
