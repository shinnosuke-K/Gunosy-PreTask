package main

import (
	"log"
	"os"

	"github.com/shinnosuke-K/Gunosy-PreTask/controller"

	"github.com/jinzhu/gorm"

	"github.com/shinnosuke-K/Gunosy-PreTask/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *gorm.DB
	Engine *gin.Engine
}

func NewServer() *Server {
	return &Server{
		Engine: gin.Default(),
	}
}

func (router *Server) Init() error {
	db, err := db.Open()
	if err != nil {
		return err
	}
	router.db = db

	ctr := &controller.Information{DB: db}

	//ユーザアカウントを作成
	router.Engine.POST("/signup", ctr.CreateHandler)
	//ユーザ情報を取得
	router.Engine.GET("/users/:user_id", ctr.GetInfoHandler)
	//ユーザ情報を更新
	router.Engine.PATCH("/users/:user_id", ctr.UpdateHandler)
	//ユーザアカウントを削除
	router.Engine.POST("/close", ctr.DeleteHandler)

	return nil

}

func (router *Server) Run(port string) {
	err := router.Engine.Run(port)
	if err != nil {
		return
	}
}

func (router *Server) Close() error {
	return router.db.Close()
}

func main() {
	s := NewServer()
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s.Run(os.Getenv(port))
}
