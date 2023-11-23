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
    // "unscale" radius
    // This is a little complex, because if the render scale or oject scale is
    // 0.0, treat it like it is 1.0.
    scale := render.Scale
    if scale == 0.0 {
      scale = 1.0
    }
    if obj.Scale != 0.0 {
      scale = scale * obj.Scale
    }
    radius := int(float64(config.PointRadius) / scale)
    writePointsToSvg(obj.points, obj.center, radius, xform)
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
