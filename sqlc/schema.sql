DROP DATABASE IF EXISTS kiber;
CREATE DATABASE kiber;
USE kiber;

CREATE TABLE card (
  uid VARCHAR(256) NOT NULL,
  user_name VARCHAR(256) NOT NULL,
  PRIMARY KEY (uid)
);

INSERT INTO card
(uid, user_name)
VALUES
('48972568', 'vardenis pavardenis');
