package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdateAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Isi Usernamenya!")
		}
		if u.Email == "" {
			return errors.New("Isi Emailnya!")
		}
		if u.Password == "" {
			return errors.New("Isi Passwordnya!")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email Tidak Valid!")
		}

		return nil
	case "login":
		if u.Email == "" {
			return errors.New("Isi Emailnya!")
		}
		if u.Password == "" {
			return errors.New("Isi Passwordnya!")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email Tidak Valid!")
		}
		return nil
	default:
		if u.Username == "" {
			return errors.New("Isi Usernamenya!")
		}
		if u.Email == "" {
			return errors.New("Isi Emailnya!")
		}
		if u.Password == "" {
			return errors.New("Isi Passwordnya!")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email Tidak Valid!")
		}
		return nil
	}
}

func (u *User) StoreUser(db *gorm.DB) (*User, error) {
	var err error

	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) AllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}

	return &users, err
}

func (u *User) UserById(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User not Found")
	}

	return u, err
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {

	err := u.BeforeSave()

	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"username":   u.Username,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, err
}

func (u *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}