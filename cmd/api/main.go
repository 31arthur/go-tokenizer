package main

import (
	"fmt"
	"server/internal/auth"
	"server/internal/server"
)

func main() {

	auth.NewAuth()
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
