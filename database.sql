CREATE TABLE match_results (
    id SERIAL PRIMARY KEY,
    competition_id VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    team_1 VARCHAR(64) NOT NULL,
    team_1_score INTEGER NOT NULL,
    team_2 VARCHAR(64) NOT NULL,
    team_2_score INTEGER NOT NULL
);
