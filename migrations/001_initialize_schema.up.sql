CREATE TABLE pack (
    id INT NOT NULL AUTO_INCREMENT,
    status INT NOT NULL,
    title VARCHAR(64) NOT NULL,
    maker VARCHAR(64) NOT NULL,
    description text,
    max_players INT, 

    PRIMARY KEY (id)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE player (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    total_coins INT,

    PRIMARY KEY (id)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE game (
    id BIGINT NOT NULL AUTO_INCREMENT,
    pack_id INT NOT NULL,
    orakki_id VARCHAR(128) NOT NULL,
    orakki_state INT NOT NULL,
    first_player_id BIGINT NOT NULL,
    joined_player_ids VARCHAR(128),
    created_at TIMESTAMP,
    
    PRIMARY KEY (id),
    FOREIGN KEY (pack_id) REFERENCES pack(id),
    FOREIGN KEY (first_player_id) REFERENCES player(id)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE signaling (
    id BIGINT NOT NULL AUTO_INCREMENT,
    token VARCHAR(24) NOT NULL,
    game_id BIGINT NOT NULL,
    player_id BIGINT NOT NULL,
    data VARCHAR(8196),
    created_at TIMESTAMP,

    PRIMARY KEY (id),
    FOREIGN KEY (game_id) REFERENCES game(id)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE feedback (
    id BIGINT NOT NULL AUTO_INCREMENT,
    content VARCHAR(512) NOT NULL,
    created_at TIMESTAMP,

    PRIMARY KEY (id)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;


