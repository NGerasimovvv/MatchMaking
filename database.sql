CREATE TABLE players (
                         id SERIAL PRIMARY KEY,
                         name TEXT NOT NULL,
                         skill FLOAT NOT NULL,
                         latency FLOAT NOT NULL,
                         join_time TIMESTAMP NOT NULL
);
