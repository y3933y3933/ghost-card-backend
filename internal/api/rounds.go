package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/y3933y3933/ghost-card/internal/database"
	"github.com/y3933y3933/ghost-card/internal/utils"
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

func (h *RoundsHandler) DrawCard(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")

	roundIDParam := c.Param("id")
	roundID, err := utils.ParseID(roundIDParam)

	if err != nil {
		BadRequest(c, "invalid round id")
		return
	}

	game, err := h.queries.GetGameByCode(ctx, gameCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
			return
		}
		InternalServerError(c, "db error")
		return
	}

	round, err := h.queries.GetRoundByID(ctx, roundID)
	if err != nil {
		InternalServerError(c, "round not found")
		return
	}

	if round.Status != "pending" {
		BadRequest(c, "round already completed")
		return
	}

	// 🎲 抽鬼牌（1/3 機率）
	isJoker := rand.Intn(3) == 0

	newStatus := "done"
	if isJoker {
		newStatus = "revealed"
	}

	err = h.queries.UpdateRoundStatus(ctx, database.UpdateRoundStatusParams{
		ID: round.ID,
		IsJoker: pgtype.Bool{
			Bool:  isJoker,
			Valid: true,
		},
		Status: newStatus,
	})
	if err != nil {
		InternalServerError(c, "failed to update round")
		return
	}

	question, err := h.queries.GetQuestionByID(ctx, round.QuestionID)
	if err != nil {
		InternalServerError(c, "failed to get question")
		return
	}

	if isJoker {
		// 👻 廣播給所有人：顯示題目
		h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
			Type: "joker_revealed",
			Data: gin.H{
				"round_id":  round.ID,
				"player_id": round.CurrentPlayerID,
				"question":  question,
			},
		})
	} else {
		// 🛡 廣播回合結束（安全）
		h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
			Type: "player_safe",
			Data: gin.H{
				"round_id":  round.ID,
				"player_id": round.CurrentPlayerID,
			},
		})
	}

	c.Status(http.StatusOK)
}

func (h *RoundsHandler) CreateNextRound(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")

	game, err := h.queries.GetGameByCode(ctx, gameCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
			return
		}
		InternalServerError(c, "db error")
		return
	}

	players, err := h.queries.ListPlayersByGameCode(ctx, game.Code)
	if err != nil || len(players) == 0 {
		InternalServerError(c, "failed to get players")
		return
	}

	lastRound, err := h.queries.GetLatestRoundInGame(ctx, game.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		InternalServerError(c, "failed to get last round")
		return
	}

	// 決定下一位玩家
	nextPlayerID := players[0].ID
	if lastRound.ID != 0 {
		for i, p := range players {
			if p.ID == lastRound.CurrentPlayerID {
				nextPlayerID = players[(i+1)%len(players)].ID
				break
			}
		}
	}

	question, err := h.queries.GetRandomQuestionByLevel(ctx, game.Level)
	if err != nil {
		InternalServerError(c, "failed to get question")
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

	// 廣播回合開始（不含題目）
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "round_started",
		Data: gin.H{
			"round_id":  round.ID,
			"player_id": nextPlayerID,
		},
	})

	// 私訊題目給當事人
	h.hub.SendToPlayer(game.Code, nextPlayerID, ws.WebSocketMessage{
		Type: "round_question",
		Data: gin.H{
			"question": question.Content,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"round_id":  round.ID,
		"player_id": nextPlayerID,
	})
}

func (h *RoundsHandler) EndGame(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")

	game, err := h.queries.GetGameByCode(ctx, gameCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
			return
		}
		InternalServerError(c, "db error")
		return
	}

	err = h.queries.UpdateGameStatus(ctx, database.UpdateGameStatusParams{
		ID:     game.ID,
		Status: "ended",
	})
	if err != nil {
		InternalServerError(c, "failed to end game")
		return
	}

	// 廣播遊戲結束
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "game_ended",
		Data: gin.H{
			"game_id": game.ID,
		},
	})

	c.Status(http.StatusOK)
}
