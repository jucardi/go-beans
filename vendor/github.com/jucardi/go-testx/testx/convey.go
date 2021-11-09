package testx

import (
	"fmt"
	"github.com/jucardi/go-terminal-colors"
	"testing"
)

var (
	descriptionColors = map[int][]fmtc.Color{
		0: {fmtc.Cyan},
		1: {fmtc.Yellow, fmtc.Bold},
		2: {fmtc.Yellow, fmtc.Bold},
		3: {fmtc.White},
		4: {fmtc.Cyan},
	}

	assertionColors = map[int][]fmtc.Color{
		0: {fmtc.Green},
	}
)

func Convey(description string, t *testing.T, f func()) {
	ctx := newCtx(t)
	contextPile = append(contextPile, ctx)
	if len(contextPile) != previous {
		ctx.Println("")
	}
	ctx.Println("").PrintIndent(description+" ", descriptionColor()...)
	f()
	if len(contextPile) <= 2 {
		ctx.Println("\n").Println(fmt.Sprintf("%d total assertions", ctx.assertions), assertionsColor()...).Println("")
	}
	if len(contextPile) < previous {
		ctx.Println("")
	}
	previous = len(contextPile)
	contextPile = contextPile[:len(contextPile)-1]
	currentCtx().assertions += ctx.assertions
}

func descriptionColor() []fmtc.Color {
	if v, ok := descriptionColors[len(contextPile)]; ok {
		return v
	}
	return descriptionColors[0]
}

func assertionsColor() []fmtc.Color {
	if v, ok := assertionColors[len(contextPile)]; ok {
		return v
	}
	return assertionColors[0]
}
