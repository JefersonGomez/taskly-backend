package models

import "gorm.io/gorm"

func MigrarTablas(db *gorm.DB) {
	db.AutoMigrate(&Usuario{})
	db.AutoMigrate(&Equipo{})
	db.AutoMigrate(&Miembro{})
	db.AutoMigrate(&Tarea{})
}
