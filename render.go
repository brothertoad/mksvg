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
    for _, curve := range(obj.rawCurves) {
      writeCurveToSvg(curve, obj.center, render, xform)
    }
    for _, bezier := range(obj.rawBeziers) {
      writeBezierToSvg(bezier, obj.center, render, xform)
    }
    for _, qbezier := range(obj.rawQBeziers) {
      writeQBezierToSvg(qbezier, obj.center, render, xform)
    }
    for _, line := range(obj.rawLines) {
      writeLineToSvg(line, obj.center, render, xform)
    }
    for _, rect := range(obj.rawRects) {
      writeRectangleToSvg(rect, obj.center, render, xform)
    }
    writeSvg("")
  }
  if config.printBorder {
    w := mask.Global.Width * 10 - 2 * config.MarginEdge
    h := mask.Global.Height * 10 - 2 * config.MarginEdge
    writePlainRectangleToSvg(config.MarginEdge, config.MarginEdge, w, h)
  }
  closeSvg()
}
