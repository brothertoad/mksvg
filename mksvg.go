package main

import (
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "time"
  "github.com/urfave/cli/v2"
  "gopkg.in/yaml.v3"
  "github.com/brothertoad/btu"
)

var config struct {
  OutputDir string `yaml:"outputDir"`
  PointSize int `yaml:"pointSize"`
  StrokeColor string `yaml:"strokeColor"`
  StrokeWidth int `yaml:"strokeWidth"`
}

// Global data
var args struct {
  inputPath string
  outputPath string
  printPoints bool
}

func main() {
  app := &cli.App{
    Name: "mksvg",
    Compiled: time.Now(),
    Usage: "create an SVG file from a text file",
    Flags: []cli.Flag{
      &cli.StringFlag{Name: "config", Usage: "path to configuration file", Required: true, EnvVars: []string{"MKSVG_CONFIG"}},
      &cli.StringFlag{Name: "input", Usage: "input file", Aliases: []string{"i"}, DefaultText: "mask.toml", Value: "mask.toml", Destination: &args.inputPath},
      &cli.StringFlag{Name: "output", Usage: "output file", Aliases: []string{"o"}, DefaultText: "mask.svg", Value: "mask.svg"},
      &cli.BoolFlag{Name: "print-points", Usage: "print the points", Aliases: []string{"p"}, Value: false, Destination: &args.printPoints},
    },
    Action: mksvg,
  }
  app.Run(os.Args)
}

func mksvg(c *cli.Context) error {
  initialize(c)
  loadMask(args.inputPath)
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
  err = yaml.Unmarshal(b, &config)
  btu.CheckError(err)
  if len(config.OutputDir) == 0 {
    config.OutputDir = "."
  }
  btu.DirMustExist(config.OutputDir)
  args.outputPath = filepath.Join(config.OutputDir, c.String("output"))
  // btu.Info("%+v\n", args)
  return nil
}
