create Table IF NOT EXISTS users(
id int primary key auto_increment,
name varchar (20),
email varchar (255) NOT NULL UNIQUE,
rol_id int not null ,
password varchar (255),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP not null,
foreign key  (rol_id) REFERENCES roles(id) 
);