CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username varchar(30) NOT NULL,
    password_hash bytea NOT NULL,
    salt_for_keys bytea NOT NULL,
    public_key_for_sharing varchar NOT NULL,
    encrypted_public_key_for_sharing varchar NOT NULL
);

create UNIQUE INDEX ON users (username);

CREATE TABLE diaries
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    owner_id uuid REFERENCES users (id) NOT NULL
);

CREATE UNIQUE INDEX ON diaries (name, owner_id);

CREATE TABLE diary_keys
(
    diary_id uuid REFERENCES diaries (id) NOT NULL,
    user_id uuid REFERENCES users (id) NOT NULL,
    encrypted_key bytea NOT NULL,
    PRIMARY KEY (diary_id, user_id)
);

CREATE TABLE diary_entries
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    diary_id uuid REFERENCES diaries (id) NOT NULL,
    name text NOT NULL,
    date date NOT NULL
);

CREATE TABLE diary_entry_contents
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    diary_entry_id uuid REFERENCES diary_entries (id) NOT NULL,
    value jsonb NOT NULL
);

CREATE TABLE sharing_tasks
(
    diary_id uuid REFERENCES diaries (id) NOT NULL,
    receiver_user_id uuid REFERENCES users (id) NOT NULL,
    encrypted_diary_key bytea NOT NULL,
    shared_at timestamp NOT NULL,
    PRIMARY KEY (diary_id, receiver_user_id)
);

CREATE VIEW diaries_with_entries AS
SELECT *,
       (SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(entries.*))) AS array_to_json
        FROM (SELECT * FROM diary_entries WHERE diary_id = diaries.id) entries) AS entries
FROM diaries;