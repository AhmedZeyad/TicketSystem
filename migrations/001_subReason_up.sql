create TAble subReasons(
id  int primary key auto_increment,
name  varchar(50),
reasonId int ,
foreign key (reasonId) REFERENCES reasons(id)
);