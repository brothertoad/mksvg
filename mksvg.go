package main

import (
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "time"
  "github.com/urfave/cli/v2"
  "github.com/pelletier/go-toml"
  "github.com/brothertoad/btu"
)

func main() {
  app := &cli.App{
    Name: "mksvg",
    Compiled: time.Now(),
    Usage: "create an SVG file from a text file",
    Flags: []cli.Flag{
      &cli.StringFlag{Name: "config", Usage: "path to configuration file", Required: true, EnvVars: []string{"MKSVG_CONFIG"}},
      &cli.StringFlag{Name: "input", Usage: "input file", Aliases: []string{"i"}, DefaultText: "mask.toml", Value: "mask.toml", Destination: &config.inputPath},
      &cli.StringFlag{Name: "output", Usage: "output file", Aliases: []string{"o"}, DefaultText: "mask.svg", Value: "mask.svg"},
      &cli.BoolFlag{Name: "points", Usage: "print the points", Aliases: []string{"p"}, Value: false, Destination: &config.printPoints},
      &cli.BoolFlag{Name: "border", Usage: "print a border", Aliases: []string{"b"}, Value: false, Destination: &config.printBorder},
    },
    Action: mksvg,
  }
  app.Run(os.Args)
}

func mksvg(c *cli.Context) error {
  initialize(c)
  parseMask(config.inputPath)
  render()
  return nil
}

func initialize(c *cli.Context) error {
  btu.SetLogLevel(btu.INFO)
  path := c.String("config")
  if !btu.FileExists(path) {
    log.Fatalf("Config file '%s' does not exist.\n", path)
  }
  b, err := ioutil.ReadFile(path)
  btu.CheckError(err)
  err = toml.Unmarshal(b, &config)
  btu.CheckError(err)
  if len(config.OutputDir) == 0 {
    config.OutputDir = "."
  }
  btu.DirMustExist(config.OutputDir)
  if config.PointRadius == 0 {
    config.PointRadius = 2
  }
  if config.MarginEdge == 0 {
    config.MarginEdge = 5
  }
  config.outputPath = filepath.Join(config.OutputDir, c.String("output"))
  return nil
}
