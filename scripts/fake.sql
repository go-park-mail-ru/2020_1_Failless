INSERT INTO events (uid, title, edate, message, is_edited, author, etype)
VALUES (2, 'Концерт Арии', '2021-03-19 12:59:13+03', 'Группа дает прощальный концерт,
        сыграет всем известные хиты - это будет просто вау', false, 'Алексей', 5)
        ,
       (3, 'Киносеанс ‘Сталкера’ Тарковского', '2021-02-18 11:55:23+03', 'Режиссерская версия,
        если вы понимаете о чем я.', false, 'Серж', 2)
        ,
       (4, 'Поход в 5 Д кинотеатр', '2021-01-17 10:50:33+03', 'Ага, а может лучше 6 Д, а может лучше 13 Д?', false, 'Сержик', 2)
        ,
       (5, 'Просмотр футбола в баре на Таганской', '2021-03-16 12:49:43+03', 'Все будут орать и напиваться - ведь именно это нам и нужно.', false, 'Егорыч', 1)
        ,
       (6, 'Зоопарк на выходных', '2021-02-15 13:48:53+03', 'Сходим посмотрим на братьев наших меньших.', false, 'Димыч', 6)
        ,
       (7, 'Вечер настольных игр у Андрея дома', '2021-01-14 14:34:53+03', 'Там и мафия, и всякие маньчкины,
        ну что я рассказываю, приходи.', false, 'Андрюха', 2)
        ,
       (8, 'Митинг за права женщин', '2021-03-13 15:24:13+03', 'Созываю менеджерок, авторок и всех,
        кого исправляет Т9 на маке', false, 'Катрина', 3)
        ,
       (9, 'Поход в компьютерный клуб', '2021-02-12 16:14:23+03', 'Зарубимся в контру и дотку, а может даже и квейк', false, 'Кирилл', 6)
        ,
       (1, 'Библио - Вечер на Китай-городе', '2021-01-11 17:47:24+03', 'Придут такие именитые писатели как Дарья Донцова и Пелевин', false, 'Миша', 4)
        ,
       (2, 'Встреча любителей мотоспорта', '2021-02-21 18:46:25+03', 'Врум - врум всем гонщикам, зажжем этот день',
        false, 'Женька', 5)
        ,
       (3, 'Собрание клуба молодых мамочек', '2021-01-22 19:44:26+03', 'Мой годовасик-тугосеря уже делает агу-агу, давайте же это обсудим', false, 'Юля', 6)
        ,
       (4, 'Собрание клуба любителей пленочной фотографии', '2021-03-23 11:45:27+03', 'Это искусство, ребят,
        тут и говорить нечего', false, 'Марго', 7)
        ,
       (5, 'Открытие БДСМ-клуба', '2021-02-24 12:43:28+03', 'Отлупим плеткой друг-дружку,
        вставим пару кляпов и пробок куда нужно', false, 'Агата', 11)
        ,
       (6, 'Постреляем с воздушки на день независимости', '2021-01-25 13:42:29+03', 'Странное занятие,
        но имеет исторические корни, между прочим', false, 'Мия', 11)
        ,
       (7, 'Охота на кабанов в Подмосковье', '2021-03-05 14:41:20+03', 'Отец даст Сайгу-12 К каждому жаждущему и оторвемся', false, 'Игорь', 12)
        ,
       (8, 'Собрание тех, чье имя никто не любит', '2021-02-04 15:54:13+03', 'Также жду Антонов', false, 'Макс', 2)
        ,
       (9, 'Собрание тех, чье имя любят все', '2021-01-03 16:14:33+03', 'Мяу - Мяу, мои котятки, всех жду,
        записываемся на ноготочки', false, 'Анечка', 2)
        ,
       (1, 'Матч Аргентина-Ямайка', '2021-02-02 17:24:43+03', 'Все знают, что будет на этом матче', false, 'Даша', 11)
        ,
       (2, 'Танцевальный мастер-класс в Парке Горького', '2021-03-01 18:34:53+03', 'Событие для молодых пар и тех,
        кто хочет научиться красиво танцевать', false, 'Марина', 15);
UPDATE events SET title_tsv = setweight(to_tsvector(title), 'A') || setweight(to_tsvector(message), 'B');

-- WITH chat_meta AS
--     (
--         SELECT uc.title, SUM(m.is_shown = TRUE) AS unseen
--         FROM user_chat uc JOIN messages m ON m.user_local_id = uc.user_local_id ORDER BY m.created LIMIT 1
--     ) SELECT
--
--
-- SELECT uc.user_local_id, uc.title, SUM(m.is_shown = TRUE) AS unseen, MAX(m.created) AS last_date,  m.message
-- FROM user_chat uc JOIN messages m ON m.user_local_id = uc.user_local_id WHERE uc.uid = $1 GROUP BY uc.user_local_id;

INSERT INTO chat_user (chat_id, admin_id, user_count, title) VALUES (1, 1, 2, 'title');
INSERT INTO user_chat (chat_local_id, uid) VALUES (1, 1), (1, 2);
UPDATE user_vote SET chat_id = 1 WHERE (uid = 1 AND user_id = 2) OR (uid = 2 AND user_id = 1);
INSERT INTO message (uid, chat_id, user_local_id, message, is_shown) VALUES (1, 1, 1, 'Hello world!', FALSE) RETURNING mid;
INSERT INTO message (uid, chat_id, user_local_id, message, is_shown) VALUES (1, 1, 2, 'Hello world!', FALSE) RETURNING mid;
