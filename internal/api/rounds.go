package api

import (
	"database/sql"
	"errors"
	"log/slog"

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

type CurrentPlayer struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}
type CurrentRoundResponse struct {
	RoundID         int64  `json:"round_id"`
	Question        string `json:"question_content"`
	GameID          int64  `json:"game_id"`
	IsJoker         *bool  `json:"is_joker"`
	Status          string `json:"status"`
	Level           string `json:"level"`
	CurrentPlayerID int64  `json:"current_player_id"`
}

func (h *RoundsHandler) GetCurrentRound(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")

	round, err := h.queries.GetCurrentRoundByGameCode(ctx, gameCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "no round found for this game")
		} else {
			h.logger.Error("get current round failed", err)
			InternalServerError(c, "failed to get current round")
		}
		return
	}

	Success(c, CurrentRoundResponse{
		RoundID:         round.ID,
		Question:        round.QuestionContent,
		GameID:          round.GameID,
		IsJoker:         &round.IsJoker.Bool,
		Status:          round.Status,
		Level:           round.Level,
		CurrentPlayerID: round.CurrentPlayerID,
	})

}
