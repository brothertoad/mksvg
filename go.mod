module brothertoad.net/mksvg

go 1.18

require (
	github.com/brothertoad/bezier v0.0.0
	github.com/brothertoad/btu v1.2.0
	github.com/pelletier/go-toml v1.9.5
	github.com/urfave/cli/v2 v2.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
)

replace (
	github.com/brothertoad/bezier => ../bezier
)
