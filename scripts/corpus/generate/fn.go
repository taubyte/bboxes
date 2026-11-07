package lib

import "github.com/taubyte/go-sdk/event"

//go:generate sh -c "echo package lib > gen.go; echo const Gen = 1 >> gen.go"

//export g
func g(e event.Event) uint32 { return uint32(Gen) }
