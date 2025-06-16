package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/y3933y3933/ghost-card/internal/database"
	"github.com/y3933y3933/ghost-card/internal/utils"
)

type GamesHandler struct {
	logger  *slog.Logger
	queries *database.Queries
}

func NewGamesHandler(queries *database.Queries, logger *slog.Logger) *GamesHandler {
	return &GamesHandler{
		logger:  logger,
		queries: queries,
	}
}

type CreateGameRequest struct {
	Level string `json:"level" binding:"required,oneof=easy normal spicy"`
}

type CreateGameResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *GamesHandler) CreateGame(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			FailedValidation(c, ve.Error())
		} else {
			BadRequest(c, "Invalid request format")
		}
		return
	}

	ctx := c.Request.Context()

	code, err := utils.GenerateUniqueGameCode(ctx, h.queries, 6, 5)

	if err != nil {
		h.logger.Error("generate unique game code error: ", err)
		if errors.Is(err, utils.ErrGenerateCode) {
			InternalServerError(c, "code collision, try again")
		} else {
			InternalServerError(c, "DB error")
		}
		return
	}

	game, err := h.queries.CreateGame(ctx, database.CreateGameParams{
		Code:   code,
		Level:  req.Level,
		Status: "waiting",
	})

	if err != nil {
		h.logger.Error("create game fail: ", err)
		InternalServerError(c, "failed to create game")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": CreateGameResponse{
			ID:        game.ID,
			Code:      game.Code,
			Level:     game.Level,
			CreatedAt: game.CreatedAt.Time,
		},
	})
}

type PlayerResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IsHost   bool   `json:"is_host"`
}

type GameResponse struct {
	ID      int64            `json:"id"`
	Code    string           `json:"code"`
	Level   string           `json:"level"`
	Status  string           `json:"status"`
	Players []PlayerResponse `json:"players"`
}

func (h *GamesHandler) GetGameByCode(c *gin.Context) {
	code := c.Param("code")
	ctx := c.Request.Context()

	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		h.logger.Error("get game by code error: ", err)
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
		} else {
			InternalServerError(c, "DB error")
		}
		return
	}

	players, err := h.queries.ListPlayersByGame(ctx, game.ID)
	if err != nil {
		h.logger.Error("failed to get players: ", err)
		InternalServerError(c, "failed to get players")
		return
	}

	var playerResponses []PlayerResponse
	for _, p := range players {
		playerResponses = append(playerResponses, PlayerResponse{
			ID:       p.ID,
			Nickname: p.Nickname,
			IsHost:   p.IsHost.Bool,
		})
	}

	resp := GameResponse{
		ID:      game.ID,
		Code:    game.Code,
		Level:   game.Level,
		Status:  game.Status,
		Players: playerResponses,
	}

	c.JSON(http.StatusOK, resp)
}
