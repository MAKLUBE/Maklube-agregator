package handlers

import (
	"net/http"

	"github.com/MAKLUBE/AP1_Final_Project/internal/app"
	"github.com/MAKLUBE/AP1_Final_Project/internal/middleware"
	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"github.com/julienschmidt/httprouter"
)

func Routes(a *app.Application) http.Handler {
	mw := middleware.New(a)
	router := httprouter.New()
	h := &Handler{App: a}
	router.HandlerFunc("GET", "/", h.home)
	router.HandlerFunc("GET", "/register", h.registerForm)
	router.HandlerFunc("POST", "/register", h.registerPost)

	router.HandlerFunc("GET", "/login", h.loginForm)
	router.HandlerFunc("POST", "/login", h.loginPost)
	router.Handler("POST", "/logout", mw.RequireAuth(http.HandlerFunc(h.logoutPost)))

	router.HandlerFunc("GET", "/restaurants", h.restaurantsList)
	router.HandlerFunc("GET", "/restaurants/:id", h.restaurantView)
	router.Handler("POST", "/restaurants/:id/reviews", mw.RequireAuth(http.HandlerFunc(h.reviewCreate)))

	router.Handler("GET", "/partner/orders", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerIncomingOrders)))
	router.Handler("GET", "/partner/restaurants/new", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.restaurantCreateForm)))
	router.Handler("POST", "/partner/restaurants/new", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.restaurantCreatePost)))

	router.Handler("GET", "/admin/halal", mw.RequireRole(models.RoleAdmin, http.HandlerFunc(h.adminHalalRequests)))
	router.Handler("GET", "/admin/reviews", mw.RequireRole(models.RoleAdmin, http.HandlerFunc(h.adminReviewsList)))
	router.Handler("POST", "/admin/reviews/:id/hide", mw.RequireRole(models.RoleAdmin, http.HandlerFunc(h.adminReviewHide)))
	router.Handler("POST", "/admin/reviews/:id/show", mw.RequireRole(models.RoleAdmin, http.HandlerFunc(h.adminReviewShow)))

	fs := http.FileServer(http.Dir("./ui/static"))
	router.Handler("GET", "/static/*filepath", http.StripPrefix("/static/", fs))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	})

	return router
}

type Handler struct {
	App *app.Application
}
