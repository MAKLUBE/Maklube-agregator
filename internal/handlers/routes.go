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

	router.HandlerFunc("GET", "/verify", h.verifyForm)
	router.HandlerFunc("POST", "/verify", h.verifyPost)

	router.HandlerFunc("GET", "/login", h.loginForm)
	router.HandlerFunc("POST", "/login", h.loginPost)
	router.Handler("POST", "/logout", mw.RequireAuth(http.HandlerFunc(h.logoutPost)))

	router.HandlerFunc("GET", "/restaurants", h.restaurantsList)
	router.HandlerFunc("GET", "/restaurants/:id", h.restaurantView)
	router.Handler("POST", "/restaurants/:id/reviews", mw.RequireAuth(http.HandlerFunc(h.reviewCreate)))
	router.Handler("POST", "/restaurants/:id/orders", mw.RequireAuth(http.HandlerFunc(h.orderCreate)))
	router.Handler("GET", "/orders", mw.RequireAuth(http.HandlerFunc(h.customerOrdersList)))

	router.Handler("GET", "/partner/orders", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerIncomingOrders)))
	router.Handler("POST", "/partner/orders/:id/status", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerOrderUpdateStatus)))
	router.Handler("GET", "/partner/restaurants", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerRestaurants)))
	router.Handler("GET", "/partner/restaurants/new", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.restaurantCreateForm)))
	router.Handler("POST", "/partner/restaurants/new", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.restaurantCreatePost)))

	router.Handler("GET", "/partner/restaurants/id/:id/edit", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerRestaurantEditForm)))
	router.Handler("POST", "/partner/restaurants/id/:id/edit", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerRestaurantEditPost)))

	router.Handler("GET", "/partner/restaurants/id/:id/menu", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerMenuList)))
	router.Handler("GET", "/partner/restaurants/id/:id/menu/new", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerMenuNewForm)))
	router.Handler("POST", "/partner/restaurants/id/:id/menu/new", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerMenuNewPost)))
	router.Handler("GET", "/partner/restaurants/id/:id/menu-items/:item/edit", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerMenuEditForm)))
	router.Handler("POST", "/partner/restaurants/id/:id/menu-items/:item/edit", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerMenuEditPost)))

	router.Handler("GET", "/partner/restaurants/id/:id/halal", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerHalalRequestForm)))
	router.Handler("POST", "/partner/restaurants/id/:id/halal", mw.RequireRole(models.RolePartner, http.HandlerFunc(h.partnerHalalRequestPost)))

	router.Handler("GET", "/admin/halal", mw.RequireRole(models.RoleAdmin, http.HandlerFunc(h.adminHalalRequests)))
	router.Handler("POST", "/admin/halal/:id/approve", mw.RequireRole(models.RoleAdmin, http.HandlerFunc(h.adminHalalApprove)))
	router.Handler("POST", "/admin/halal/:id/reject", mw.RequireRole(models.RoleAdmin, http.HandlerFunc(h.adminHalalReject)))
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
