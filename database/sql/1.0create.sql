DROP ROLE IF EXISTS admin;
CREATE ROLE admin WITH LOGIN ENCRYPTED PASSWORD 'admin';
DROP SCHEMA IF EXISTS "places4all" CASCADE;
CREATE SCHEMA "places4all";
SET SCHEMA 'places4all';

CREATE TABLE country (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL,
  iso2 CHAR(2) UNIQUE NOT NULL
);

CREATE TABLE address (
  id SERIAL PRIMARY KEY,
  id_country INTEGER NOT NULL REFERENCES country(id),
  address_line1 VARCHAR(200) NOT NULL,
  address_line2 VARCHAR(200),
  address_line3 VARCHAR(200),
  town_city VARCHAR(100),
  county VARCHAR(150),
  postcode VARCHAR(50),
  latitude DECIMAL(8,6),
  longitude DECIMAL(9,6)
);

CREATE TABLE tag (
  id SERIAL PRIMARY KEY,
  name VARCHAR(30) UNIQUE NOT NULL
);

CREATE TABLE entity (
  id SERIAL PRIMARY KEY,
  id_country INTEGER NOT NULL REFERENCES country(id),
  name VARCHAR(100) NOT NULL,
  email VARCHAR(254) UNIQUE NOT NULL
    CHECK (email ~* '^[^\s@]+@[^\s@]+\.[^\s@.]+$'),
  username VARCHAR UNIQUE NOT NULL
    CHECK (username ~* '^[A-Za-z][A-Za-z0-9\.\-_]{2,15}$'),
  password VARCHAR NOT NULL
    CHECK (LENGTH(password) >= 6),
  image BYTEA,
  banned BOOLEAN,
  banned_date TIMESTAMP,
  reason TEXT,
  mobilephone VARCHAR(20),
  telephone VARCHAR(20),
  created_date TIMESTAMP NOT NULL
);

CREATE TABLE superadmin (
  id SERIAL PRIMARY KEY,
  id_entity INTEGER NOT NULL REFERENCES entity(id)
);

CREATE TABLE localadmin (
  id SERIAL PRIMARY KEY,
  id_entity INTEGER NOT NULL REFERENCES entity(id)
);

CREATE TABLE auditor (
  id SERIAL PRIMARY KEY,
  id_entity INTEGER NOT NULL REFERENCES entity(id)
);

CREATE TABLE client (
  id SERIAL PRIMARY KEY,
  id_entity INTEGER NOT NULL REFERENCES entity(id)
);

CREATE TABLE property (
  id SERIAL PRIMARY KEY,
  id_address INTEGER NOT NULL REFERENCES address(id),
  name VARCHAR(150) NOT NULL,
  details TEXT NOT NULL,
  created_date TIMESTAMP NOT NULL
);

CREATE TABLE gallery (
  id SERIAL PRIMARY KEY,
  id_property INTEGER NOT NULL REFERENCES property(id),
  name VARCHAR(100) NOT NULL,
  description TEXT,
  created_date TIMESTAMP NOT NULL
);

CREATE TABLE image (
  id SERIAL PRIMARY KEY,
  id_gallery INTEGER NOT NULL REFERENCES gallery(id),
  name VARCHAR(100) NOT NULL,
  description TEXT,
  image BYTEA NOT NULL,
  created_date TIMESTAMP NOT NULL
);

CREATE TABLE property_tag (
  id_property INTEGER REFERENCES property(id),
  id_tag INTEGER REFERENCES tag(id),
  PRIMARY KEY(id_property, id_tag)
);

CREATE TABLE property_client (
  id_property INTEGER REFERENCES property(id),
  id_client INTEGER REFERENCES client(id),
  PRIMARY KEY(id_property, id_client)
);

CREATE TABLE template (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  created_date TIMESTAMP NOT NULL
);

CREATE TABLE maingroup (
  id SERIAL PRIMARY KEY,
  id_template INTEGER NOT NULL REFERENCES template(id),
  name VARCHAR(100) NOT NULL,
  weight INTEGER NOT NULL,
  created_date TIMESTAMP NOT NULL
);

CREATE TABLE subgroup (
  id SERIAL PRIMARY KEY,
  id_maingroup INTEGER NOT NULL REFERENCES maingroup(id),
  name VARCHAR(100) NOT NULL,
  weight INTEGER NOT NULL,
  created_date TIMESTAMP NOT NULL
);


CREATE TABLE legislation (
  id SERIAL PRIMARY KEY,
  name VARCHAR(400) NOT NULL,
  description TEXT,
  url VARCHAR(2083)
);

CREATE TABLE criterion (
  id SERIAL PRIMARY KEY,
  id_subgroup INTEGER NOT NULL REFERENCES subgroup(id),
  id_legislation INTEGER REFERENCES legislation(id),
  name VARCHAR(100) NOT NULL,
  weight INTEGER NOT NULL,
  created_date TIMESTAMP NOT NULL
);

CREATE TABLE accessibility (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE criterion_accessibility (
  id_criterion INTEGER REFERENCES criterion(id),
  id_accessibility INTEGER REFERENCES accessibility(id),
  weight INTEGER NOT NULL,
  PRIMARY KEY(id_criterion, id_accessibility)
);

CREATE TABLE audit (
  id SERIAL PRIMARY KEY,
  id_property INTEGER NOT NULL REFERENCES property(id),
  id_auditor INTEGER NOT NULL REFERENCES auditor(id),
  id_template INTEGER NOT NULL REFERENCES template(id),
  rating INTEGER,
  observation TEXT,
  created_date TIMESTAMP NOT NULL,
  finished_date TIMESTAMP
);

CREATE TABLE audit_subgroup (
  id_audit INTEGER REFERENCES audit(id),
  id_subgroup INTEGER REFERENCES subgroup(id),
  PRIMARY KEY(id_audit, id_subgroup)
);

CREATE TABLE audit_criterion (
  id_audit INTEGER REFERENCES audit(id),
  id_criterion INTEGER REFERENCES criterion(id),
  value INTEGER,
  observation TEXT,
  PRIMARY KEY(id_audit, id_criterion)
);

GRANT ALL ON DATABASE "places4all" to admin;
GRANT ALL ON SCHEMA "places4all" TO admin;
GRANT ALL ON ALL TABLES IN SCHEMA "places4all" TO admin;
GRANT SELECT, USAGE ON ALL SEQUENCES IN SCHEMA "places4all" to admin;
