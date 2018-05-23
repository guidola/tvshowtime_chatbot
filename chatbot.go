package main

import (

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/guidola/tvshowtime_chatbot/engine"
	"bytes"
)

var (
	upgrader = websocket.Upgrader{}
	bot = engine.CreateBotInstance()
)



func ws_handler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	//send initial message
	err = ws.WriteMessage(websocket.TextMessage, bot.Init())

	for {
		// wait for user input. then process it and send response back
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}

		err = ws.WriteMessage(websocket.TextMessage, bot.ProcessInput(bytes.NewBuffer(msg).String()))
		if err != nil {
			c.Logger().Error(err)
		}
		c.Logger().Info("Received msg: %s\n", msg)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./public/")
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.GET("/ws", ws_handler)
	e.Logger.Fatal(e.Start("127.0.0.1:1714"))
}