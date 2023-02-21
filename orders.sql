CREATE DATABASE IF NOT EXISTS orders;

USE orders;

CREATE TABLE IF NOT EXISTS orders (
  id varchar(36) NOT NULL,
  status varchar(255) NOT NULL,
  items json NOT NULL,
  total float(10,2) NOT NULL,
  currency_unit varchar(255) NOT NULL,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL,
  PRIMARY KEY (id)
);