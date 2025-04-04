package users

import (
	// "log"
	"errors"
	"time"

	"github.com/AhmedZeyad/TicketSystem/utilities"


)

func getUsers() ([]User, error) {
	var users []User
	println("hi form db 124")
	err := utilities.DB.Select(&users, "SELECT * from Users ")
	if err != nil {

		return nil, err
	}
	return users, nil

}
func(u *User) getUserById()  error {

	err := utilities.DB.Get(u, "SELECT * from Users WHERE id = ? LIMIT 1", u.ID)
	if err != nil {

		return  err
	}
	return  nil
}

func (u *User) InsertUser() error {
	u.CreatedAt = time.Now()
	resoult, err := utilities.DB.Exec("INSERT INTO Users (id, name, email, created_at, role_id) VALUES (null,?,?,?,?)", u.Name, u.Email, u.CreatedAt, u.RoleID)
	if err != nil {
		println(err.Error(), time.Now().String())
		return err
	}
	id, err := resoult.LastInsertId()
	if err != nil {
		return err

	}
	u.ID = int(id)
print(u.ID)
print(u.CreatedAt.String())
	return nil
}
func(u *User)EditUser()error{
	currentUser:=User{ID: u.ID}
	err:=currentUser.getUserById()
if err!=nil{
	return errors.New("user not found")
	
}
kk_,err := utilities.DB.Exec("UPDATE Users SET name = ?, email = ?  WHERE id = ?", u.Name, u.Email, u.ID)
if err!=nil{
	println("erer",kk_,err.Error())
	return errors.New("can't update user")
}
*u=currentUser
return  nil
}
func dropUser(id int) error {
	result, err := utilities.DB.Exec("DELETE FROM Users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err	 // unexpected error checking rows
	}

	if rowsAffected == 0 {
		return errors.New("user not found",)
	}

	return nil

}