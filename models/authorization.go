package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Ability struct {
	ID            uint   `json:"id"`
	Key           string `json:"key"`
	Value         string `json:"value"`
	ModuleID      uint   `json:"module_id"`
	ModuleSubject string `gorm:"<-:false;->;-:migration;omitempty;" json:"module_subject"`
}

type UserAbilities struct {
	UserID    uint `gorm:"primary_key"  json:"user_id" `
	AbilityID uint `gorm:"primary_key" json:"ability_id"  `
	RoleID    int  `json:"role_id"`
}

type User struct {
	ID            uint                   `json:"id"`
	Email         string                 `gorm:"unique" json:"email" validate:"required,email"`
	Password      []byte                 `json:"-"`
	UserRole      []Role                 `gorm:"many2many:user_roles;" json:"user_role,omitempty"`
	Abilities     []Ability              `gorm:"many2many:user_abilities;" json:"-"`
	CaslAbilities map[string]interface{} `gorm:"-" json:"casl_abilities,omitempty"`
	Level         uint                   `gorm:"-" json:"level"`
}

type Role struct {
	ID        uint      `json:"id"`
	Abilities []Ability `gorm:"many2many:roles_abilities;" json:",omitempty"`
	Name      string    `gorm:"type:character varying(50)" json:"name"`
	RoleLevel uint      `gorm:"type:integer" json:"role_level"`
}

func HashPassword(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, 14)
	return bytes, err
}

func CheckPasswordHash(password []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

func (User) Seed(db *gorm.DB) {
	password, err := HashPassword([]byte("auction."))
	if err != nil {
		panic(err)
	}
	user := User{
		Email:    "info@auction.com",
		Password: password,
		UserRole: []Role{
			{
				Name:      "admin",
				RoleLevel: 1,
			},
		},
	}

	seller := Seller{
		User: user,
		Identification: Identification{
			Name:      "Yasin",
			Surname:   "Yumrutepe",
			BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Email: "seller@seller.com",
	}

	db.Create(&seller)
	var abilitiesID []uint
	db.Model(&Ability{}).Pluck("id", &abilitiesID)

	// if err != nil {
	// 	panic(err)
	// }
}
