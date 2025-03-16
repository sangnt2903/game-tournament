package game_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type gameController struct {
	client *redis.Client
}

var _gameController *gameController

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,    // use default DB
		PoolSize: 1000, // 1000 connections x 100 ms /
	})

	_gameController = &gameController{
		client: client,
	}
}

func Score() gin.HandlerFunc {
	return func(c *gin.Context) {
		type request struct {
			Username string `json:"username"`
			Score    int    `json:"score"`
		}

		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := _gameController.client.ZIncrBy(c, "leaderboard", float64(req.Score), req.Username).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Score added successfully"})
	}
}

func Leaderboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		top := c.DefaultQuery("top", "10")
		topK, _ := strconv.ParseInt(top, 10, 64)

		username := c.Query("username")
		ctx := c.Request.Context()

		// Prepare variables to store results
		var leaderboard []redis.Z
		var rank int64
		var score float64

		// Run Redis pipeline
		_, err := _gameController.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			leaderboardCmd := pipe.ZRevRangeWithScores(ctx, "leaderboard", 0, topK-1)
			rankCmd := pipe.ZRevRank(ctx, "leaderboard", username)
			scoreCmd := pipe.ZScore(ctx, "leaderboard", username)

			pipe.Exec(ctx)

			// Execute the pipeline
			var iErr error
			leaderboard, iErr = leaderboardCmd.Result()
			if iErr != nil && iErr != redis.Nil {
				return iErr
			}

			rank, iErr = rankCmd.Result()
			if iErr != nil {
				if iErr != redis.Nil {
					return iErr
				}
				rank = -1
			}

			score, iErr = scoreCmd.Result()
			if iErr != nil {
				if iErr != redis.Nil {
					return iErr
				}
				score = 0
			}

			return nil
		})
		if err != nil && err != redis.Nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(leaderboard) == 0 {
			c.JSON(200, gin.H{"data": gin.H{
				"leaderboard": []*Ranking{},
				"user_ranking": &Ranking{
					Username: username,
					Score:    score,
					Rank:     rank,
				},
			}})
			return
		}

		var result []*Ranking
		for i := 0; i < len(leaderboard); i++ {
			result = append(result, &Ranking{
				Username: leaderboard[i].Member.(string),
				Score:    leaderboard[i].Score,
				Rank:     int64(i + 1),
			})
		}

		c.JSON(200, gin.H{
			"data": gin.H{
				"leaderboard": result,
				"user_ranking": &Ranking{
					Username: username,
					Score:    score,
					Rank: func() int64 {
						if rank >= 0 {
							return rank + 1
						}
						return -1
					}(),
				},
			},
		})
	}
}

func Rank() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")

		var (
			rank  int64
			score float64
		)

		_, err := _gameController.client.Pipelined(c, func(pipe redis.Pipeliner) error {
			rankCmd := pipe.ZRevRank(c, "leaderboard", username)
			scoreCmd := pipe.ZScore(c, "leaderboard", username)

			var iErr error
			_, iErr = pipe.Exec(c)

			rank, iErr = rankCmd.Result()
			if iErr != nil {
				if iErr != redis.Nil {
					return iErr
				}
				rank = -1
			}

			score, iErr = scoreCmd.Result()
			if iErr != nil {
				if iErr != redis.Nil {
					return iErr
				}
				score = 0
			}

			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"data": gin.H{
				"rank": func() int64 {
					if rank >= 0 {
						return rank + 1
					}
					return -1
				}(),
				"score": score,
			},
		})
	}
}

type Ranking struct {
	Username string  `json:"username"`
	Score    float64 `json:"score"`
	Rank     int64   `json:"rank"`
}

func newRanking(username string, score float64, rank int64) Ranking {
	return Ranking{
		Username: username,
		Score:    score,
		Rank:     rank,
	}
}
