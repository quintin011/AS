package controllers

import "gorm.io/gorm"

type Controller struct {
	DB *gorm.DB
}

func NewController(DB *gorm.DB) Controller {
	return Controller{DB}
<<<<<<< HEAD
}
=======
}
>>>>>>> v0.0.2
