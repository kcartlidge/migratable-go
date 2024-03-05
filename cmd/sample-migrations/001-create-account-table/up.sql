CREATE TABLE IF NOT EXISTS account
(
	id            BIGSERIAL                NOT NULL,
	created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,
	email_address CHARACTER VARYING(200)   NOT NULL,
	username      CHARACTER VARYING(200)   NOT NULL,
	is_admin      BOOLEAN                  NOT NULL DEFAULT FALSE,
	verified_at   TIMESTAMP WITH TIME ZONE,
	last_login_at TIMESTAMP WITH TIME ZONE,
	closed_at     TIMESTAMP WITH TIME ZONE,

	PRIMARY KEY (id),
	CONSTRAINT uq_account_email_address UNIQUE (email_address),
	CONSTRAINT uq_account_username UNIQUE (username)
);

CREATE INDEX IF NOT EXISTS ix_account_email_address
	ON account USING btree (email_address ASC NULLS LAST)
	WITH (deduplicate_items=True);

COMMENT ON TABLE account IS 'User accounts.';
