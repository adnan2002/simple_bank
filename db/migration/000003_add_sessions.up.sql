CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  username VARCHAR NOT NULL,
  refresh_token VARCHAR NOT NULL,
  user_agent VARCHAR NOT NULL,
  client_ip VARCHAR NOT NULL,
  is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_sessions_user FOREIGN KEY (username) REFERENCES users (username)
);
