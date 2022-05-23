CREATE TABLE IF NOT EXISTS 'queries'
(
    'id'         TEXT    NOT NULL PRIMARY KEY,
    'query_text' TEXT    NOT NULL,
    'chat_type'  TEXT    NOT NULL,
    'user_id'    INTEGER NOT NULL,
    'user_name'  TEXT    NOT NULL,
    'user_lang'  TEXT    NOT NULL,
    'timestamp'  TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS 'messages'
(
    'id'        TEXT    NOT NULL PRIMARY KEY,
    'query_id'  TEXT    NOT NULL,
    'user_id'   INTEGER NOT NULL,
    'user_name' TEXT    NOT NULL,
    'timestamp' TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS 'members'
(
    'user_id'       INTEGER NOT NULL,
    'user_name'     TEXT    NOT NULL,
    'user_lang'     TEXT    NOT NULL,
    'from_id'       INTEGER NOT NULL,
    'message_id'    TEXT    NOT NULL,
    'chat_instance' TEXT    NOT NULL,
    'contest_id'    INTEGER NOT NULL,
    'timestamp'     TEXT    NOT NULL,
    UNIQUE (user_id, contest_id)
);

CREATE TABLE IF NOT EXISTS 'contests'
(
    'id'             INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    'creator_id'     INTEGER NOT NULL,
    'contest_name'   TEXT    NOT NULL,
    'contest_start'  TEXT    NOT NULL,
    'contest_end'    TEXT    NOT NULL,
    'contest_active' INTEGER NOT NULL,
    'username'       TEXT    NOT NULL,
    'timestamp'      TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS 'settings'
(
    'token'     TEXT    NOT NULL,
    'admin'     INTEGER NOT NULL,
    'is_public' INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS 'creators'
(
    'id'        INTEGER NOT NULL PRIMARY KEY,
    'user_name' TEXT    NOT NULL,
    'status'    INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS 'posts'
(
    'contest_id'  INTEGER NOT NULL,
    'type'        TEXT    NOT NULL,
    'title'       TEXT    NOT NULL,
    'message'     TEXT    NOT NULL,
    'description' TEXT    NOT NULL,
    'image'       TEXT    NOT NULL,
    UNIQUE ('contest_id', 'type')
);