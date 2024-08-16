CREATE DATABASE go-user-settings-example

USE go-user-settings-example

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  name varchar(50) NOT NULL,
  authority varchar(10) NOT NULL,
  email varchar(50) NOT NULL UNIQUE,
  password varchar(100) NOT NULL UNIQUE,
  createdAt timestamp NOT NULL DEFAULT NOW(),
  updatedAt timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE user_settings (
  id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  user_id uuid REFERENCES users(id),
  theme varchar(10) NOT NULL DEFAULT 'dark',
);
