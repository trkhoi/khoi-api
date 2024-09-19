package main

import (
	"github.com/trkhoi/khoi-api/internal/app/api"
	"github.com/trkhoi/khoi-api/internal/appmain"
)

func main() {
	appmain.RunApplication("api", api.BindService)
}
