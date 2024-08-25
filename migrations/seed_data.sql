-- Insertar datos en la tabla 'events'
INSERT INTO events (id, name, start_time_mili, end_time_mili) VALUES
('event1', 'Charity Marathon', 1677830400000, 1677916800000);

-- Insertar datos en la tabla 'schedules'
INSERT INTO schedules (id, name, start_time_mili, end_time_mili, setup_time_mili, event_id) VALUES
('schedule1', 'Morning Session', 1677830400000, 1677834000000, 300000, 'event1');

-- Insertar datos de prueba en la tabla runs
INSERT INTO runs (id, name, start_time_mili, estimate_string, estimate_mili, setup_time_mili, status, schedule_id) VALUES
('run1', 'Speedrun 1', 0, '2h', 7200000, 300000, 'default', 'schedule1'),
('run2', 'Speedrun 2', 0, '1h 30m', 5400000, 300000, 'default', 'schedule1'),
('run4', 'Speedrun 12', 0, '3h', 10800000, 300000, 'active', 'schedule1'),
('run5', 'Speedrun 13', 0, '4h 30m', 16200000, 300000, 'default', 'schedule1'),
('run6', 'Speedrun 14', 0, '5h', 18000000, 300000, 'default', 'schedule1');

-- Insertar datos de prueba en la tabla run_metadata
INSERT INTO run_metadata (id, run_id, category, platform, twitch_game_name, twitch_game_id, run_name, note) VALUES
('metadata1', 'run1', 'Any%', 'PC', 'Game 1', 33214, 'Run 1', 'Note for Run 1'),
('metadata2', 'run2', '100%', 'Console', 'Game 2', 33214, 'Run 2', 'Note for Run 2'),
('metadata4', 'run4', 'Any%', 'PC', 'Super Mario Bros', 509508, 'Run 12', 'Note for Run 12'),
('metadata5', 'run5', '100%', 'Console', 'Zelda', 'Run 13', 33214, 'Note for Run 13'),
('metadata6', 'run6', 'All%', 'PC', 'Minecraft', 'Run 14', 33214, 'Note for Run 14');

-- Insertar datos en la tabla 'teams'
INSERT INTO teams (id, name, run_id) VALUES
('team1', 'Team Alpha', 'run1'),
('team2', 'Team Beta', 'run1'),
('team3', 'Team3', 'run2'),
('team4', 'Team4', 'run3');


-- Insertar datos de prueba en la tabla users
INSERT INTO users (id, name, username, created_at) VALUES
('user1', 'John Doe', 'johndoe', datetime('now')),
('user2', 'Jane Smith', 'janesmith', datetime('now')),
('user3', 'Alice Johnson', 'alicejohnson', datetime('now')),
('user4', 'Bob Brown', 'bobbrown', datetime('now')),
('user5', 'Charlie White', 'charliewhite', datetime('now')),
('user6', 'Diana Green', 'dianagreen', datetime('now')),
('user7', 'Eve Black', 'eveblack', datetime('now')),
('user8', 'Frank Blue', 'frankblue', datetime('now'));

-- Insertar datos de prueba en la tabla user_socials
INSERT INTO user_socials (id, user_id, twitch, twitter, youtube, facebook) VALUES
('social1', 'user1', 'twitch_user1', 'twitter_user1', 'youtube_user1', 'facebook_user1'),
('social2', 'user2', 'twitch_user2', 'twitter_user2', 'youtube_user2', 'facebook_user2'),
('social3', 'user3', 'twitch_user3', 'twitter_user3', 'youtube_user3', 'facebook_user3'),
('social4', 'user4', 'twitch_user4', 'twitter_user4', 'youtube_user4', 'facebook_user4'),
('social5', 'user5', 'twitch_user5', 'twitter_user5', 'youtube_user5', 'facebook_user5'),
('social6', 'user6', 'twitch_user6', 'twitter_user6', 'youtube_user6', 'facebook_user6'),
('social7', 'user7', 'twitch_user7', 'twitter_user7', 'youtube_user7', 'facebook_user7'),
('social8', 'user8', 'twitch_user8', 'twitter_user8', 'youtube_user8', 'facebook_user8');

-- Insertar datos de prueba en la tabla admins
INSERT INTO admins (id, username, password, created_at) VALUES
('admin1', 'admin', '$2a$10$u8fl6qTAR.Fd/ZmGZqJYm.dihtVBcXw/WCuXB/ZkmP8gvqDgqSMmG', datetime('now'));

-- Insertar datos de prueba en la tabla prizes
INSERT INTO prizes (id, name, description, url, min_amount, status, international_delivery, event_id) VALUES
('prize1', 'Prize 1', 'Description for Prize 1', 'http://example.com/prize1', 50, 'Available', true, 'event1'),
('prize2', 'Prize 2', 'Description for Prize 2', 'http://example.com/prize2', 100, 'Available', false, 'event2'),
('prize3', 'Prize 3', 'Description for Prize 3', 'http://example.com/prize3', 150, 'Available', true, 'event3'),
('prize4', 'Prize 4', 'Description for Prize 4', 'http://example.com/prize4', 200, 'Available', false, 'event4');

-- Insertar datos en la tabla bids
INSERT INTO bids (id, bidname, goal, current_amount, description, type, create_new_options, status, run_id) VALUES
('bid1', 'Bid War 1', 1000, 0, 'First Bid War', 'bidwar', true, 'active', 'run1'),
('bid2', 'Total Donation 1', 2000, 0, 'Total Donation Goal', 'total', false, 'active', 'run2'),
('bid3', 'Goal 1', 500, 0, 'Achieve this Goal', 'goal', false, 'default', 'run4'),
('bid4', 'Bid War 2', 1500, 0, 'Second Bid War', 'bidwar', true, 'active', 'run5'),
('bid5', 'Total Donation 2', 2500, 0, 'Another Total Donation Goal', 'total', false, 'default', 'run6');

-- Insertar datos de prueba en la tabla bid_options
INSERT INTO bid_options (id, bid_id, name, current_amount) VALUES
('option1', 'bid1', 'Option A', 0),
('option2', 'bid1', 'Option B', 0),
('option3', 'bid4', 'Option X', 0),
('option4', 'bid4', 'Option Y', 0);

-- Insertar datos de prueba en la tabla donations
INSERT INTO donations (id, name, email, time_mili, amount, description, to_bid, event_id) VALUES
('donation1', 'Alice', 'alice@example.com', 1609459200000, 20, 'Donation for Run 1', true, 'event1'),
('donation2', 'Bob', 'bob@example.com', 1609545600000, 70, 'Donation for Run 2', true, 'event2'),
('donation3', 'Pep', 'pep@example.com', 1609545600000, 20, 'Donation for Run 1 But Better', true, 'event1'),
('donation4', 'Mary', 'mary@example.com', 1609632000000, 20, 'Another Donation for Bid1', true, 'event1'),
('donation5', 'Rhino', 'rhino@example.com', 1609632000000, 80, 'Another Donation for Bid3', true, 'event1'),
('donation6', 'Chemi', 'chemi@example.com', 1609632000000, 30, 'No bid donation!', false, 'event1'),
('donation7', 'George', 'george@example.com', 1609724800000, 50, 'Donation for Run 12', true, 'event1'),
('donation8', 'Hannah', 'hannah@example.com', 1609811200000, 75, 'Donation for Run 13', true, 'event2'),
('donation9', 'Ian', 'ian@example.com', 1609897600000, 100, 'Donation for Run 14', true, 'event2');

-- Insertar datos de prueba en la tabla donation_bids
INSERT INTO donation_bids (donation_id, bid_id, bid_option_id) VALUES
('donation1', 'bid1', 'option1'),
('donation2', 'bid2', 'option3'),
('donation3', 'bid1', 'option2'),
('donation4', 'bid1', 'option4'),
('donation5', 'bid3', NULL),
('donation7', 'bid5', 'option6'),
('donation8', 'bid6', 'option7'),
('donation9', 'bid7', 'option8');

-- Insertar datos de prueba en la tabla players
INSERT INTO players (team_id, user_id) VALUES
('team1', 'user1'),
('team1', 'user2'),
('team2', 'user3'),
('team3', 'user4'),
('team7', 'user5'),
('team8', 'user6'),
('team8', 'user7'),
('team9', 'user8'),
('team9', 'user1');

-- Asegurar la consistencia de los datos

-- Verificar que todas las runs tienen al menos un equipo asignado
DELETE FROM runs WHERE id NOT IN (SELECT run_id FROM teams);

-- Actualizar current_amount en bids basado en las donaciones asignadas
UPDATE bids
SET current_amount = (
    SELECT COALESCE(SUM(d.amount), 0)
    FROM donations d
    JOIN donation_bids db ON d.id = db.donation_id
    WHERE db.bid_id = bids.id
);

-- Actualizar current_amount en bid_options basado en las donaciones asignadas
UPDATE bid_options
SET current_amount = (
    SELECT COALESCE(SUM(d.amount), 0)
    FROM donations d
    JOIN donation_bids db ON d.id = db.donation_id
    WHERE db.bid_option_id = bid_options.id
);