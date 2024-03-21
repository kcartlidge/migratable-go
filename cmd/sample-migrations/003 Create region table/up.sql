CREATE TABLE IF NOT EXISTS region
(
	id         BIGSERIAL                NOT NULL,
	name       CHARACTER VARYING(200)   NOT NULL,
	slug       CHARACTER VARYING(200)   NOT NULL,

	PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS ix_region_slug
	ON region USING btree (slug ASC NULLS LAST)
	WITH (deduplicate_items=True);

COMMENT ON TABLE region IS 'Random world regions.';
