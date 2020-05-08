INSERT INTO
    mid_events
    (admin_id,
     title,
     description,
     date,
     tags,
     photos,
     member_limit,
     is_public)
VALUES
    (
        5,
        'Концерт Арии',
        'Группа дает прощальный концерт, сыграет всем известные хиты - это будет просто вау',
        '2021-03-19 12:59:13+03',
        '{5,15}',
        '{mocks/aria-concert.jpg}',
        5,
        TRUE
    ),(
        6,
        'Киносеанс ‘Сталкера’ Тарковского',
        'Режиссерская версия, если вы понимаете о чем я.',
        '2021-02-18 11:55:23+03',
        '{1,2}',
        '{mocks/tarkovskiy-stalker.jpg}',
        15,
        TRUE
    ),(
        7,
        'Поход в 5 Д кинотеатр',
        'Ага, а может лучше 6 Д, а может лучше 13 Д?',
        '2021-01-17 10:50:33+03',
        '{2}',
        '{mocks/5d-kino.jpg}',
        3,
        FALSE
    ),(
        8,
        'Просмотр футбола в баре на Таганской',
        'Все будут орать и напиваться - ведь именно это нам и нужно.',
        '2021-03-16 12:49:43+03',
        '{1,11}',
        '{mocks/soccer.jpg,mocks/bar.jpg}',
        10,
        TRUE
    ),(
        9,
        'Зоопарк на выходных',
        'Сходим посмотрим на братьев наших меньших.',
        '2021-02-15 13:48:53+03',
        NULL,
        '{mocks/zoo.jpg,mocks/cat.jpg}',
        4,
        TRUE
    ),(
        10,
        'Вечер настольных игр у Андрея дома',
        'Там и мафия, и всякие маньчкины, ну что я рассказываю, приходи.',
        '2021-01-14 14:34:53+03',
        '{6,15}',
        NULL,
        6,
        TRUE
    ),(
        5,
        'Митинг за права женщин',
        'Созываю менеджерок, авторок и всех,кого исправляет Т9 на маке',
        '2021-03-13 15:24:13+03',
        NULL,
        '{mocks/meeting.jpg,mocks/apple-logo.jpg}',
        15,
        FALSE
    ),(
        6,
        'Поход в компьютерный клуб',
        'Зарубимся в контру и дотку, а может даже и квейк. Команды 5 на 5. Залетаем!',
        '2021-02-12 16:14:23+03',
        NULL,
        '{mocks/computer-room.jpg,mocks/pro-gamer.jpg}',
        10,
        TRUE
    ),(
        7,
        'Библио - Вечер на Китай-городе',
        'Придут такие именитые писатели как Дарья Донцова и Пелевин',
        '2021-01-11 17:47:24+03',
        '{14}',
        '{mocks/book-night.jpg}',
        6,
        TRUE
    ),(
        8,
        'Встреча любителей мотоспорта',
        'Врум - врум всем гонщикам, зажжем этот день',
        '2021-02-21 18:46:25+03',
        '{11}',
        '{mocks/moto.jpg}',
        10,
        TRUE
    ),(
        9,
        'Собрание клуба молодых мамочек',
        'Мой годовасик-тугосеря уже делает агу-агу, давайте же это обсудим',
        '2021-01-22 19:44:26+03',
        NULL,
        '{mocks/baby.jpg}',
        10,
        TRUE
    ),(
        10,
        'Собрание клуба любителей пленочной фотографии',
        'Это искусство, ребят, тут и говорить нечего',
        '2021-03-23 11:45:27+03',
        '{6,15}',
        '{mocks/photographer.jpg,mocks/zenit.jpg}',
        12,
        TRUE
    ),(
        5,
        'Открытие БДСМ-клуба',
        'Отлупим плеткой друг-дружку, вставим пару кляпов и пробок куда нужно',
        '2021-02-24 12:43:28+03',
        '{13}',
        '{mocks/50-shades-of-gray.jpg}',
        9,
        FALSE
    ),(
        6,
        'Постреляем с воздушки на день независимости',
        'Странное занятие, но имеет исторические корни, между прочим',
        '2021-01-25 13:42:29+03',
        '{15}',
        NULL,
        5,
        TRUE
    ),(
        7,
        'Охота на кабанов в Подмосковье',
        'Отец даст Сайгу-12 К каждому жаждущему и оторвемся',
        '2021-03-05 14:41:20+03',
        '{6}',
        '{mocks/kaban.jpg}',
        7,
        TRUE
    ),(
        8,
        'Собрание тех, чье имя никто не любит',
        'Также жду Антонов',
        '2021-02-04 15:54:13+03',
        NULL,
        '{mocks/lapenko.jpg}',
        7,
        TRUE
    ),(
        9,
        'Собрание тех, чье имя любят все',
        'Мяу - Мяу, мои котятки, всех жду, записываемся на ноготочки',
        '2021-01-03 16:14:33+03',
        NULL,
        '{mocks/putin.jpg}',
        7,
        TRUE
    ),(
        10,
        'Матч Аргентина-Ямайка',
        'Все знают, что будет на этом матче',
        '2021-02-02 17:24:43+03',
        '{1}',
        '{mocks/soccer.jpg,mocks/pain.jpg}',
        15,
        TRUE
    ),(
        5,
        'Танцевальный мастер-класс в Парке Горького',
        'Событие для молодых пар и тех, кто хочет научиться красиво танцевать',
        '2021-03-01 18:34:53+03',
        '{15}',
        '{mocks/dance.jpg,mocks/gorky-park.jpg}',
        12,
        TRUE
    );

UPDATE
    mid_events
SET
    title_tsv = setweight(to_tsvector(title), 'A') || setweight(to_tsvector(description), 'B');

-- WITH chat_meta AS
--     (
--         SELECT uc.title, SUM(m.is_shown = TRUE) AS unseen
--         FROM user_chat uc JOIN messages m ON m.user_local_id = uc.user_local_id ORDER BY m.created LIMIT 1
--     ) SELECT
--
--
-- SELECT uc.user_local_id, uc.title, SUM(m.is_shown = TRUE) AS unseen, MAX(m.created) AS last_date,  m.message
-- FROM user_chat uc JOIN messages m ON m.user_local_id = uc.user_local_id WHERE uc.uid = $1 GROUP BY uc.user_local_id;

-- INSERT INTO chat_user (chat_id, admin_id, user_count, title) VALUES (1, 1, 2, 'title');
-- INSERT INTO user_chat (chat_local_id, uid) VALUES (1, 1), (1, 2);
-- UPDATE user_vote SET chat_id = 1 WHERE (uid = 1 AND user_id = 2) OR (uid = 2 AND user_id = 1);
-- INSERT INTO message (uid, chat_id, user_local_id, message, is_shown) VALUES (1, 1, 1, 'Hello world!', FALSE) RETURNING mid;
-- INSERT INTO message (uid, chat_id, user_local_id, message, is_shown) VALUES (1, 1, 2, 'Hello world!', FALSE) RETURNING mid;
