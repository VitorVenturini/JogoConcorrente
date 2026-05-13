package main

import (
	"log"
	"runtime/debug"
)

func logPanic(component string) {
	if r := recover(); r != nil {
		log.Printf("panic in %s: %v\n%s", component, r, debug.Stack())
	}
}
