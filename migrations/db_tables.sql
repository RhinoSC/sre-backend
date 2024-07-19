-- Eliminar tablas si existen
DROP TABLE IF EXISTS donation_bids;
DROP TABLE IF EXISTS bid_options;
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS run_metadata;
DROP TABLE IF EXISTS runs;
DROP TABLE IF EXISTS prizes;
DROP TABLE IF EXISTS user_socials;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams_runs;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS donations;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS events;

-- Migration for creating the 'users' table
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Migration for creating the 'user_socials' table
CREATE TABLE user_socials (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    twitch VARCHAR(255) DEFAULT NULL,
    twitter VARCHAR(255) DEFAULT NULL,
    youtube VARCHAR(255) DEFAULT NULL,
    facebook VARCHAR(255) DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Migration for creating the 'prizes' table
CREATE TABLE prizes (
    id VARCHAR(255) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    url VARCHAR(255),
    min_amount NUMERIC DEFAULT 0,
    status VARCHAR(255) NOT NULL,
    international_delivery BOOLEAN DEFAULT FALSE,
    event_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- Migration for creating the 'runs' table
CREATE TABLE runs (
    id VARCHAR(255) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    start_time_mili NUMERIC NOT NULL,
    estimate_string VARCHAR(255),
    estimate_milliseconds NUMERIC,
    metadata TEXT,
    schedule_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);

-- Migration for creating the 'run_metadata' table
CREATE TABLE run_metadata (
    id VARCHAR(255) PRIMARY KEY,
    run_id VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    platform VARCHAR(255) NOT NULL,
    twitch_game_name VARCHAR(255),
    run_name VARCHAR(255),
    note TEXT DEFAULT NULL,
    FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
);

-- Migration for creating the 'teams' table
CREATE TABLE teams (
    id VARCHAR(255) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL
);

-- Migration for creating the 'players' table
CREATE TABLE players (
    team_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Migration for creating the 'bids' table
CREATE TABLE bids (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    bidname VARCHAR(255) NOT NULL,
    goal NUMERIC DEFAULT 0,
    current_amount NUMERIC NOT NULL,
    `description` TEXT,
    type VARCHAR(255) NOT NULL CHECK (type IN ('bidwar', 'total', 'goal')),
    create_new_options BOOLEAN NOT NULL,
    run_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
);

-- Migration for creating the 'bid_options' table
CREATE TABLE bid_options (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    bid_id VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    current_amount NUMERIC NOT NULL DEFAULT 0,
    FOREIGN KEY (bid_id) REFERENCES bids(id) ON DELETE CASCADE
);

-- Migration for creating the 'donations' table
CREATE TABLE donations (
    id VARCHAR(255) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    time_mili NUMERIC NOT NULL,
    amount NUMERIC NOT NULL,
    `description` TEXT,
    to_bid BOOLEAN DEFAULT FALSE,
    event_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- Migration for creating the 'donation_bids' table
CREATE TABLE donation_bids (
    donation_id VARCHAR(255) NOT NULL,
    bid_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (donation_id, bid_id),
    FOREIGN KEY (donation_id) REFERENCES donations(id) ON DELETE CASCADE,
    FOREIGN KEY (bid_id) REFERENCES bids(id) ON DELETE CASCADE
);

-- Migration for creating the 'teams_runs' table
CREATE TABLE teams_runs (
    run_id VARCHAR(255) NOT NULL,
    team_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (run_id, team_id),
    FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- Migration for creating the 'schedules' table
CREATE TABLE schedules (
    id VARCHAR(255) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    start_time_mili NUMERIC NOT NULL,
    end_time_mili NUMERIC NOT NULL,
    event_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- Migration for creating the 'events' table
CREATE TABLE events (
    id VARCHAR(255) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    start_time_mili NUMERIC NOT NULL,
    end_time_mili NUMERIC NOT NULL
);
