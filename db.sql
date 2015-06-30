DROP DATABASE IF EXISTS `persona` ;
CREATE SCHEMA `persona` CHARACTER SET UTF8;

CREATE TABLE `persona`.`project` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(120) NOT NULL,
  `start_date` DATE NOT NULL,
  `end_date` DATE NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `persona`.`person` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(120) NOT NULL,
  `email` VARCHAR(200) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `person_email_unq` (`email` ASC)
);

CREATE TABLE `persona`.`cost_center` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(120) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `cost_center_name_unq` (`name` ASC)
);

CREATE TABLE `persona`.`person_allocation` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `person_id` INT NOT NULL,
  `cost_center_id` INT NOT NULL,
  `allocation` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`person_id`) REFERENCES `persona`.`person` (`id`),
  FOREIGN KEY (`cost_center_id`) REFERENCES `persona`.`cost_center` (`id`) 
);

CREATE TABLE `persona`.`ability` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `ability_name_unq` (`name` ASC)
);

CREATE TABLE `persona`.`experience` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `experience_name_unq` (`name` ASC)
);

CREATE TABLE `persona`.`person_abilities` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `person_id` INT NOT NULL,
  `ability_id` INT NOT NULL,
  `experience_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`person_id`) REFERENCES `persona`.`person` (`id`),
  FOREIGN KEY (`ability_id`) REFERENCES `persona`.`ability` (`id`),
  FOREIGN KEY (`experience_id`) REFERENCES `persona`.`experience` (`id`),
  UNIQUE INDEX `person_ability_unq` (`person_id`, `ability_id`, `experience_id`)
);
