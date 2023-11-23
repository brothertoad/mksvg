package main

import (
  "path"
)

func render() {
  openSvg(path.Join(config.OutputDir, "mask.svg"))
  for _, render := range(mask.Renders) {
    if render.Hide {
      continue
    }
    if render.Comment != "" {
      writeSvgF(`<!-- %s -->`, render.Comment)
      writeSvg("")
    }
    obj := mask.Objects[render.Object]
    xform := createTransformString(render, obj)
    writePathToSvg(obj.d, xform)
  }
  w := mask.Global.Width - 2 * config.MarginEdge
  h := mask.Global.Height - 2 * config.MarginEdge
  if config.printBorder {
    writePlainRectangleToSvg(config.MarginEdge, config.MarginEdge, w, h)
  }
  if config.printGrid {
    spacing := config.GridSpacing
    if spacing == 0 {
      spacing = 25
    }
    writeGridToSvg(config.MarginEdge, config.MarginEdge, w, h, spacing)
  }
  closeSvg()
}
