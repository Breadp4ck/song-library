-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs (
    song_id UUID DEFAULT gen_random_uuid(),
    song_name VARCHAR NOT NULL,
    song_text TEXT,
    group_name VARCHAR NOT NULL,
    link VARCHAR,
    release_date DATE,
    PRIMARY KEY (song_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE songs;
-- +goose StatementEnd
