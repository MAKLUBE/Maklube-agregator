package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) partnerIncomingOrders(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rests, err := h.App.Restaurants.ListByOwner(ctx, u.ID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	groups := make([]map[string]any, 0, len(rests))
	for _, rest := range rests {
		orders, err := h.App.Orders.ListByRestaurant(ctx, rest.ID)
		if err != nil {
			h.serverError(w, err)
			return
		}
		groups = append(groups, map[string]any{
			"restaurant": rest,
			"orders":     orders,
		})
	}

	h.render(w, r, "partner_orders.tmpl", &templateData{
		User: u,
		Data: groups,
	})
}

func (h *Handler) partnerRestaurants(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.App.Restaurants.ListByOwner(ctx, u.ID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.render(w, r, "partner_restaurants.tmpl", &templateData{
		User: u,
		Data: items,
	})
}

func (h *Handler) partnerRestaurantEditForm(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	h.render(w, r, "partner_restaurant_edit.tmpl", &templateData{
		User: u,
		Data: rest,
	})
}

func (h *Handler) partnerRestaurantEditPost(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	rest.Name = strings.TrimSpace(r.PostForm.Get("name"))
	rest.Description = strings.TrimSpace(r.PostForm.Get("description"))
	rest.Phone = strings.TrimSpace(r.PostForm.Get("phone"))
	rest.Address.AddressText = strings.TrimSpace(r.PostForm.Get("address_text"))
	rest.Address.City = strings.TrimSpace(r.PostForm.Get("city"))
	rest.Address.District = strings.TrimSpace(r.PostForm.Get("district"))
	rest.Address.MapURL = strings.TrimSpace(r.PostForm.Get("map_url"))
	rest.Instagram = strings.TrimSpace(r.PostForm.Get("instagram"))

	if rest.Name == "" || rest.Address.AddressText == "" || rest.Address.City == "" {
		h.render(w, r, "partner_restaurant_edit.tmpl", &templateData{
			User: u,
			Data: rest,
			Form: map[string]string{"error": "Fill required fields"},
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.Restaurants.Update(ctx, rest); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/partner/restaurants", http.StatusSeeOther)
}

func (h *Handler) partnerMenuNewForm(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	h.render(w, r, "partner_menu_form.tmpl", &templateData{
		User: u,
		Data: map[string]any{
			"restaurant": rest,
		},
	})
}

func (h *Handler) partnerHalalRequestForm(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	parts := splitPath(r.URL.Path)
	if len(parts) < 4 {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	restaurantID, err := primitive.ObjectIDFromHex(parts[3])
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	rest, err := h.App.Restaurants.FindByID(ctx, restaurantID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}
	if rest.OwnerUserID != u.ID {
		h.clientError(w, http.StatusForbidden)
		return
	}
	h.render(w, r, "partner_halal_request.tmpl", &templateData{
		User: u,
		Data: map[string]any{
			"restaurant": rest,
		},
	})
}

func (h *Handler) partnerHalalRequestPost(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	var ingredientProofs []models.IngredientProof

	maxIndex := -1
	for key := range r.PostForm {
		if strings.HasPrefix(key, "products[") && strings.Contains(key, "][name]") {
			start := strings.Index(key, "[") + 1
			end := strings.Index(key, "]")
			if start > 0 && end > start {
				if idx, err := strconv.Atoi(key[start:end]); err == nil {
					if idx > maxIndex {
						maxIndex = idx
					}
				}
			}
		}
	}

	if maxIndex < 0 {
		h.render(w, r, "partner_halal_request.tmpl", &templateData{
			User: u,
			Data: map[string]any{"restaurant": rest},
			Form: map[string]string{"error": "Fill required fields"},
		})
		return
	}

	hasValidProduct := false
	for i := 0; i <= maxIndex; i++ {
		nameKey := "products[" + strconv.Itoa(i) + "][name]"
		proofTypeKey := "products[" + strconv.Itoa(i) + "][proof_type]"
		urlKey := "products[" + strconv.Itoa(i) + "][url]"

		name := strings.TrimSpace(r.PostForm.Get(nameKey))
		proofType := strings.TrimSpace(r.PostForm.Get(proofTypeKey))
		url := strings.TrimSpace(r.PostForm.Get(urlKey))

		if name == "" || url == "" {
			continue
		}

		if proofType == "" {
			proofType = "certificate"
		}

		hasValidProduct = true
		ingredientProofs = append(ingredientProofs, models.IngredientProof{
			Ingredient: name,
			ProofType:  proofType,
			ProofURLs:  []string{url},
		})
	}

	if !hasValidProduct {
		h.render(w, r, "partner_halal_request.tmpl", &templateData{
			User: u,
			Data: map[string]any{"restaurant": rest},
			Form: map[string]string{"error": "Fill required fields"},
		})
		return
	}

	req := &models.HalalVerification{
		ID:               primitive.NewObjectID(),
		RestaurantID:     rest.ID,
		RequestedBy:      u.ID,
		Status:           "pending",
		IngredientProofs: ingredientProofs,
		CreatedAt:        time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.Halal.Insert(ctx, req); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/partner/restaurants", http.StatusSeeOther)
}

func (h *Handler) partnerMenuList(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.App.MenuItems.ListByRestaurant(ctx, rest.ID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.render(w, r, "partner_menu_list.tmpl", &templateData{
		User: u,
		Data: map[string]any{
			"restaurant": rest,
			"items":      items,
		},
	})
}

func (h *Handler) partnerMenuNewPost(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	item, err := h.menuItemFromForm(r, rest.ID)
	if err != nil {
		h.render(w, r, "partner_menu_form.tmpl", &templateData{
			User: u,
			Data: map[string]any{"restaurant": rest},
			Form: map[string]string{"error": "Fill required fields"},
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.MenuItems.Insert(ctx, item); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/partner/restaurants", http.StatusSeeOther)
}

func (h *Handler) partnerMenuEditForm(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	item, err := h.getPartnerMenuItem(r, rest.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	h.render(w, r, "partner_menu_form.tmpl", &templateData{
		User: u,
		Data: map[string]any{
			"restaurant": rest,
			"item":       item,
		},
	})
}

func (h *Handler) partnerMenuEditPost(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rest, err := h.getPartnerRestaurant(r, u.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	item, err := h.getPartnerMenuItem(r, rest.ID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	updated, err := h.menuItemFromForm(r, rest.ID)
	if err != nil {
		h.render(w, r, "partner_menu_form.tmpl", &templateData{
			User: u,
			Data: map[string]any{"restaurant": rest, "item": item},
			Form: map[string]string{"error": "Fill required fields"},
		})
		return
	}

	item.Name = updated.Name
	item.Description = updated.Description
	item.Category = updated.Category
	item.Price = updated.Price
	item.PrepTimeMin = updated.PrepTimeMin
	item.IsAvailable = updated.IsAvailable
	item.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.MenuItems.Update(ctx, item); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/partner/restaurants", http.StatusSeeOther)
}

func (h *Handler) getPartnerRestaurant(r *http.Request, ownerID primitive.ObjectID) (*models.Restaurant, error) {
	parts := splitPath(r.URL.Path)
	if len(parts) < 4 {
		return nil, errors.New("invalid path")
	}

	oid, err := primitive.ObjectIDFromHex(parts[3])
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rest, err := h.App.Restaurants.FindByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if rest.OwnerUserID != ownerID {
		return nil, errors.New("forbidden")
	}
	return rest, nil
}

func (h *Handler) getPartnerMenuItem(r *http.Request, restaurantID primitive.ObjectID) (*models.MenuItem, error) {
	parts := splitPath(r.URL.Path)
	if len(parts) < 6 {
		return nil, errors.New("invalid path")
	}

	itemID, err := primitive.ObjectIDFromHex(parts[5])
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	item, err := h.App.MenuItems.FindByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item.RestaurantID != restaurantID {
		return nil, errors.New("forbidden")
	}
	return item, nil
}

func (h *Handler) menuItemFromForm(r *http.Request, restaurantID primitive.ObjectID) (*models.MenuItem, error) {
	name := strings.TrimSpace(r.PostForm.Get("name"))
	description := strings.TrimSpace(r.PostForm.Get("description"))
	category := strings.TrimSpace(r.PostForm.Get("category"))
	priceStr := strings.TrimSpace(r.PostForm.Get("price"))
	prepStr := strings.TrimSpace(r.PostForm.Get("prep_time_min"))
	availableStr := strings.TrimSpace(r.PostForm.Get("is_available"))

	if name == "" || priceStr == "" {
		return nil, errors.New("missing fields")
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, err
	}

	prepTime := 0
	if prepStr != "" {
		prepTime, _ = strconv.Atoi(prepStr)
	}

	isAvailable := availableStr == "on"

	item := &models.MenuItem{
		RestaurantID: restaurantID,
		Name:         name,
		Description:  description,
		Category:     category,
		Price:        price,
		PrepTimeMin:  prepTime,
		IsAvailable:  isAvailable,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	return item, nil
}

func parseIngredientProofs(raw string) ([]models.IngredientProof, error) {
	if strings.TrimSpace(raw) == "" {
		return []models.IngredientProof{}, nil
	}
	lines := strings.Split(raw, "\n")
	out := make([]models.IngredientProof, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			return nil, errors.New("invalid ingredient proof line")
		}
		ingredient := strings.TrimSpace(parts[0])
		proofType := strings.TrimSpace(parts[1])
		urlParts := strings.Split(parts[2], ",")
		urls := make([]string, 0, len(urlParts))
		for _, p := range urlParts {
			v := strings.TrimSpace(p)
			if v != "" {
				urls = append(urls, v)
			}
		}
		if ingredient == "" || proofType == "" || len(urls) == 0 {
			return nil, errors.New("invalid ingredient proof values")
		}
		out = append(out, models.IngredientProof{Ingredient: ingredient, ProofType: proofType, ProofURLs: urls})
	}
	return out, nil
}

func (h *Handler) partnerOrderUpdateStatus(w http.ResponseWriter, r *http.Request) {
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

	orderID, err := primitive.ObjectIDFromHex(parts[2])
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}
	status := strings.TrimSpace(r.PostForm.Get("status"))
	if status == "" {
		http.Redirect(w, r, "/partner/orders", http.StatusSeeOther)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	order, err := h.App.Orders.FindByID(ctx, orderID)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	rest, err := h.App.Restaurants.FindByID(ctx, order.RestaurantID)
	if err != nil || rest.OwnerUserID != u.ID {
		h.clientError(w, http.StatusForbidden)
		return
	}

	if err := h.App.Orders.UpdateStatus(ctx, orderID, status); err != nil {
		h.serverError(w, err)
		return
	}

	statusRow := &models.OrderStatusHistory{
		OrderID:   orderID,
		Status:    status,
		ChangedBy: u.ID,
		Comment:   "partner update",
		CreatedAt: time.Now().UTC(),
	}
	_ = h.App.OrderStatus.Insert(ctx, statusRow)

	http.Redirect(w, r, "/partner/orders", http.StatusSeeOther)
}
