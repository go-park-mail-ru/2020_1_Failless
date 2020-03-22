-- CREATE DATABASE eventum WITH OWNER = postgres;
-- CREATE ROLE eventum WITH SUPERUSER PASSWORD 'eventum' LOGIN CONNECTION LIMIT -1;

-- CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

CREATE TYPE SEX_T AS ENUM ('male', 'female', 'other');
CREATE TYPE ETYPE_T AS ENUM ('concert', 'museum', 'bar', 'theater', 'walk', 'tour');

CREATE TEXT SEARCH DICTIONARY russian_ispell
(
    TEMPLATE = ispell,
    DictFile = russian,
    AffFile = russian,
    StopWords = russian
);

CREATE TEXT SEARCH CONFIGURATION ru (COPY =russian);

ALTER TEXT SEARCH CONFIGURATION ru
    ALTER MAPPING FOR hword, hword_part, word
        WITH russian_ispell, russian_stem;

SET default_text_search_config = 'ru';

CREATE TABLE IF NOT EXISTS profile
(
    uid      SERIAL PRIMARY KEY,
    name     VARCHAR(32)        NOT NULL check ( name <> '' ),
    phone    VARCHAR(12) UNIQUE NOT NULL CHECK ( phone <> '' ),
    email    VARCHAR(64) UNIQUE NOT NULL CHECK ( email <> '' ),
    password BYTEA              NOT NULL CHECK ( octet_length(password) <> 0 )
);

CREATE TABLE IF NOT EXISTS tag
(
    tag_id SERIAL PRIMARY KEY,
    name   VARCHAR(32) UNIQUE NOT NULL CHECK ( name <> '' )
);

CREATE TABLE IF NOT EXISTS user_tag
(
    tag_id INTEGER REFERENCES tag (tag_id),
    uid    INTEGER REFERENCES profile (uid)
);

CREATE TABLE IF NOT EXISTS profile_info
(
    pid        INTEGER PRIMARY KEY REFERENCES profile (uid),
    about      VARCHAR(512)                         DEFAULT '',
    photos     VARCHAR(64)[]                        DEFAULT NULL,
    rating     FLOAT                                DEFAULT 0,
    location   GEOGRAPHY                            DEFAULT NULL, -- ST_POINT(latitude, longitude)
    birthday   DATE                                 DEFAULT NULL,
    gender     SEX_T,
    login_date TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS timetable
(
    table_id  SERIAL PRIMARY KEY,
    title     VARCHAR(128)                NOT NULL CHECK ( title <> '' ),
    edate     TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    message   VARCHAR(1024)               NOT NULL,
    is_edited BOOLEAN                              DEFAULT FALSE,
    sponsor   CITEXT                               DEFAULT NULL,
    photos    VARCHAR(64)[]                        DEFAULT NULL,
    etype     ETYPE_T
);

CREATE TABLE IF NOT EXISTS events
(
    eid       SERIAL PRIMARY KEY,
    uid       INTEGER REFERENCES profile_info (pid),
    title     VARCHAR(128)                NOT NULL CHECK ( title <> '' ),
    edate     TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    message   VARCHAR(1024)               NOT NULL,
    is_edited BOOLEAN                              DEFAULT FALSE,
    is_public BOOLEAN                              DEFAULT FALSE,
    author    CITEXT                               DEFAULT NULL,
    etype     INTEGER REFERENCES tag (tag_id),
    range     SMALLINT                             DEFAULT 1,
    title_tsv TSVECTOR,
    photos    VARCHAR(64)[]                        DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS event_vote
(
    uid   INTEGER  NOT NULL REFERENCES profile (uid),
    eid   INTEGER  NOT NULL REFERENCES events (eid),
    value SMALLINT NOT NULL DEFAULT 0,
    bid   INTEGER           DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS subscribe
(
    uid      INTEGER  NOT NULL REFERENCES profile (uid),
    table_id INTEGER  NOT NULL REFERENCES timetable (table_id),
    value    SMALLINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_vote
(
    author_id INTEGER  NOT NULL REFERENCES profile (uid),
    eid       INTEGER  NOT NULL REFERENCES events (eid),
    user_id   INTEGER  NOT NULL REFERENCES profile (uid),
    is_edited BOOLEAN  NOT NULL DEFAULT FALSE,
    value     SMALLINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS band
(
    bid SERIAL PRIMARY KEY,
    eid INTEGER REFERENCES events (eid)
);

CREATE TABLE IF NOT EXISTS message
(
    mid     BIGSERIAL PRIMARY KEY,
    uid     INTEGER REFERENCES profile (uid),
    bid     INTEGER REFERENCES band (bid),
    message TEXT,
    created TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);

CREATE OR REPLACE PROCEDURE add_location(uid INT, latitude FLOAT, longitude FLOAT)
    LANGUAGE plpgsql AS
$$
BEGIN
    UPDATE profile_info
    SET location = ST_POINT(latitude, longitude)
    WHERE pid = uid;
    COMMIT;
END;
$$;

CREATE INDEX IF NOT EXISTS event_title_idx ON events USING GIN (title);

CREATE OR REPLACE FUNCTION make_tsvector()
    RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE events
        SET title_tsv = setweight(to_tsvector('russian', title), 'A')
            || setweight(to_tsvector('russian', message), 'B')
        WHERE eid = NEW.eid;
        RETURN NULL;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tg_add_event
    AFTER INSERT
    ON events
    FOR EACH ROW
EXECUTE FUNCTION make_tsvector();


GRANT ALL PRIVILEGES ON DATABASE eventum TO eventum;
GRANT USAGE ON SCHEMA public TO eventum;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO eventum;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO eventum;


INSERT INTO tag (name)
VALUES ('хочувБАР'),
       ('хочувКИНО'),
       ('хочувТЕАТР'),
       ('хочувКЛУБ'),
       ('хочунаКОНЦЕРТ'),
       ('хочуГУЛЯТЬ'),
       ('хочунаКАТОК'),
       ('хочунаВЫСТАВКУ'),
       ('хочуСПАТЬ'),
       ('хочунаСАЛЮТ'),
       ('хочувСПОРТ'),
       ('хочувМУЗЕЙ'),
       ('хочунаЛЕКЦИЮ'),
       ('хочуБОТАТЬ'),
       ('хочувПАРК')
;
