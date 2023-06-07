package main

import (
  "bufio"
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
      &cli.StringFlag{Name: "list", Usage: "file containing list of input files", Aliases: []string{"l"}},
      &cli.StringFlag{Name: "output", Usage: "output file", Aliases: []string{"o"}, DefaultText: "mask.svg", Value: "mask.svg"},
      &cli.BoolFlag{Name: "points", Usage: "print the points", Aliases: []string{"p"}, Value: false, Destination: &config.printPoints},
      &cli.BoolFlag{Name: "border", Usage: "print a border", Aliases: []string{"b"}, Value: false, Destination: &config.printBorder},
      &cli.BoolFlag{Name: "grid", Usage: "print a grid", Aliases: []string{"g"}, Value: false, Destination: &config.printGrid},
      &cli.StringFlag{Name: "image", Usage: "set background image"},
    },
    Action: mksvg,
  }
  app.Run(os.Args)
}

func mksvg(c *cli.Context) error {
  initialize(c)
  args := getArgs(c)
  parseMask(args)
  if c.String("image") != "" {
    initFromImage(c.String("image"))
  }
  render()
  return nil
}

func getArgs(c *cli.Context) []string {
  args := c.Args().Slice()
  listPath := c.String("list")
  if listPath != "" {
    f := btu.OpenFile(listPath)
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
      args = append(args, scanner.Text())
    }
  }
  if len(args) == 0 {
    // default is mask.toml
    args = append(args, "mask.toml")
  }
  return args
}

func initialize(c *cli.Context) {
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
}
