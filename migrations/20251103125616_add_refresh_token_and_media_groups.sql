-- +goose Up
-- +goose StatementBegin

CREATE TYPE group_members_roles AS ENUM ('owner', 'editor', 'viewer');
CREATE TYPE file_status AS ENUM ('new', 'active', 'uploading', 'deleted');
CREATE TYPE mime_types AS ENUM ('image/jpeg', 'image/png', 'image/gif', 'video/mp4');

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    token TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    issued_at TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP without time zone NOT NULL,
    revoked BOOLEAN DEFAULT FALSE
);

CREATE TABLE groups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    owner_id        UUID REFERENCES users(id) ON DELETE SET NULL,
    is_shared       BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);


CREATE TABLE group_members (
   group_id    UUID REFERENCES groups(id) ON DELETE CASCADE,
   user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
   role        group_members_roles DEFAULT 'viewer',
   PRIMARY KEY (group_id, user_id)
);

CREATE TABLE files (
   id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   group_id        UUID REFERENCES groups(id) ON DELETE CASCADE,
   uploader_id     UUID REFERENCES users(id) ON DELETE SET NULL,
   file_name       VARCHAR(255) NOT NULL,
   storage_path    VARCHAR(255) NOT NULL,
   mime_type       mime_types,
   size_bytes      BIGINT,
   uploaded_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
   status          file_status DEFAULT 'new',
   is_shared       BOOLEAN DEFAULT FALSE,
   meta            JSONB DEFAULT '{}'::jsonb
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE group_members;
DROP TABLE files;
DROP TABLE groups;
DROP TABLE refresh_tokens;
DROP TYPE group_members_roles;
DROP TYPE file_status;
DROP TYPE mime_types;
-- +goose StatementEnd
