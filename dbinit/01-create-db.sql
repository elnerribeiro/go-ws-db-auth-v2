\c test;
drop table if exists user_db;
drop table if exists insert_batch;
drop table if exists ins_id;
create table user_db (id serial not null, email varchar(50) not null, role varchar(20) not null, password varchar(128) not null, primary key (id));
create table ins_id (id serial not null, type varchar(20) not null, quantity int not null, status varchar(20) not null, tstampinit bigint, tstampend bigint, primary key (id));
create table insert_batch(id serial not null, id_ins_id int not null, pos int not null, primary key(id), foreign key (id_ins_id) references ins_id(id));
--PWD abc
insert into user_db (email, role, password) values ('user@user.com', 'user', 'DDAF35A193617ABACC417349AE20413112E6FA4E89A97EA20A9EEEE64B55D39A2192992A274FC1A836BA3C23A3FEEBBD454D4423643CE80E2A9AC94FA54CA49F');
--PWD 123
insert into user_db (email, role, password) values ('admin@admin.com', 'admin', '3C9909AFEC25354D551DAE21590BB26E38D53F2173B8D3DC3EEE4C047E7AB1C1EB8B85103E3BE7BA613B31BB5C9C36214DC9F14A42FD7A2FDB84856BCA5C44C2');
