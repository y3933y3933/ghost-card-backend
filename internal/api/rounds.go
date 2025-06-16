package api

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/ghost-card/internal/database"
	"github.com/y3933y3933/ghost-card/internal/ws"
)

type RoundsHandler struct {
	logger  *slog.Logger
	queries *database.Queries
	hub     *ws.Hub
}

func NewRoundsHandler(queries *database.Queries, logger *slog.Logger, hub *ws.Hub) *RoundsHandler {
	return &RoundsHandler{
		logger:  logger,
		queries: queries,
		hub:     hub,
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

type CreateRoundRequest struct {
	PlayerID int64 `json:"player_id" binding:"required"`
}

type CreateRoundResponse struct {
	RoundID  int64  `json:"round_id"`
	PlayerID int64  `json:"player_id"`
	Question string `json:"question"`
}

func (h *RoundsHandler) CreateRound(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	var req CreateRoundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid player_id")
		return
	}

	game, err := h.queries.GetGameByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
		} else {
			InternalServerError(c, "db error")
		}
		return
	}

	question, err := h.queries.GetRandomQuestionByLevel(ctx, game.Level)
	if err != nil {
		InternalServerError(c, "failed to pick question")
		return
	}

	round, err := h.queries.CreateRound(ctx, database.CreateRoundParams{
		GameID:          game.ID,
		QuestionID:      question.ID,
		CurrentPlayerID: req.PlayerID,
	})
	if err != nil {
		InternalServerError(c, "failed to create round")
		return
	}

	// ✅ WebSocket 廣播
	// 廣播誰是出題者（全體看到）
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "round_started",
		Data: gin.H{
			"round_id":  round.ID,
			"player_id": round.CurrentPlayerID,
		},
	})

	// 私訊題目給該玩家（只有他看到）
	h.hub.SendToPlayer(game.Code, round.CurrentPlayerID, ws.WebSocketMessage{
		Type: "round_question",
		Data: gin.H{
			"question": question.Content,
		},
	})

	// ✅ 回傳給建立 round 的前端（主持人）
	Success(c, CreateRoundResponse{
		RoundID:  round.ID,
		PlayerID: round.CurrentPlayerID,
		Question: question.Content,
	})

}
