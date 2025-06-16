package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
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

func (h *RoundsHandler) DrawCard(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	var req struct {
		PlayerID int64 `json:"player_id" binding:"required"`
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

	round, err := h.queries.GetLatestRoundByGameID(ctx, game.ID)
	if err != nil {
		NotFound(c, "no current round")
		return
	}

	if round.CurrentPlayerID != req.PlayerID {
		Forbidden(c, "not your turn")
		return
	}

	if round.Status != "pending" {
		BadRequest(c, "round already resolved")
		return
	}

	isJoker := rand.Intn(3) == 0

	err = h.queries.RevealRound(ctx, database.RevealRoundParams{
		ID: round.ID,
		IsJoker: pgtype.Bool{
			Bool:  isJoker,
			Valid: true,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reveal round"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_joker":    isJoker,
		"question":    round.QuestionID,
		"next_action": "next_round",
	})

}

func (h *RoundsHandler) GetCurrentRound(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	round, err := h.queries.GetCurrentRoundByGameCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "no round found")
		} else {
			InternalServerError(c, "failed to fetch round")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"round_id": round.RoundID,
		"question": round.Question,
		"current_player": gin.H{
			"id":       round.PlayerID,
			"nickname": round.Nickname,
		},
		"is_joker": round.IsJoker,
		"status":   round.Status,
	})
}
