package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lakeside763/github-repo/internal/adapters/repository"
	"github.com/lakeside763/github-repo/internal/adapters/http/handlers"
)

func RepometadataRoutes(router *httprouter.Router, repo *repository.DataRepository) {
	repoHandler := handlers.NewRepometadataHandler(repo.Repometadata)

	router.POST("/repositories", repoHandler.CreateRepometadata)
}