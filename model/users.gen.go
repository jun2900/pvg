// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	UserID      int32     `gorm:"column:user_id;primaryKey;autoIncrement:true" json:"user_id"`
	Username    string    `gorm:"column:username;not null" json:"username"`
	Password    string    `gorm:"column:password;not null" json:"password"`
	Firstname   string    `gorm:"column:firstname;not null" json:"firstname"`
	Lastname    string    `gorm:"column:lastname;not null" json:"lastname"`
	PhoneNumber string    `gorm:"column:phone_number;not null" json:"phone_number"`
	Email       string    `gorm:"column:email;not null" json:"email"`
	Birthday    time.Time `gorm:"column:birthday;not null" json:"birthday"`
	Status      string    `gorm:"column:status;not null;default:PENDING" json:"status"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}