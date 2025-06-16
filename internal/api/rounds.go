package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/ghost-card/internal/database"
)

type RoundsHandler struct {
	logger  *slog.Logger
	queries *database.Queries
}

func NewRoundsHandler(queries *database.Queries, logger *slog.Logger) *RoundsHandler {
	return &RoundsHandler{
		logger:  logger,
		queries: queries,
	}
}

func (h *RoundsHandler) CreateRoundInOrder(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

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
	if err != nil || len(players) == 0 {
		BadRequest(c, "no players found")
		return
	}

	var nextPlayerID int64

	// 取得上一 round，如果沒有就從第一位開始
	lastRound, err := h.queries.GetLatestRoundByGameID(ctx, game.ID)
	if err != nil {
		// no previous round → pick first player
		nextPlayerID = players[0].ID
	} else {
		// 找目前玩家在 list 的 index
		currentIndex := -1
		for i, p := range players {
			if p.ID == lastRound.CurrentPlayerID {
				currentIndex = i
				break
			}
		}
		// 決定下一位玩家 index（環狀）
		nextIndex := (currentIndex + 1) % len(players)
		nextPlayerID = players[nextIndex].ID
	}

	question, err := h.queries.GetUnusedQuestion(ctx, database.GetUnusedQuestionParams{
		Level:  game.Level,
		GameID: game.ID,
	})
	if err != nil {
		InternalServerError(c, "no questions available")
		return
	}

	round, err := h.queries.CreateRound(ctx, database.CreateRoundParams{
		GameID:          game.ID,
		QuestionID:      question.ID,
		CurrentPlayerID: nextPlayerID,
	})
	if err != nil {
		InternalServerError(c, "failed to create round")
		return
	}

	c.JSON(http.StatusOK, round)
}
