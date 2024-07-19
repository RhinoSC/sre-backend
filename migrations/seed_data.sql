-- Insertar datos de prueba en la tabla users
INSERT INTO users (id, name, username, created_at) VALUES
('user1', 'John Doe', 'johndoe', datetime('now')),
('user2', 'Jane Smith', 'janesmith', datetime('now')),
('user3', 'Alice Johnson', 'alicejohnson', datetime('now')),
('user4', 'Bob Brown', 'bobbrown', datetime('now'));

-- Insertar datos de prueba en la tabla user_socials
INSERT INTO user_socials (id, user_id, twitch, twitter, youtube, facebook) VALUES
('social1', 'user1', 'twitch_user1', 'twitter_user1', 'youtube_user1', 'facebook_user1'),
('social2', 'user2', 'twitch_user2', 'twitter_user2', 'youtube_user2', 'facebook_user2'),
('social3', 'user3', 'twitch_user3', 'twitter_user3', 'youtube_user3', 'facebook_user3'),
('social4', 'user4', 'twitch_user4', 'twitter_user4', 'youtube_user4', 'facebook_user4');

-- Insertar datos de prueba en la tabla teams
INSERT INTO teams (id, name) VALUES
('team1', 'Team Alpha'),
('team2', 'Team Beta'),
('team3', 'Team Gamma'),
('team4', 'Team Delta');

-- Insertar datos de prueba en la tabla runs
INSERT INTO runs (id, name, start_time_mili, estimate_string, estimate_milliseconds, metadata, schedule_id) VALUES
('run1', 'Speedrun 1', 1000, '2h', 7200000, 'Metadata for Run 1', 'schedule1'),
('run2', 'Speedrun 2', 2000, '1h 30m', 5400000, 'Metadata for Run 2', 'schedule2');

-- Insertar datos de prueba en la tabla run_metadata
INSERT INTO run_metadata (id, run_id, category, platform, twitch_game_name, run_name, note) VALUES
('metadata1', 'run1', 'Any%', 'PC', 'Game 1', 'Run 1', 'Note for Run 1'),
('metadata2', 'run2', '100%', 'Console', 'Game 2', 'Run 2', 'Note for Run 2');

-- Insertar datos de prueba en la tabla players (Para un run con un equipo y varios jugadores)
INSERT INTO players (team_id, user_id) VALUES
('team1', 'user1'),
('team1', 'user2');

-- Insertar datos de prueba en la tabla players (Para un run con varios equipos, cada uno con un jugador)
INSERT INTO players (team_id, user_id) VALUES
('team2', 'user3'),
('team3', 'user4');

-- Insertar datos de prueba en la tabla bids
INSERT INTO bids (id, bidname, goal, current_amount, description, type, create_new_options, run_id) VALUES
('bid1', 'Bid 1', 100, 50, 'Description for Bid 1', 'bidwar', true, 'run1'),
('bid2', 'Bid 2', 200, 100, 'Description for Bid 2', 'total', false, 'run2');

-- Insertar datos de prueba en la tabla bid_options
INSERT INTO bid_options (id, bid_id, name, current_amount) VALUES
('option1', 'bid1', 'Option 1', 30),
('option2', 'bid1', 'Option 2', 20),
('option3', 'bid2', 'Option 3', 70);

-- Insertar datos de prueba en la tabla donations
INSERT INTO donations (id, name, email, time_mili, amount, description, to_bid, event_id) VALUES
('donation1', 'Alice', 'alice@example.com', 1609459200000, 50, 'Donation for Run 1', true, 'event1'),
('donation2', 'Bob', 'bob@example.com', 1609545600000, 100, 'Donation for Run 2', false, 'event2');

-- Insertar datos de prueba en la tabla donation_bids
INSERT INTO donation_bids (donation_id, bid_id) VALUES
('donation1', 'bid1'),
('donation2', 'bid2');

-- Insertar datos de prueba en la tabla teams_runs
INSERT INTO teams_runs (run_id, team_id) VALUES
('run1', 'team1'),
('run2', 'team2'),
('run2', 'team3'),
('run2', 'team4');

-- Insertar datos de prueba en la tabla schedules
INSERT INTO schedules (id, name, start_time_mili, end_time_mili, event_id) VALUES
('schedule1', 'Schedule 1', 1609459200000, 1609545600000, 'event1'),
('schedule2', 'Schedule 2', 1609632000000, 1609718400000, 'event2');

-- Insertar datos de prueba en la tabla events
INSERT INTO events (id, name, start_time_mili, end_time_mili) VALUES
('event1', 'Event 1', 1609459200000, 1609545600000),
('event2', 'Event 2', 1609632000000, 1609718400000);

-- Insertar datos de prueba en la tabla prizes
INSERT INTO prizes (id, name, description, url, min_amount, status, international_delivery, event_id) VALUES
('prize1', 'Prize 1', 'Description for Prize 1', 'http://example.com/prize1', 50, 'Available', true, 'event1'),
('prize2', 'Prize 2', 'Description for Prize 2', 'http://example.com/prize2', 100, 'Available', false, 'event2');
