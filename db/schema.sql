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
INSERT INTO status (id, name) VALUES (2, "Processing");
INSERT INTO status (id, name) VALUES (3, "Sent");
INSERT INTO status (id, name) VALUES (4, "Failed to Sent");
INSERT INTO status (id, name) VALUES (5, "Opened");

CREATE TABLE IF NOT EXISTS email (
  id INT NOT NULL AUTO_INCREMENT,
  `from` INT NOT NULL,
  `to` INT NOT NULL,
  subject VARCHAR(128) NOT NULL,
  body LONGTEXT NOT NULL,
  url_unsubscribe VARCHAR(256) NOT NULL,
  status INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  opened_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX email_from (`from`),
  INDEX email_to (`to`),
  INDEX email_status (status),
  FOREIGN KEY (`from`) REFERENCES address(id) ON DELETE CASCADE,
  FOREIGN KEY (`to`) REFERENCES address(id) ON DELETE CASCADE,
  FOREIGN KEY (status) REFERENCES status(id) ON DELETE CASCADE
) ENGINE=InnoDB;
