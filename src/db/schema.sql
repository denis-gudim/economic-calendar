DROP TABLE IF EXISTS languages CASCADE;
DROP TABLE IF EXISTS countries CASCADE;
DROP TABLE IF EXISTS country_translations CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS event_translations CASCADE;
DROP TABLE IF EXISTS event_schedule CASCADE;
DROP TABLE IF EXISTS event_schedule_translations CASCADE;

/* Languages and ISO 639-1 codes */
CREATE TABLE languages
(
	id			INTEGER NOT NULL,
	code		VARCHAR(16) NOT NULL,
	name		VARCHAR(256) NOT NULL,
	native_name	VARCHAR(256) NOT NULL,
	domain		VARCHAR(16) NOT NULL,
	CONSTRAINT pk_languages PRIMARY KEY (id)
);

CREATE UNIQUE INDEX iux_languages_code
	ON languages (code);

CREATE UNIQUE INDEX iux_languages_domain
	ON languages (domain);

/* Countries and ISO 3166-1 alpha-2 codes*/
CREATE TABLE countries
(
	id				INTEGER NOT NULL,
	code			CHAR(2) NOT NULL,
	continent_code	CHAR(2)	NOT NULL,
	name			VARCHAR(256) NOT NULL,
	currency		CHAR(3) NOT NULL,
	CONSTRAINT pk_countries PRIMARY KEY (id)
);

CREATE UNIQUE INDEX iux_countries_code
	ON countries (code);

CREATE UNIQUE INDEX iux_countries_name
	ON countries (name);


/* Country title translations */
CREATE TABLE country_translations
(
	country_id		INTEGER NOT NULL,
	language_id		INTEGER NOT NULL,
	title			VARCHAR(256) NOT NULL,
	CONSTRAINT pk_country_translations PRIMARY KEY (country_id, language_id),
	CONSTRAINT fk_country_translations_countries FOREIGN KEY(country_id)
		REFERENCES countries ON DELETE CASCADE,
	CONSTRAINT fk_country_translations_languages FOREIGN KEY(language_id)
		REFERENCES languages ON DELETE CASCADE
);

/* Calendar events main data */
CREATE TABLE events
(
	id				INTEGER NOT NULL,
	country_id		INTEGER,
	impact_level	VARCHAR(10) NOT NULL,
	unit			VARCHAR(16),
	source			TEXT,
	source_url		TEXT,
	CONSTRAINT pk_events PRIMARY KEY (id),
	CONSTRAINT fk_events_countries FOREIGN KEY(country_id)
		REFERENCES countries ON DELETE CASCADE
);

CREATE INDEX ix_events_country_id ON events (country_id);

CREATE INDEX ix_events_impact_level ON events (impact_level);

/* Calendar events main data */
CREATE TABLE event_translations
(
	event_id			INTEGER NOT NULL,
	language_id			INTEGER NOT NULL,
	title				VARCHAR(1024) NOT NULL,
	overview			TEXT,
	CONSTRAINT pk_event_translations PRIMARY KEY (event_id, language_id),
	CONSTRAINT fk_event_translations_events FOREIGN KEY(event_id)
		REFERENCES events ON DELETE CASCADE,
	CONSTRAINT fk_event_translations_languages FOREIGN KEY(language_id)
		REFERENCES languages ON DELETE CASCADE
);

/* Calendar schedule history*/
CREATE TABLE event_schedule
(
	id				INTEGER NOT NULL,
	timestamp_utc 	TIMESTAMP NOT NULL,
	actual			DOUBLE PRECISION,
	forecast		DOUBLE PRECISION,
	previous		DOUBLE PRECISION,
	done			BOOLEAN NOT NULL,
	type			INTEGER NOT NULL,
	event_id		INTEGER,
	CONSTRAINT pk_event_schedule PRIMARY KEY (id),
	CONSTRAINT fk_event_schedule_events FOREIGN KEY(event_id)
		REFERENCES events ON DELETE CASCADE
);

CREATE INDEX ix_event_schedule_event_id ON event_schedule (event_id);

CREATE INDEX ix_event_schedule_done ON event_schedule (done);

CREATE INDEX ix_event_schedule_timestamp_utc ON event_schedule (timestamp_utc DESC);

/* Calendar schedule event translations*/
CREATE TABLE event_schedule_translations
(
	event_schedule_id	INTEGER NOT NULL,
	language_id			INTEGER NOT NULL,
	title				VARCHAR(1024) NOT NULL,
	CONSTRAINT pk_event_schedule_translations PRIMARY KEY (event_schedule_id, language_id),
	CONSTRAINT fk_event_schedule_translations_event_schedule FOREIGN KEY(event_schedule_id)
		REFERENCES event_schedule ON DELETE CASCADE,
	CONSTRAINT fk_event_schedule_translations_languages FOREIGN KEY(language_id)
		REFERENCES languages ON DELETE CASCADE
);