package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lakeside763/github-repo/internal/adapters/repository"
)

func SetupRoutes(router *httprouter.Router, repo *repository.DataRepository) {
	RepometadataRoutes(router, repo)
}