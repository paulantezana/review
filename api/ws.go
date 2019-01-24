package api

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
)

func PublicWs(e *echo.Echo) {
	mel := melody.New()
	mel.Config.MaxMessageSize = 1024 * 1024 * 1024
	ws := e.Group("/api/v1/ws")

	// Comment comment
	ws.GET("/comment", func(c echo.Context) error {
		mel.HandleRequest(c.Response(), c.Request())
		return nil
	})

	mel.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println(msg)
		mel.Broadcast(msg)
	})
}
