package main

import (
	"fmt"

	"github.com/thearyanahmed/nordsec/core/handler"
)

func main() {
	fmt.Println("[cmd] init core")
	handler.NewRouter()
}
