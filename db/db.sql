CREATE EXTENSION IF NOT EXISTS citext;

----------------------------------------------------------------- USER
CREATE UNLOGGED TABLE IF NOT EXISTS "user" (
  id serial,
  nickname citext COLLATE "ucs_basic" NOT NULL PRIMARY KEY,
  fullname text NOT NULL,
  about text NOT NULL DEFAULT '',
  email citext NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS index_user_nickname_hash ON "user" USING HASH ("nickname");
CREATE INDEX IF NOT EXISTS index_user_email_hash ON "user" (nickname, email);


----------------------------------------------------------------- FORUM
CREATE UNLOGGED TABLE IF NOT EXISTS "forum" (
    id      serial,
    slug    citext NOT NULL PRIMARY KEY,
    title   text  NOT NULL,
    "user"  citext NOT NULL,
    posts   int NOT NULL DEFAULT 0,
    threads int NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS index_forum_slug_hash ON "forum" USING HASH ("slug");

----------------------------------------------------------------- FORUM-USER
CREATE UNLOGGED TABLE IF NOT EXISTS "forum_user" (
     nickname citext NOT NULL REFERENCES "user" (nickname),
     forum  citext NOT NULL REFERENCES "forum" (slug),
     PRIMARY KEY (forum, nickname)
);

CREATE FUNCTION add_forum_user() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO forum_user (forum, nickname)
     VALUES (NEW.forum, NEW.author)
     ON conflict do nothing;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER forum_user
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();

CREATE TRIGGER forum_user
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();

----------------------------------------------------------------- THREAD
CREATE UNLOGGED TABLE IF NOT EXISTS "thread" (
    id     serial PRIMARY KEY,
    slug    citext NOT NULL,
    title   text NOT NULL,
    author  text NOT NULL,
    forum   text NOT NULL,
    message text NOT NULL,
    votes  int NOT NULL DEFAULT 0,
    created timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS index_thread_slug_hash ON "thread" USING HASH ("slug");
CREATE INDEX IF NOT EXISTS index_thread_id_hash ON "thread" USING HASH ("forum");
CREATE INDEX IF NOT EXISTS index_thread_forum_created ON "thread" ("forum", "created");


----------------------------------------------------------------- POST
CREATE UNLOGGED TABLE IF NOT EXISTS "post" (
    id   serial PRIMARY KEY,
    parent   int DEFAULT 0,
    path     int [] NOT NULL,
    author   text NOT NULL,
    forum    text  NOT NULL,
    thread   int  NOT NULL,
    message  text NOT NULL,
    isEdited bool NOT NULL DEFAULT FALSE,
    created  timestamptz NOT NULL
);


CREATE INDEX IF NOT EXISTS index_post_thread_id ON "post"("thread", "path");
CREATE INDEX IF NOT EXISTS index_post_path_complex ON "post" ((path[1]), path);


CREATE OR REPLACE FUNCTION increment_thread() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE "forum"
    SET threads = threads + 1
    WHERE slug=NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE  plpgsql;

CREATE TRIGGER insert_thread
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE increment_thread();


CREATE OR REPLACE FUNCTION new_path() RETURNS TRIGGER AS $$
BEGIN
    new.path = (SELECT path
                FROM "post"
                WHERE id = new.parent) || new.id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_path
    BEFORE INSERT ON "post"
    FOR EACH ROW
    EXECUTE PROCEDURE new_path();


CREATE OR REPLACE FUNCTION increment_forum_posts() RETURNS TRIGGER AS $$
BEGIN
    UPDATE "forum"
    SET posts = "forum".posts + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_count_posts
    AFTER INSERT
    ON "post"
    FOR EACH ROW
EXECUTE PROCEDURE increment_forum_posts();

----------------------------------------------------------------- VOTE
CREATE UNLOGGED TABLE IF NOT EXISTS "vote" (
  nickname text NOT NULL,
  thread int NOT NULL,
  voice int NOT NULL
);

CREATE INDEX IF NOT EXISTS index_vote_update ON "vote" ("nickname", "thread", "voice");


CREATE OR REPLACE FUNCTION make_votes() RETURNS TRIGGER AS $$
BEGIN
    UPDATE "thread"
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_votes
    AFTER INSERT
    ON "vote"
    FOR EACH ROW
EXECUTE PROCEDURE make_votes();


CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS $$
BEGIN
    UPDATE "thread"
    SET votes = votes + NEW.voice - OLD.voice
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_votes
    AFTER UPDATE
    ON "vote"
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();

VACUUM ANALYZE;