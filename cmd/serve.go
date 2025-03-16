package cmd

import (
	"game-tournament/internal/controller/game_controller"
	"game-tournament/internal/middleware/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()

		r.GET("/leaderboard", game_controller.Leaderboard())
		r.GET("/player_ranking", game_controller.Rank())
		r.POST("/score", ratelimiter.RateLimitMiddleware(), game_controller.Score())

		r.Run(":9000")
	},
}
