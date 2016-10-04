
-- +goose Up
CREATE TABLE IF NOT EXISTS address (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(128) NULL DEFAULT NULL,
  email VARCHAR(512) NOT NULL,
  is_sender TINYINT(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE INDEX address_email (email, is_sender)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS status (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(128) NULL DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE INDEX status_name (name)
) ENGINE=InnoDB;
TRUNCATE TABLE status;
INSERT INTO status (id, name) VALUES (1, "Pending");
INSERT INTO status (id, name) VALUES (2, "Sent");
INSERT INTO status (id, name) VALUES (3, "Failed to Sent");
INSERT INTO status (id, name) VALUES (4, "Opened");

CREATE TABLE IF NOT EXISTS email (
  id INT NOT NULL AUTO_INCREMENT,
  `from` INT NOT NULL,
  `to` INT NOT NULL,
  subject VARCHAR(128) NOT NULL,
  html TEXT NOT NULL,
  status INT NOT NULL DEFAULT 1,
  PRIMARY KEY (id),
  INDEX email_from (`from`),
  INDEX email_to (`to`),
  INDEX email_status (status),
  FOREIGN KEY (`from`) REFERENCES address(id) ON DELETE CASCADE,
  FOREIGN KEY (`to`) REFERENCES address(id) ON DELETE CASCADE,
  FOREIGN KEY (status) REFERENCES status(id) ON DELETE CASCADE
) ENGINE=InnoDB;


-- +goose Down
DROP TABLE IF EXISTS email;
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS status;
