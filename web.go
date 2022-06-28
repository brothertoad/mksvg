package main

import (
  "path"
  "strconv"
  "strings"
  "unicode"
  "unicode/utf8"
  "github.com/urfave/cli/v2"
)

var webCommand = cli.Command {
  Name: "web",
  Usage: "create SVG file for the web",
  Action: doWeb,
}

func doWeb(c *cli.Context) error {
  openSvg(path.Join(config.OutputDir, "mask.svg"))
  // Need to output the actual SVG here.
  // writeSvg(`<polyline points="419,329 894,632 835,610 654,592 361,528 246,413 111,443 114,549"/>`)
  for _, obj := range(mask.Objects) {
    for _, line := range(obj.Lines) {
      // For each line, just find the point names and change them to the actual coordinates.
      writeSvg(`<polyline points="` + substitutePoints(line) + `"/>` + "\n")
    }
  }
  closeSvg()
  return nil
}

func substitutePoints(s string) string {
  var sb strings.Builder
  words := strings.Fields(s)
  for j, word := range(words) {
    if j != 0 {
      sb.WriteRune(' ')
    }
    rune, _ := utf8.DecodeRuneInString(word)
    if unicode.IsDigit(rune) {
      sb.WriteString(word)
    } else {
      p := mask.Points[word]
      sb.WriteString(strconv.Itoa(p.X) + "," + strconv.Itoa(p.Y))
    }
  }
  return sb.String()
}
