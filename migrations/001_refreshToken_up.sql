CREATE  TABLE IF NOT EXISTS  refreshToken(
id int primary key auto_increment,
userId int not null,
token varChar (255) not null,
foreign key (userid) REFERENCES users(id)
)