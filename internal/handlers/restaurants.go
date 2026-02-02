package handlers

import (
	"context"
	"net/http"
	"time"

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

	menu, err := h.App.MenuItems.ListByRestaurant(ctx, oid)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.render(w, r, "restaurant_view.tmpl", &templateData{
		User: h.currentUser(r),
		Data: map[string]any{
			"restaurant": rest,
			"menu":       menu,
		},
	})
}

func splitPath(path string) []string {
	// "/restaurants/xxx" => ["restaurants", "xxx"]
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
