CREATE TABLE IF NOT EXISTS TicketImages (
  id int primary key auto_increment,
  fileName varchar(25) NOT NULL,
  filePath varchar(255) NOT NULL,
  ticketId int NOT NULL,
  uplode_by int NOT NULL,
  uplode_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  foreign key (tiketId) REFERENCES tickets(id),
  FOREIGN KEY (uplode_by) REFERENCES users(id)
)