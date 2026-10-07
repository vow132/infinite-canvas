package main

import (
	"log"

	"github.com/vow132/infinite-canvas/config"
	"github.com/vow132/infinite-canvas/handler"
	"github.com/vow132/infinite-canvas/router"
	"github.com/vow132/infinite-canvas/service"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	if err := service.EnsureDefaultAdmin(); err != nil {
		log.Fatal(err)
	}
	if err := service.EnsureDefaultAgentSkills(); err != nil {
		log.Fatal(err)
	}
	service.StartPromptSyncScheduler()
	service.StartCanvasProjectCleanupScheduler()
	handler.StartVideoTaskPoller()
	log.Fatal(router.New().Run(":" + config.Cfg.Port))
}
