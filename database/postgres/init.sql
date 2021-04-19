CREATE DATABASE test1;
\connect test1

CREATE TABLE account1(
   user_id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (355) UNIQUE NOT NULL,
   created_on TIMESTAMP NOT NULL,
   last_login TIMESTAMP
);

CREATE TABLE link (
   ID serial PRIMARY KEY,
   url VARCHAR (255) NOT NULL,
   name VARCHAR (255) NOT NULL,
   description VARCHAR (255),
   rel VARCHAR (50)
);

INSERT INTO link (url, name)
VALUES
   ('http://www.postgresqltutorial.com','PostgreSQL Link 1.1'),
   ('http://www.postgresqltutorial.com','PostgreSQL Link 1.2');

CREATE DATABASE test2;
\connect test2

CREATE TABLE account2(
   user_id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (355) UNIQUE NOT NULL,
   created_on TIMESTAMP NOT NULL,
   last_login TIMESTAMP
);

CREATE TABLE link (
   ID serial PRIMARY KEY,
   url VARCHAR (255) NOT NULL,
   name VARCHAR (255) NOT NULL,
   description VARCHAR (255),
   rel VARCHAR (50)
);

INSERT INTO link (url, name)
VALUES
   ('http://www.postgresqltutorial.com','PostgreSQL Link 2.1'),
   ('http://www.postgresqltutorial.com','PostgreSQL Link 2.2');

CREATE DATABASE test3;
\connect test3

CREATE TABLE account3(
   user_id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (355) UNIQUE NOT NULL,
   created_on TIMESTAMP NOT NULL,
   last_login TIMESTAMP
);

CREATE TABLE link (
   ID serial PRIMARY KEY,
   url VARCHAR (255) NOT NULL,
   name VARCHAR (255) NOT NULL,
   description VARCHAR (255),
   rel VARCHAR (50)
);

INSERT INTO link (url, name)
VALUES
   ('http://www.postgresqltutorial.com','PostgreSQL Link 3.1'),
   ('http://www.postgresqltutorial.com','PostgreSQL Link 3.2');
