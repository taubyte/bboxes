package lib

import (
	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
	"github.com/taubyte/go-sdk/http/client"
)

//export dbping
func dbping(e event.Event) uint32 {
	db, err := database.New("/db")
	if err != nil {
		return 1
	}
	db.Put("k", []byte("v"))
	c, _ := client.New()
	_ = c
	return 0
}
