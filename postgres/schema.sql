CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
    user_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
    name varchar(80) NOT NULL,
    email varchar(80) NOT NULL,
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS followers (
    followee_user_id uuid NOT NULL references users(user_id),
    follower_user_id uuid NOT NULL references users(user_id),
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT followers_pk PRIMARY KEY (followee_user_id, follower_user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
    sender_user_id uuid NOT NULL references users(user_id),
    content varchar(200) NOT NULL,
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT messages_pk PRIMARY KEY (message_id)
);

CREATE TABLE IF NOT EXISTS upvotes (
    upvote_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL references users(user_id),
    message_id uuid NOT NULL references messages(message_id),
    created_at timestamp NOT NULL, -- DEFAULT NOW() NOT NULL,
    CONSTRAINT upvotes_pk PRIMARY KEY (upvote_id)
);

-- ?? derived tables ??

CREATE TABLE IF NOT EXISTS topics (
   topic_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
   name varchar(80) NOT NULL,
   description varchar(80) NOT NULL,
--   created_at timestamp DEFAULT NOW() NOT NULL,
   CONSTRAINT topics_pk PRIMARY KEY (topic_id)
);

CREATE TABLE IF NOT EXISTS pings (
   ping_id uuid NOT NULL, -- DEFAULT uuid_generate_v4() NOT NULL,
   user_id uuid NOT NULL references users(user_id),
   message_id uuid NOT NULL references messages(message_id),
--   created_at timestamp DEFAULT NOW() NOT NULL,
   CONSTRAINT pings_pk PRIMARY KEY (ping_id)
);
