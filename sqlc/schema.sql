DROP DATABASE IF EXISTS kiber;
CREATE DATABASE kiber;
USE kiber;

CREATE TABLE card (
  uid VARCHAR(256) NOT NULL,
  user_name VARCHAR(256) NOT NULL,
  PRIMARY KEY (uid)
);

CREATE TABLE log (
  id INT NOT NULL AUTO_INCREMENT,
  uid VARCHAR(256) NOT NULL,
  permitted BOOL NOT NULL,
  time timestamp NOT NULL,
  PRIMARY KEY(id)
);

INSERT INTO card
(uid, user_name)
VALUES
('48972568', 'vardenis pavardenis');
