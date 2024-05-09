package main

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestHello(t *testing.T) {
	g := Goblin(t)
	g.Describe("Hello", func() {
		g.It("Should print hello world", func() {
			g.Assert(hello()).Equal("Hello World!")
		})
	})
	g.Describe("Hello your name", func() {
		g.It("Should print hello your name", func() {
			g.Assert(helloYourName("hi")).Equal("Hello John!")
		})
	})
}
