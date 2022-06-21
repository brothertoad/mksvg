package main

import (
  "os"
)

var svgPrefix =`<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="WIDTHpx" height="HEIGHTpx" viewBox="0 0 WIDTH HEIGHT" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink= "http://www.w3.org/1999/xlink">
  <style>
  * {
    fill: none;
    stroke: STROKE-COLOR;
    stroke-width: STROKE-WIDTH;
  }
  </style>
`
var svgSuffix = `</svg>
`

var svgFile *os.File

func openSvg(path string) {
  file, err := os.Create(path)
  checkError(err)
  svgFile = file
  svgFile.WriteString(filterString(svgPrefix))
}

func writeSvg(s string) {
  svgFile.WriteString(s + "\n")
}

func closeSvg() {
  svgFile.WriteString(filterString(svgSuffix))
  err := svgFile.Close()
  checkError(err)
}
