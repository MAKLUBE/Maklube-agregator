package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) restaurantsList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.App.Restaurants.List(ctx)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.render(w, r, "restaurants.tmpl", &templateData{
		User: h.currentUser(r),
		Data: items,
	})
}

func (h *Handler) restaurantView(w http.ResponseWriter, r *http.Request) {
	// /restaurants/<id>
	parts := splitPath(r.URL.Path)
	if len(parts) < 2 {
		h.clientError(w, http.StatusBadRequest)
		return
	}
	idHex := parts[1]

	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rest, err := h.App.Restaurants.FindByID(ctx, oid)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	h.render(w, r, "restaurant_view.tmpl", &templateData{
		User: h.currentUser(r),
		Data: h.buildRestaurantViewData(ctx, rest, oid),
	})
}

func (h *Handler) restaurantCreateForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "create_restaurant.tmpl", &templateData{
		User: h.currentUser(r),
	})
}

func (h *Handler) restaurantCreatePost(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.PostForm.Get("name"))
	description := strings.TrimSpace(r.PostForm.Get("description"))
	phone := strings.TrimSpace(r.PostForm.Get("phone"))
	addressText := strings.TrimSpace(r.PostForm.Get("address_text"))
	city := strings.TrimSpace(r.PostForm.Get("city"))
	district := strings.TrimSpace(r.PostForm.Get("district"))
	mapURL := strings.TrimSpace(r.PostForm.Get("map_url"))
	instagram := strings.TrimSpace(r.PostForm.Get("instagram"))

	if name == "" || addressText == "" || city == "" {
		h.render(w, r, "create_restaurant.tmpl", &templateData{
			User: u,
			Form: map[string]string{"error": "Fill required fields"},
		})
		return
	}

	rest := &models.Restaurant{
		Name:        name,
		Description: description,
		Phone:       phone,
		HalalStatus: "pending",
		Address: struct {
			AddressText string `bson:"address_text" json:"address_text"`
			City        string `bson:"city" json:"city"`
			District    string `bson:"district,omitempty" json:"district,omitempty"`
			MapURL      string `bson:"map_url,omitempty" json:"map_url,omitempty"`
		}{
			AddressText: addressText,
			City:        city,
			District:    district,
			MapURL:      mapURL,
		},
		OwnerUserID: u.ID,
		Instagram:   instagram,
		CreatedAt:   time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.Restaurants.Insert(ctx, rest); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/restaurants", http.StatusSeeOther)
}

func (h *Handler) reviewCreate(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	parts := splitPath(r.URL.Path)
	if len(parts) < 3 {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	oid, err := primitive.ObjectIDFromHex(parts[1])
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	ratingStr := strings.TrimSpace(r.PostForm.Get("rating"))
	comment := strings.TrimSpace(r.PostForm.Get("comment"))
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating < 1 || rating > 5 {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		rest, restErr := h.App.Restaurants.FindByID(ctx, oid)
		if restErr != nil {
			h.clientError(w, http.StatusNotFound)
			return
		}
		h.render(w, r, "restaurant_view.tmpl", &templateData{
			User: u,
			Form: map[string]string{"error": "Rating must be 1 to 5"},
			Data: h.buildRestaurantViewData(ctx, rest, oid),
		})
		return
	}

	review := &models.Review{
		RestaurantID: oid,
		UserID:       u.ID,
		Rating:       rating,
		Comment:      comment,
		IsHidden:     false,
		CreatedAt:    time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.Reviews.Insert(ctx, review); err != nil {
		h.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/restaurants/"+parts[1], http.StatusSeeOther)
}

func (h *Handler) buildRestaurantViewData(ctx context.Context, rest *models.Restaurant, oid primitive.ObjectID) map[string]any {
	menu, err := h.App.MenuItems.ListByRestaurant(ctx, oid)
	if err != nil {
		menu = []models.MenuItem{}
	}

	reviews, err := h.App.Reviews.ListByRestaurant(ctx, oid)
	if err != nil {
		reviews = []models.Review{}
	}

	return map[string]any{
		"restaurant": rest,
		"menu":       menu,
		"reviews":    reviews,
	}
}

func (h *Handler) restaurantCreate(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "restaurant_create.tmpl", &templateData{
		User: h.currentUser(r),
		Data: map[string]any{
			"title": r.FormValue("title"),
		},
	})
}

func splitPath(path string) []string {
	out := []string{}
	cur := ""
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(path[i])
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
