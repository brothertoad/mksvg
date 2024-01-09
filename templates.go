package main

var htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="mask.css">
    <script src="jquery-3.6.0.slim.js"></script>
    <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
    <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
    <link rel="manifest" href="/site.webmanifest">
    <title>%s mask</title>
</head>
<body>
    <div id="mainContainer">
        <div id="maskContainer">
            <div id="maskDiv">
                <!-- <img src="mask.svg" width="%dpx" height="%dpx"> -->
            </div>
        </div>
        <div id="textContainer">
            <div id="pointContainer">
                <input type="text" id="points" readonly="true">
            </div>
            <div id="buttonContainer">
                <button id="clearButton">Clear</button>
            </div>
        </div>
    </div>

    <script>

    let points = []

    function clickOnMask(e) {
        addPoint(e)
    }

    function addPoint(e) {
        let p = {}
        p.x = e.offsetX
        p.y = e.offsetY
        points.push(p)
        let caption = ""
        points.forEach(function(p, j) {
            if (j != 0) {
                caption = caption + " "
            }
            caption = caption + p.x + "," + p.y
        })
        $("#points").val(caption)
    }

    function clearPoints(e) {
        points = []
        $("#points").val("")
    }

    $(function() {
        $("#maskDiv").click(clickOnMask)
        $("#clearButton").click(clearPoints)
    });
    </script>
</body>
</html>
`

var cssTemplate = `#mainContainer {
    /* display: flex; */
    flex-direction: column;
    width: 100%%;
    height: 100%%;
}

#maskContainer {
    margin-top: 40px;
    width: 100%%;
    height: 80%%;
}

#maskDiv {
    width: %dpx;
    height: %dpx;
    background: url("mask.jpg") no-repeat center;
    margin: auto;
}

#textContainer {
    margin-top: 40px;
}

#pointContainer {
    text-align: center;
}

#buttonContainer {
    padding: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
}

#buttonContainer button {
    align-self: center;
    margin-left: 5px;
    margin-right: 5px;
    width: 150px;
    height: 40px;
}

#points {
    margin: auto;
    width: 80%%;
    height: 2em;
    text-align: center;
    display: inline-block;
}
`

var maskTomlTemplate = `[global]
title = "title"
# printName = "title"
physicalWidth = "120mm"
physicalHeight = "90mm"
width = 1200
height = 900
scale = 1.0
strokeWidth = 1
# strokeColor = "black"

[points]
# p = { x = , y =  }

[objects]
obj.curves = [
""
]
obj.beziers = [
""
]
obj.qbeziers = [
""
]
obj.lines = [
""
]
# obj.scale = 1.0

[[renders]]
object = ""
# comment = ""
translate = { x = , y =  }
# flip = ""
# scale = 1.0
# hide = false
`
