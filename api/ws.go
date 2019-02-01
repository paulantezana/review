package api

import (
	"github.com/labstack/echo"
	"github.com/olahol/melody"
    "net/http"
)

func PublicWs(e *echo.Echo) {
	ws := e.Group("/api/v1/ws")

	// Comment comment
	mel := melody.New()
	mel.Config.MaxMessageSize = 1024 * 1024 * 1024
	mel.Upgrader.CheckOrigin = func(r *http.Request) bool {
       return  true
    }
	// Path
	ws.GET("/comment", func(c echo.Context) error {
		mel.HandleRequest(c.Response(), c.Request())
		return nil
	})
	// Handle message
	mel.HandleMessage(func(s *melody.Session, msg []byte) {
		mel.Broadcast(msg)
	})


    me := melody.New()
    me.Config.MaxMessageSize = 1024 * 1024 * 1024
    me.Upgrader.CheckOrigin = func(r *http.Request) bool {
        return  true
    }

	// Chat
    ws.GET("/chat", func(c echo.Context) error {
        me.HandleRequest(c.Response(), c.Request())
        return nil
    })

    // Handle message
    me.HandleMessage(func(s *melody.Session, msg []byte) {
    	me.Broadcast(msg)
    })
}
