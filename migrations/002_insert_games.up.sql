insert into game (title, maker, description, max_players)
values
    ('Tekken 3', 'NAMCO', '', 2),
    ('Bobl Bubl', 'TAITO', '', 2),
    ('Cadilacs Dinosours', 'Capcom', '', 3);

insert into player (name, total_coins)
values
    ('eddy', 10);

insert into running_game (peer_name, orakki_id, orakki_state, game_id, first_player_id, joined_player_ids)
values 
    ('peer name 1', 'orakki1', 0, 1, 1, '');

insert into connection_info (orakki_id, player_id, state, server_data) 
values 
    ('orakki1', 1, 0, '');