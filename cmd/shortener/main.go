package main

import (
	"github.com/maxzhirnov/urlshort/internal/app"
	data "github.com/maxzhirnov/urlshort/internal/data/inmemory"
)

func main() {
	//urlData := make(map[string]models.URL) // переделать в thread safe
	m := data.NewSafeMap()
	store := data.NewInMemoryStore(m)
	urlShortener := app.NewURLShortener(store)
	server := NewServer(urlShortener)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
