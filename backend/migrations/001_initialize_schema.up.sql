create table game (
    id int not null auto_increment,
    title varchar(64) not null,
    description text,
    max_players int, 

    primary key (id)
) character set utf8mb4 collate utf8mb4_unicode_ci;

create table player (
    id bigint not null auto_increment,
    name varchar(64) not null,
    total_coins int,

    primary key (id)
) character set utf8mb4 collate utf8mb4_unicode_ci;

create table running_game (
    id bigint not null auto_increment,
    game_id int not null,
    owner_player_id bigint not null,
    joined_player_ids varchar(128),
    created_at timestamp,
    
    primary key (id),
    foreign key (game_id) references game(id),
    foreign key (owner_player_id) references player(id)
) character set utf8mb4 collate utf8mb4_unicode_ci;
