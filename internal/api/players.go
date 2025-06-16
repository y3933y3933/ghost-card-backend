package api

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/ghost-card/internal/database"
)

type PlayersHandler struct {
	logger  *slog.Logger
	queries *database.Queries
}

func NewPlayersHandler(queries *database.Queries, logger *slog.Logger) *PlayersHandler {
	return &PlayersHandler{
		logger:  logger,
		queries: queries,
	}
}

type PlayerResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsHost   bool   `json:"is_host"`
}

func (h *PlayersHandler) ListPlayers(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	players, err := h.queries.ListPlayersByGameCode(ctx, code)
	if err != nil {
		h.logger.Error("list players by game code error: ", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			BadRequest(c, "game not found")
		default:
			InternalServerError(c, "something went wrong")
		}
		return
	}

	Success(c, transformToPlayerResponse(players))
}

func transformToPlayerResponse(players []database.ListPlayersByGameCodeRow) []PlayerResponse {
	var playerResponses []PlayerResponse
	for _, p := range players {
		playerResponses = append(playerResponses, PlayerResponse{
			ID:       p.ID,
			Nickname: p.Nickname,
			IsHost:   p.IsHost.Bool,
		})
	}
	return playerResponses
}
