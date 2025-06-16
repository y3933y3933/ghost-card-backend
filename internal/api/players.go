package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"

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

func (h *PlayersHandler) JoinGameHandler(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")

	game, err := h.queries.GetGameByCode(ctx, gameCode)
	if err != nil {
		h.logger.Error("get game by code error: ", err)
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
			return
		}
		InternalServerError(c, "DB error")
		return
	}

	var req struct {
		Nickname string `json:"nickname" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			FailedValidation(c, ve.Error())
		} else {
			BadRequest(c, "Invalid request format")
		}
		return
	}

	count, err := h.queries.CountPlayersInGame(ctx, game.ID)
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}

	isHost := count == 0

	player, err := h.queries.CreatePlayer(ctx, database.CreatePlayerParams{
		GameID:   game.ID,
		Nickname: req.Nickname,
		IsHost: pgtype.Bool{
			Bool:  isHost,
			Valid: true,
		},
	})

	if err != nil {
		h.logger.Error("failed to join game: ", err)
		InternalServerError(c, "failed to join game")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        player.ID,
		"nickname":  player.Nickname,
		"is_host":   player.IsHost,
		"joined_at": player.JoinedAt,
	})
}
