package handler

import (
	"effective_mobile/internal/dto/request"
	"effective_mobile/internal/dto/response"
	"effective_mobile/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

type SubscriptionHandler struct {
	subService *service.SubscriptionService
}

func NewSubscriptionHandler(subService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: subService}
}

// CreateSubscription godoc
// @Summary Create new subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body request.SubCreateRequest true "Subscription data"
// @Success 201 {object} response.SubResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription/ [post]
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub request.SubCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if createdSub, err := h.subService.CreateSub(&sub); err != nil {
		http.Error(w, "Failed to create subscription: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdSub)
	}
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Get subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" format(uuid)
// @Success 200 {object} response.SubResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscription/{id} [get]
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	subID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.subService.GetSub(subID)
	if err != nil || subscription == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// UpdateSubscription godoc
// @Summary Update subscription by ID
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" format(uuid)
// @Param request body request.SubUpdateRequest true "Subscription update request data"
// @Success 200 {object} response.SubResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	subID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	var subUpdate request.SubUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&subUpdate); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if updatedSub, err := h.subService.UpdateSub(subID, &subUpdate); err != nil {
		http.Error(w, "Failed to update subscription: "+err.Error(), http.StatusInternalServerError)
	} else {
		if updatedSub == nil {
			http.Error(w, "Subscription not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedSub)
	}
}

// DeleteSubscription godoc
// @Summary Delete subscription by ID
// @Description Delete subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" format(uuid)
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	subID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	if err := h.subService.DeleteSub(subID); err != nil {
		http.Error(w, "Failed to delete subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListSubscription godoc
// @Summary List subscription
// @Description Pageable subscription list
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} []response.SubResponse
// @Failure 500 {object} map[string]string
// @Router /subscription/ [get]
func (h *SubscriptionHandler) ListSubscription(w http.ResponseWriter, r *http.Request) {
	pageNumber := 0

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			pageNumber = p
		}
	}

	subscriptions, err := h.subService.GetSubsList(pageNumber)
	if err != nil {
		http.Error(w, "Failed to load subscriptions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// GetSubscriptionsSum godoc
// @Summary Sum subscriptions price
// @Description Sum subscriptions price by specified period
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param start_period query string true "Start period date" format(date)
// @Param end_period query string true "End period date" format(date)
// @Param userId query string false "User ID filter" format(uuid)
// @Param service_title query string false "Service title filter"
// @Success 200 {object} response.SubSumResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscription/sum [get]
func (h *SubscriptionHandler) GetSubscriptionsSum(w http.ResponseWriter, r *http.Request) {
	var filter request.SubFilterRequest

	if err := schema.NewDecoder().Decode(&filter, r.URL.Query()); err != nil {
		http.Error(w, "Invalid filters: "+err.Error(), http.StatusBadRequest)
		return
	}

	sum, err := h.subService.GetSubsSum(&filter)
	if err != nil {
		http.Error(w, "Failed to calculate subscriptions price sum: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response.SubSumResponse{Result: sum})
}
