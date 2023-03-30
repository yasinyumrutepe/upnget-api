package models

import (
	"auction/globals"

	"gorm.io/gorm"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

type Module struct {
	ID           uint      `json:"id"`
	ParentID     *uint     `json:"parent_id"`
	ParentModule *Module   `gorm:"foreignKey:ParentID"`
	Name         string    `json:"name"`
	Subject      string    `json:"subject"`
	Abilities    []Ability `gorm:"foreignKey:ModuleID"`
	Children     []Module  `gorm:"foreignKey:ParentID"`
}

func (md Module) Seed(db *gorm.DB) {

	var modules []Module

	modules = append(modules, ORNEKMODULE)
	modules = append(modules, abilities)

	module := Module{
		ID:       uint(MaxInt - 1),
		Name:     "main",
		Subject:  "main",
		ParentID: globals.Ptr(uint(MaxInt - 1)),
		Children: modules,
	}
	err := db.Create(&module).Error
	if err != nil {
		panic(err)
	}
}

var ORNEKMODULE = Module{
	Subject: "module",
	Name:    "Modül",
	Abilities: []Ability{
		{
			Key:   "GET",
			Value: "Görüntüle",
		},
		{
			Key:   "POST",
			Value: "EKLE",
		},
		{
			Key:   "PUT",
			Value: "GÜNCELLE",
		},
		{
			Key:   "DELETE",
			Value: "SİL",
		},
	},
	Children: []Module{
		{
			Subject: "subModule",
			Name:    "Yetkilendirme",
			Abilities: []Ability{
				{
					Key:   "GET",
					Value: "Görüntüle",
				},
				{
					Key:   "POST",
					Value: "EKLE",
				},
				{
					Key:   "PUT",
					Value: "GÜNCELLE",
				},
				{
					Key:   "DELETE",
					Value: "SİL",
				},
			},
		},
	},
}

var abilities = Module{
	Subject: "abilities",
	Name:    "Yetkiler",
	Abilities: []Ability{
		{
			Key:   "GET",
			Value: "Görüntüle",
		},
		{
			Key:   "POST",
			Value: "EKLE",
		},
		{
			Key:   "PUT",
			Value: "GÜNCELLE",
		},
		{
			Key:   "DELETE",
			Value: "SİL",
		},
	},
}
