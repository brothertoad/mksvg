package main

import (
  "io/ioutil"
  "path"
  "strconv"
  "strings"
  "github.com/urfave/cli/v2"
)

var initCommand = cli.Command {
  Name: "init",
  Usage: "create initial output files",
  Action: doInit,
}

func doInit(c *cli.Context) error {
  filterTemplates("mask.html", "mask.css")
  dest := path.Join(config.OutputDir, "mask.jpg")
  copyFile(mask.Global.Image, dest)
  // Create an empty SVG file.
  openSvg(path.Join(config.OutputDir, "mask.svg"))
  closeSvg()
  return nil
}

func filterTemplates(templates ...string) {
  for _, template := range(templates) {
    inputPath := path.Join(config.TemplateDir, template)
    outputPath := path.Join(config.OutputDir, template)
    b, err := ioutil.ReadFile(inputPath)
    checkError(err)
    ioutil.WriteFile(outputPath, filterBytes(b), 0644)
  }
}

func filterBytes(src []byte) []byte {
  return []byte(filterString(string(src)))
}

func filterString(src string) string {
  // Have to do STROKE-WIDTH first, so the part after the hypen is not
  // changed to the image width.
  s := strings.ReplaceAll(src, "STROKE-COLOR", mask.Global.StrokeColor)
  s = strings.ReplaceAll(s, "STROKE-WIDTH", strconv.Itoa(mask.Global.StrokeWidth))
  s = strings.ReplaceAll(s, "WIDTH", strconv.Itoa(mask.width))
  s = strings.ReplaceAll(s, "HEIGHT", strconv.Itoa(mask.height))
  return strings.ReplaceAll(s, "TITLE", mask.Global.Title)
}
