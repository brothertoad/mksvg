package main

import (
  "io/ioutil"
  "log"
  "os"
  "time"
  "github.com/urfave/cli/v2"
  "gopkg.in/yaml.v3"
)

var config struct {
  OutputDir string `yaml:"outputDir"`
  TemplateDir string `yaml:"templateDir"`
  // Not really part of config, but a global parameter
  inputPath string
}

func main() {
  app := &cli.App{
    Name: "mksvg",
    Compiled: time.Now(),
    Usage: "create an SVG file from a text file",
    Flags: []cli.Flag{
      &cli.StringFlag{Name: "config", Usage: "path to configuration file", Required: true, EnvVars: []string{"MKSVG_CONFIG"}},
      &cli.StringFlag{Name: "input", Usage: "path to input file", Aliases: []string{"i"}, DefaultText: "mask.toml"},
    },
    Commands: []*cli.Command {
      &initCommand,
      &webCommand,
      // need print command
    },
    Before: initialize,
  }
  app.Run(os.Args)
}

func initialize(c *cli.Context) error {
  path := c.String("config")
  if !fileExists(path) {
    log.Fatalf("Config file '%s' does not exist.\n", path)
  }
  b, err := ioutil.ReadFile(path)
  checkError(err)
  err = yaml.Unmarshal(b, &config)
  checkError(err)
  if len(config.OutputDir) == 0 {
    log.Fatalf("No output directory specified in config file %s.\n", path)
  }
  dirMustExist(config.OutputDir)
  if len(config.TemplateDir) == 0 {
    log.Fatalf("No template directory specified in config file %s.\n", path)
  }
  dirMustExist(config.TemplateDir)
  config.inputPath = c.String("input")
  if len(config.inputPath) == 0 {
    config.inputPath = "mask.toml"
  }
  loadMask(config.inputPath)
  return nil
}
