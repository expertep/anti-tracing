package main

import (
	"anti-tracing/internal"
)

func main() {

	// internal.GetFiles("../test-jaeger/controller")
	internal.ClearDir()
	internal.GetFiles("example")

}
