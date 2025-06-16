-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rounds (
    id BIGSERIAL PRIMARY KEY,
    game_id BIGSERIAL NOT NULL,
    question_id BIGSERIAL NOT NULL,
    current_player_id BIGSERIAL NOT NULL,
    is_joker BOOLEAN DEFAULT FALSE,
    status TEXT NOT NULL CHECK (status IN ('pending','revealed','done')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id),
    FOREIGN KEY (current_player_id) REFERENCES players(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rounds;
-- +goose StatementEnd
