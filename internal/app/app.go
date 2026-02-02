package app

import (
	"html/template"
	"log"
)

type Application struct {
	Logger *log.Logger

	Templates map[string]*template.Template

	Users       UserStore
	Sessions    SessionStore
	Restaurants RestaurantStore
	MenuItems   MenuItemStore
}
