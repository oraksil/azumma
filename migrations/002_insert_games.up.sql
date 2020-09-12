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
    ('peer1', 'orakki1', 0, 1, 1, '');

insert into signaling_info (orakki_id, data, is_last)
values
    ('orakki1', 'Ice Candidate 1', false);


insert into signaling_info (orakki_id, data, is_last)
values
    ('orakki1', 'Ice Candidate 2', false);

insert into signaling_info (orakki_id, data, is_last)
values
    ('orakki1', 'Ice Candidate 3', true);
