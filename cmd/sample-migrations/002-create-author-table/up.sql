CREATE TABLE IF NOT EXISTS author
(
	id         BIGSERIAL                NOT NULL,
	account_id BIGINT                   NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
	name       CHARACTER VARYING(200)   NOT NULL,
	slug       CHARACTER VARYING(200)   NOT NULL,
	sort_name  CHARACTER VARYING(200)   NOT NULL,

	PRIMARY KEY (id),
	CONSTRAINT fk_author_account FOREIGN KEY (account_id)
		REFERENCES account (id) MATCH SIMPLE
		ON UPDATE NO ACTION ON DELETE NO ACTION NOT VALID
);

CREATE INDEX IF NOT EXISTS ix_author_account_id
	ON author USING btree (account_id ASC NULLS LAST)
	WITH (deduplicate_items=True);

CREATE INDEX IF NOT EXISTS ix_author_name
	ON author USING btree (name ASC NULLS LAST)
	WITH (deduplicate_items=True);

CREATE INDEX IF NOT EXISTS ix_author_slug
	ON author USING btree (slug ASC NULLS LAST)
	WITH (deduplicate_items=True);

COMMENT ON TABLE author IS 'Authors per account.';
