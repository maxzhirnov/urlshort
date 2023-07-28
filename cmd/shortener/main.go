package main

import (
	"github.com/maxzhirnov/urlshort/internal/app"
	data "github.com/maxzhirnov/urlshort/internal/data/inmemory"
)

func main() {
	safeMap := data.NewSafeMap()
	store := data.NewInMemoryStore(safeMap)
	urlShortener := app.NewURLShortener(store)
	server := NewServer(urlShortener)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
