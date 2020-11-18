INSERT INTO pack (status, title, maker, description, max_players, poster_url, rom_name)
VALUES
    (1, 'Tekken 3', 'NAMCO', 'Tekken 3 (鉄拳3) is a fighting game, the third installment in the Tekken series.', 2, 'https://oraksil.s3.ap-northeast-2.amazonaws.com/packs/tekken3/tekken3_poster.jpg', 'tekken3'),
    (1, 'Bobl Bubl', 'TAITO', '', 2, 'https://oraksil.s3.ap-northeast-2.amazonaws.com/packs/bublbobl/bublbobl_poster.jpg', 'bublbobl'),
    (1, 'Cadilacs Dinosours', 'Capcom', "It is a side-scrolling beat 'em up based on the comic book series Xenozoic Tales. The game was produced as a tie-in to the short-lived Cadillacs and Dinosaurs animated series which was aired during the same year the game was released.", 3, 'https://oraksil.s3.ap-northeast-2.amazonaws.com/packs/dino/dino_poster.jpg', 'dino'),
    (0, 'Final Fight II', 'Capcom', '', 3, 'https://oraksil.s3.ap-northeast-2.amazonaws.com/packs/ffight2b/ffight2b_poster.png', 'ffight2b'),
    (1, 'Super Tank', 'Video Games GmbH', '', 2, 'https://oraksil.s3.ap-northeast-2.amazonaws.com/packs/supertnk/supertnk_poster.png', 'supertnk');

INSERT INTO player (name, total_coins)
VALUES
    ('gamz', 10);
