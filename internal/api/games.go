package api

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/y3933y3933/joker/internal/database"
	"github.com/y3933y3933/joker/internal/utils"
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
	if err := bindCreateGameRequest(c, &req); err != nil {
		return
	}

	ctx := c.Request.Context()

	code, err := generateGameCode(ctx, h)
	if err != nil {
		handleGameCodeError(c, h.logger, err)
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

	Success(c, CreateGameResponse{
		ID:        game.ID,
		Code:      game.Code,
		Level:     game.Level,
		CreatedAt: game.CreatedAt.Time,
	})

}

func generateGameCode(ctx context.Context, h *GamesHandler) (string, error) {
	return utils.GenerateUniqueGameCode(ctx, h.queries, 6, 5)
}

func bindCreateGameRequest(c *gin.Context, req *CreateGameRequest) error {
	if err := c.ShouldBindJSON(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			FailedValidation(c, ve.Error())
		} else {
			BadRequest(c, "Invalid request format")
		}
		return err
	}
	return nil
}

func handleGameCodeError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Error("generate unique game code error: ", err)
	if errors.Is(err, utils.ErrGenerateCode) {
		InternalServerError(c, "code collision, try again")
	} else {
		InternalServerError(c, "DB error")
	}
}
