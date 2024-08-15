CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    skill REAL NOT NULL,
    latency REAL NOT NULL,
    searching_match BOOLEAN NOT NULL,
    search_start_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

create index user_name on users (name);
create index user_searching_match on users (searching_match);