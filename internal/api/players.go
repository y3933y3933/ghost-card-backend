package api

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/y3933y3933/joker/internal/database"
	"github.com/y3933y3933/joker/internal/ws"
)

type PlayersHandler struct {
	logger  *slog.Logger
	queries *database.Queries
	hub     *ws.Hub
}

func NewPlayersHandler(queries *database.Queries, logger *slog.Logger, hub *ws.Hub) *PlayersHandler {
	return &PlayersHandler{
		logger:  logger,
		queries: queries,
		hub:     hub,
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

type JoinGameRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

func (h *PlayersHandler) JoinGame(c *gin.Context) {
	ctx := c.Request.Context()
	gameCode := c.Param("code")

	var req JoinGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid nickname")
		return
	}

	game, err := h.queries.GetGameByCode(ctx, gameCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			NotFound(c, "game not found")
		} else {
			InternalServerError(c, "DB error")
		}
		return
	}

	count, err := h.queries.CountPlayersInGame(ctx, game.ID)
	if err != nil {
		InternalServerError(c, "Count error")
		return
	}
	isHost := count == 0

	player, err := h.queries.CreatePlayer(ctx, database.CreatePlayerParams{
		GameID:   game.ID,
		Nickname: req.Nickname,
		IsHost:   pgtype.Bool{Bool: isHost, Valid: true},
	})
	if err != nil {
		InternalServerError(c, "Create player failed")
		return
	}

	// ✅ WebSocket 廣播
	h.hub.BroadcastToGame(game.Code, ws.WebSocketMessage{
		Type: "player_joined",
		Data: gin.H{
			"id":        player.ID,
			"nickname":  player.Nickname,
			"is_host":   player.IsHost,
			"joined_at": player.JoinedAt,
		},
	})

	// ✅ 回傳該玩家資訊
	Success(c, PlayerResponse{
		ID:       player.ID,
		Nickname: player.Nickname,
		IsHost:   player.IsHost.Bool,
	})

}
