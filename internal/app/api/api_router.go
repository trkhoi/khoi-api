package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/trkhoi/khoi-api/internal/appmain"
)

func setupRouter(p *appmain.Params) *gin.Engine {
	r := gin.New()

	r.Use(
		gin.LoggerWithWriter(p.Logger().Writer(), "/healthz"),
		gin.Recovery(),
	)

	r.Use(func(c *gin.Context) {
		allowOrigins := []string{"*"}

		cors.New(
			cors.Config{
				AllowOrigins: allowOrigins,
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
				AllowHeaders: []string{"Origin", "Host",
					"Content-Type", "Content-Length",
					"Accept-Encoding", "Accept-Language", "Accept",
					"X-CSRF-Token", "Authorization", "X-Requested-With", "X-Access-Token", "X-Request-Id"},
				ExposeHeaders:    []string{"MeAllowMethodsntent-Length"},
				AllowCredentials: true,
			},
		)(c)
	})

	// r.Use(func(ctx *gin.Context) {
	// 	cr := mdwgin.CaptureRequest(ctx, &mdwgin.CaptureRequestOptions{
	// 		ExcludePaths: []string{"/healthz"},
	// 	})
	// 	if cr == nil {
	// 		ctx.Next()
	// 		return
	// 	}
	// 	b, err := json.Marshal(cr)
	// 	if err != nil {
	// 		p.Logger().Error(err, "cannot marshal capture request")
	// 		ctx.Next()
	// 		return
	// 	}
	// 	go func() {
	// 		kfkMsg := typesetqueue.KafkaMessage{
	// 			Type:   typesetqueue.KAFKA_MESSAGE_TYPE_AUDIT,
	// 			Data:   b,
	// 			Sender: typesetservice.SERVICE_MOCHI_API,
	// 		}
	// 		body, err := json.Marshal(kfkMsg)
	// 		if err != nil {
	// 			return
	// 		}
	// 		requestutils.SendRequest(requestutils.SendRequestQuery{
	// 			Method: http.MethodPost,
	// 			URL:    fmt.Sprintf("%v/api/v1/audit", p.Config().GetString("MOCHI_AUDIT_BASE_URL")),
	// 			Headers: map[string]string{
	// 				"Content-Type": "application/json",
	// 			},
	// 			Body: bytes.NewReader(body),
	// 		})
	// 	}()
	// })

	// handlers
	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	// use ginSwagger middleware to serve the API docs
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
