package main

import (
  _ "io/ioutil"
  _ "path"
  "strconv"
  "strings"
  "github.com/urfave/cli/v2"
  _ "github.com/brothertoad/btu"
)

var initCommand = cli.Command {
  Name: "init",
  Usage: "create initial output files",
  Action: doInit,
}

func doInit(c *cli.Context) error {
  /*
  filterTemplates("mask.html", "mask.css")
  dest := path.Join(config.OutputDir, "mask.jpg")
  btu.CopyFile(mask.Global.Image, dest)
  // Create an empty SVG file.
  openSvg(path.Join(config.OutputDir, "mask.svg"))
  closeSvg()
  */
  return nil
}

/*
func filterTemplates(templates ...string) {
  for _, template := range(templates) {
    inputPath := path.Join(config.TemplateDir, template)
    outputPath := path.Join(config.OutputDir, template)
    b, err := ioutil.ReadFile(inputPath)
    btu.CheckError(err)
    ioutil.WriteFile(outputPath, filterBytes(b), 0644)
  }
}

func filterBytes(src []byte) []byte {
  return []byte(filterString(string(src)))
}
*/

func filterString(src string) string {
  // Have to do STROKE-WIDTH first, so the part after the hypen is not
  // changed to the image width.
  s := strings.ReplaceAll(src, "STROKE-COLOR", mask.Global.StrokeColor)
  s = strings.ReplaceAll(s, "STROKE-WIDTH", strconv.Itoa(mask.Global.StrokeWidth))
  s = strings.ReplaceAll(s, "WIDTH", strconv.Itoa(mask.Global.Width))
  s = strings.ReplaceAll(s, "HEIGHT", strconv.Itoa(mask.Global.Height))
  return strings.ReplaceAll(s, "TITLE", mask.Global.Title)
}
