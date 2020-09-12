create table game (
    id int not null auto_increment,
    title varchar(64) not null,
    maker varchar(64) not null,
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
    peer_name varchar(128) not null,
    orakki_id varchar(128) not null,
    orakki_state int not null,
    game_id int not null,
    first_player_id bigint not null,
    joined_player_ids varchar(128),
    created_at timestamp,
    
    primary key (id),
    unique (peer_name),
    foreign key (game_id) references game(id),
    foreign key (first_player_id) references player(id)
) character set utf8mb4 collate utf8mb4_unicode_ci;

create table signaling_info (
    id bigint not null auto_increment,
    orakki_id varchar(128) not null,
    data varchar(8196),
    created_at timestamp,
    is_last boolean not null,

    primary key (id)
) character set utf8mb4 collate utf8mb4_unicode_ci;
