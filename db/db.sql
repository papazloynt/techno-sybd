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

-- Change
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
--CREATE INDEX IF NOT EXISTS index_forum_id_hash ON "forum" USING HASH ("id");







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

-- change
CREATE INDEX IF NOT EXISTS index_thread_slug_hash ON "thread" USING HASH ("slug");

CREATE INDEX IF NOT EXISTS index_thread_id_hash ON "thread" USING HASH ("forum");
CREATE INDEX IF NOT EXISTS index_thread_forum_created ON "thread" ("forum", "created");









----------------------------------------------------------------- POST
CREATE UNLOGGED TABLE IF NOT EXISTS "post" (
    id   serial PRIMARY KEY,
    parent   int DEFAULT 0,
    path     int [] DEFAULT ARRAY [] :: INTEGER [],
    author   text NOT NULL,
    forum    text  NOT NULL,
    thread   int  NOT NULL,
    message  text NOT NULL,
    isEdited bool NOT NULL DEFAULT FALSE,
    created  timestamptz NOT NULL
);


--CREATE INDEX IF NOT EXISTS index_post_id ON "post" USING HASH ("id");
CREATE INDEX IF NOT EXISTS index_post_thread_id ON "post"("thread", "path");
-- new
CREATE INDEX IF NOT EXISTS post_path_complex ON "post" ((path[1]), path);


CREATE OR REPLACE FUNCTION increment_forum_thread() RETURNS TRIGGER AS $$
BEGIN
    UPDATE "forum" SET threads = threads + 1 where slug=NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE  plpgsql;

CREATE TRIGGER insert_thread
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE increment_forum_thread();


CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS $$
BEGIN
    new.path = (SELECT path FROM "post" WHERE id = new.parent) || new.id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_path
    BEFORE INSERT ON "post"
    FOR EACH ROW
    EXECUTE PROCEDURE update_path();


CREATE OR REPLACE FUNCTION increment_forum_posts() RETURNS TRIGGER AS $$
BEGIN
    UPDATE "forum" SET posts = "forum".posts + 1 WHERE slug = NEW.forum;
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

CREATE INDEX IF NOT EXISTS index_vote_exist ON "vote" ("thread", "nickname");
CREATE INDEX IF NOT EXISTS index_vote_update ON "vote" ("nickname", "thread", "voice");



CREATE OR REPLACE FUNCTION set_threads_votes() RETURNS TRIGGER AS $$
BEGIN
    UPDATE "thread" SET votes = votes + NEW.voice WHERE id = NEW.thread;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_votes
    AFTER INSERT
    ON "vote"
    FOR EACH ROW
EXECUTE PROCEDURE set_threads_votes();



CREATE OR REPLACE FUNCTION update_threads_votes() RETURNS TRIGGER AS $$
BEGIN
    UPDATE "thread" SET votes = votes + NEW.voice - OLD.voice WHERE id = NEW.thread;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_votes
    AFTER UPDATE
    ON "vote"
    FOR EACH ROW
EXECUTE PROCEDURE update_threads_votes();


----------------------------------------------------------------- FORUM-USER
CREATE UNLOGGED TABLE IF NOT EXISTS "forum_user" (
    nickname citext COLLATE "ucs_basic" NOT NULL,
    forum citext NOT NULL,
    fullname text NOT NULL,
    about text NOT NULL,
    email citext NOT NULL,
    CONSTRAINT forum_user_key UNIQUE (nickname, forum)
);

CREATE UNIQUE INDEX IF NOT EXISTS index_fast ON "forum_user"(forum, nickname);


CREATE OR REPLACE FUNCTION update_forum_user() RETURNS TRIGGER AS
$$
DECLARE
    nickname citext;
    fullname text;
    about    text;
    email    citext;
BEGIN
    SELECT u.nickname, u.fullname, u.about, u.email FROM "user" u WHERE u.nickname = NEW.author
    INTO nickname, fullname, about, email;

    INSERT INTO "forum_user" (nickname, fullname, about, email, forum)
    VALUES (nickname, fullname, about, email, NEW.forum) ON CONFLICT do nothing;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_forum_user_on_post
    AFTER INSERT
    ON "post"
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_user();

CREATE TRIGGER update_forum_user_on_thread
    AFTER INSERT
    ON "thread"
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_user();

