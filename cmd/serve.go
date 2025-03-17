package cmd

import (
	"game-tournament/internal/controller/game_controller"
	"game-tournament/internal/middleware/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.Mkdir("bin", 0755)
		pid := os.Getpid()
		if err := os.WriteFile("bin/pid", []byte(strconv.Itoa(pid)), 0644); err != nil {
			panic(err)
		}

		r := gin.Default()

		r.GET("/leaderboard", game_controller.Leaderboard())
		r.GET("/player_ranking", game_controller.Rank())
		r.POST("/score", ratelimiter.RateLimitMiddleware(), game_controller.Score())

		r.Run(":9000")
	},
}
