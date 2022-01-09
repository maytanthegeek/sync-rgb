package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/urfave/cli/v2"
)

func toggleSwitch() {
	rawState := hd.Strip.Status("20")
	state := false
	if rawState != nil {
		state = rawState.(bool)
	}
	hd.Strip.SetW("20", !state, 5*time.Second)
}

func changeColor(color string) cli.ExitCoder {
	rgb := strings.SplitN(color, ",", 3)
	if len(rgb) != 3 {
		return cli.Exit("invalid color", 1)
	}
	r, errR := strconv.ParseFloat(rgb[0], 32)
	g, errG := strconv.ParseFloat(rgb[1], 32)
	b, errB := strconv.ParseFloat(rgb[2], 32)

	if errR != nil || errG != nil || errB != nil {
		return cli.Exit("invalid color", 1)
	}

	c := colorful.Color{R: r / 255.0, G: g / 255.0, B: b / 255.0}
	h, s, v := c.Hsv()
	hInt := int(h)
	sInt := int(s * 1000)
	vInt := int(v * 1000)

	color = fmt.Sprintf("%04x%04x%04x", hInt, sInt, vInt)
	hd.Strip.SetW("24", color, 5*time.Second)

	return nil
}
