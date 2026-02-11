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

func (h *Handler) customerOrdersList(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	orders, err := h.App.Orders.ListByCustomer(ctx, u.ID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.render(w, r, "customer_orders.tmpl", &templateData{
		User: u,
		Data: orders,
	})
}

func (h *Handler) orderCreate(w http.ResponseWriter, r *http.Request) {
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

	restID, err := primitive.ObjectIDFromHex(parts[1])
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	menuIDHex := strings.TrimSpace(r.PostForm.Get("menu_item_id"))
	qtyStr := strings.TrimSpace(r.PostForm.Get("qty"))
	deliveryType := strings.TrimSpace(r.PostForm.Get("delivery_type"))
	deliveryAddr := strings.TrimSpace(r.PostForm.Get("delivery_address"))
	paymentType := strings.TrimSpace(r.PostForm.Get("payment_type"))
	customerNote := strings.TrimSpace(r.PostForm.Get("customer_comment"))
	if menuIDHex == "" || qtyStr == "" {
		h.renderOrderFormError(w, r, restID, u, "Select item and quantity")
		return
	}

	menuID, err := primitive.ObjectIDFromHex(menuIDHex)
	if err != nil {
		h.renderOrderFormError(w, r, restID, u, "Invalid menu item")
		return
	}

	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty < 1 {
		h.renderOrderFormError(w, r, restID, u, "Quantity must be 1 or more")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rest, err := h.App.Restaurants.FindByID(ctx, restID)
	if err != nil {
		h.renderOrderFormError(w, r, restID, u, "Restaurant not found")
		return
	}

	if paymentType == "kaspi_transfer" && strings.TrimSpace(rest.KaspiNumber) == "" {
		h.renderOrderFormError(w, r, restID, u, "Kaspi transfer is not available for this")
		return
	}

	menu, err := h.App.MenuItems.FindByID(ctx, menuID)
	if err != nil || menu.RestaurantID != restID {
		h.renderOrderFormError(w, r, restID, u, "Menu item not found")
		return
	}

	lineTotal := menu.Price * float64(qty)
	order := &models.Order{
		ID:           primitive.NewObjectID(),
		CustomerID:   u.ID,
		RestaurantID: restID,
		OrderStatus:  "new",
		DeliveryType: deliveryType,
		DeliveryAddr: deliveryAddr,
		Subtotal:     lineTotal,
		Total:        lineTotal,
		DeliveryFee:  0,
		PaymentType:  paymentType,
		CustomerNote: customerNote,
		Items: []models.OrderItem{
			{
				MenuItemID: menu.ID,
				Name:       menu.Name,
				Price:      menu.Price,
				Qty:        qty,
				LineTotal:  lineTotal,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := h.App.Orders.Insert(ctx, order); err != nil {
		h.serverError(w, err)
		return
	}

	status := &models.OrderStatusHistory{
		OrderID:   order.ID,
		Status:    "new",
		ChangedBy: u.ID,
		Comment:   "order created",
		CreatedAt: time.Now().UTC(),
	}
	_ = h.App.OrderStatus.Insert(ctx, status)

	http.Redirect(w, r, "/orders", http.StatusSeeOther)
}

func (h *Handler) renderOrderFormError(w http.ResponseWriter, r *http.Request, restID primitive.ObjectID, u *models.User, msg string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rest, restErr := h.App.Restaurants.FindByID(ctx, restID)
	if restErr != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	h.render(w, r, "restaurant_view.tmpl", &templateData{
		User: u,
		Form: map[string]string{"order_error": msg},
		Data: h.buildRestaurantViewData(ctx, rest, restID),
	})
}
