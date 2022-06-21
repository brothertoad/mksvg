package main

import (
  "path"
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
  writeSvg(`<polyline points="419,329 894,632 835,610 654,592 361,528 246,413 111,443 114,549"/>`)
  closeSvg()
  return nil
}
