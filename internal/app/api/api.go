package api

import (
	"net/http"

	"github.com/gin-contrib/pprof"

	"github.com/trkhoi/khoi-api/internal/app/handler"
	"github.com/trkhoi/khoi-api/internal/app/routes"
	"github.com/trkhoi/khoi-api/internal/appmain"
)

// BindService creates the backend service and binds it to the serving harness
func BindService(p *appmain.Params, b *appmain.Bindings) appmain.IServer {
	router := setupRouter(p)
	h := handler.New(p)
	routes.LoadV1(router, h, p)
	if !p.Config().IsSet("PORT") {
		p.Logger().Fatal("PORT environment variable not set")
	}
	pprof.Register(router)

	srv := &http.Server{
		Addr:    ":" + p.Config().GetString("PORT"),
		Handler: router,
	}
	return srv
}
