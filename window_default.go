//go:build !darwin

package main

import "fmt"

func SetWindowAspectRatio() {
	fmt.Println("Aspect ratio enforcement not supported on this platform")
}
