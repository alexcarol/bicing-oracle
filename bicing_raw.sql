CREATE TABLE station
(
id int PRIMARY KEY,
latitude float,
longitude float,
street varchar(255),
height int,
street_number varchar(255),
nearby_station_list varchar(255),
last_updatetime int
);

CREATE TABLE station_state
(
id int,
updatetime int,
slots int,
bikes int,
FOREIGN KEY (id)
	REFERENCES station(id)
	ON DELETE CASCADE,
PRIMARY KEY (id, updatetime)
);
