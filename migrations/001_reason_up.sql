CREATE TABLE IF NOT EXISTS reasons(
id int primary key auto_increment,
name varchar (50),
status bool not null default false,
created_at timestamp default current_timestamp not null defalult current_timestamp ,
created_by int not null ,
updated_at timestamp default current_timestamp not null,
updated_by int not null,
deleted_at  timestamp ,
deleted_by int 
);