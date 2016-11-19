DROP SCHEMA IF EXISTS "places4all" CASCADE;
CREATE SCHEMA "places4all";
SET SCHEMA 'places4all';

CREATE TABLE Country (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE Address (
  id SERIAL PRIMARY KEY,
  idCountry INTEGER NOT NULL REFERENCES Country(id),
  addressLine1 VARCHAR(200) NOT NULL,
  addressLine2 VARCHAR(200),
  addressLine3 VARCHAR(200),
  townCity VARCHAR(100),
  county VARCHAR(150),
  postCode VARCHAR(50),
  latitude DECIMAL(8,6),
  longitude DECIMAL(9,6)
);

CREATE TABLE Tag (
  id SERIAL PRIMARY KEY,
  name VARCHAR(30) UNIQUE NOT NULL
);

CREATE TABLE Person (
  id SERIAL PRIMARY KEY,
  idCountry INTEGER NOT NULL REFERENCES Country(id),
  name VARCHAR(100) NOT NULL,
  email VARCHAR(254) UNIQUE NOT NULL
    CHECK (email ~* '^[^\s@]+@[^\s@]+\.[^\s@.]+$'),
  username VARCHAR UNIQUE NOT NULL
    CHECK (username ~* '^[A-Za-z][A-Za-z0-9\.\-_]{2,15}$'),
  password VARCHAR NOT NULL
    CHECK (LENGTH(password) >= 6),
  imagePath VARCHAR(255),
  banned BOOLEAN,
  dateBanned TIMESTAMP,
  reason TEXT,
  mobilephone VARCHAR(20),
  telephone VARCHAR(20),
  created TIMESTAMP NOT NULL
);

CREATE TABLE SuperAdmin (
  id SERIAL PRIMARY KEY,
  idPerson INTEGER NOT NULL REFERENCES Person(id)
);

CREATE TABLE LocalAdmin (
  id SERIAL PRIMARY KEY,
  idPerson INTEGER NOT NULL REFERENCES Person(id)
);

CREATE TABLE Auditor (
  id SERIAL PRIMARY KEY,
  idPerson INTEGER NOT NULL REFERENCES Person(id)
);

CREATE TABLE Client (
  id SERIAL PRIMARY KEY,
  idPerson INTEGER NOT NULL REFERENCES Person(id)
);

CREATE TABLE Property (
  id SERIAL PRIMARY KEY,
  idAddress INTEGER NOT NULL REFERENCES Address(id),
  name VARCHAR(150) NOT NULL,
  details TEXT NOT NULL,
  created TIMESTAMP NOT NULL
);

CREATE TABLE Gallery (
  id SERIAL PRIMARY KEY,
  idProperty INTEGER NOT NULL REFERENCES Property(id),
  name VARCHAR(100) NOT NULL,
  description TEXT,
  created TIMESTAMP NOT NULL
);

CREATE TABLE Image (
  id SERIAL PRIMARY KEY,
  idGallery INTEGER NOT NULL REFERENCES Gallery(id),
  name VARCHAR(100) NOT NULL,
  description TEXT,
  imagePath VARCHAR(255) NOT NULL,
  created TIMESTAMP NOT NULL
);

CREATE TABLE PropertyTag (
  idProperty INTEGER REFERENCES Property(id),
  idTag INTEGER REFERENCES Tag(id),
  PRIMARY KEY(idProperty, idTag)
);

CREATE TABLE PropertyClient (
  idProperty INTEGER REFERENCES Property(id),
  idClient INTEGER REFERENCES Client(id),
  PRIMARY KEY(idProperty, idClient)
);

CREATE TABLE Template (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  created TIMESTAMP NOT NULL
);

CREATE TABLE Maingroup (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  weight INTEGER NOT NULL,
  description TEXT,
  imagePath VARCHAR(255),
  created TIMESTAMP NOT NULL
);

CREATE TABLE TemplateMaingroup (
  idTemplate INTEGER REFERENCES Template(id),
  idMaingroup INTEGER REFERENCES Maingroup(id),
  PRIMARY KEY(idTemplate, idMaingroup)
);

CREATE TABLE Subgroup (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  weight INTEGER NOT NULL,
  description TEXT,
  imagePath VARCHAR(100),
  created TIMESTAMP NOT NULL
);

CREATE TABLE MaingroupSubgroup (
  idMaingroup INTEGER REFERENCES Maingroup(id),
  idSubgroup INTEGER REFERENCES Subgroup(id),
  PRIMARY KEY(idMaingroup, idSubgroup)
);

CREATE TABLE Legislation (
  id SERIAL PRIMARY KEY,
  name VARCHAR(200),
  url VARCHAR(2083)
);

CREATE TABLE Criterion (
  id SERIAL PRIMARY KEY,
  idLegislation INTEGER REFERENCES Legislation(id),
  name VARCHAR(100) NOT NULL,
  weight INTEGER NOT NULL,
  description TEXT,
  imagePath VARCHAR(255),
  created TIMESTAMP NOT NULL
);

CREATE TABLE SubgroupCriterion (
  idSubgroup INTEGER REFERENCES Subgroup(id),
  idCriterion INTEGER REFERENCES Criterion(id),
  PRIMARY KEY(idSubgroup, idCriterion)
);

CREATE TABLE Accessibility (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  imagePath VARCHAR(255)
);

CREATE TABLE CriterionAccessibility (
  idCriterion INTEGER REFERENCES Criterion(id),
  idAccessibility INTEGER REFERENCES Accessibility(id),
  weight INTEGER NOT NULL,
  PRIMARY KEY(idCriterion, idAccessibility)
);

CREATE TABLE Audit (
  id SERIAL PRIMARY KEY,
  idProperty INTEGER NOT NULL REFERENCES Property(id),
  idAuditor INTEGER NOT NULL REFERENCES Auditor(id),
  idTemplate INTEGER NOT NULL REFERENCES Template(id),
  rating INTEGER,
  observation TEXT,
  created TIMESTAMP NOT NULL,
  finished TIMESTAMP
);

CREATE TABLE AuditCriterion (
  idAudit INTEGER REFERENCES Audit(id),
  idCriterion INTEGER REFERENCES Criterion(id),
  value INTEGER,
  observation VARCHAR,
  PRIMARY KEY(idAudit, idCriterion)
);
