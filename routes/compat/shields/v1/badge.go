package v1

import (
	"github.com/gorilla/mux"
	config2 "github.com/muety/wakapi/config"
	"github.com/muety/wakapi/models"
	v1 "github.com/muety/wakapi/models/compat/shields/v1"
	"github.com/muety/wakapi/services"
	"github.com/muety/wakapi/utils"
	"net/http"
	"regexp"
	"strings"
)

const (
	intervalPattern     = `interval:([a-z0-9_]+)`
	entityFilterPattern = `(project|os|editor|language|machine):([_a-zA-Z0-9-]+)`
)

type BadgeHandler struct {
	userSrvc    *services.UserService
	summarySrvc *services.SummaryService
	config      *config2.Config
}

func NewBadgeHandler(summaryService *services.SummaryService, userService *services.UserService) *BadgeHandler {
	return &BadgeHandler{
		summarySrvc: summaryService,
		userSrvc:    userService,
		config:      config2.Get(),
	}
}

func (h *BadgeHandler) ApiGet(w http.ResponseWriter, r *http.Request) {
	intervalReg := regexp.MustCompile(intervalPattern)
	entityFilterReg := regexp.MustCompile(entityFilterPattern)

	if userAgent := r.Header.Get("user-agent"); !strings.HasPrefix(userAgent, "Shields.io/") && !h.config.IsDev() {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	requestedUserId := mux.Vars(r)["user"]
	user, err := h.userSrvc.GetUserById(requestedUserId)
	if err != nil || !user.BadgesEnabled {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var filterEntity, filterKey string
	if groups := entityFilterReg.FindStringSubmatch(r.URL.Path); len(groups) > 2 {
		filterEntity, filterKey = groups[1], groups[2]
	}

	var interval = models.IntervalPast30Days
	if groups := intervalReg.FindStringSubmatch(r.URL.Path); len(groups) > 1 {
		interval = groups[1]
	}

	filters := &models.Filters{}
	switch filterEntity {
	case "project":
		filters.Project = filterKey
	case "os":
		filters.OS = filterKey
	case "editor":
		filters.Editor = filterKey
	case "language":
		filters.Language = filterKey
	case "machine":
		filters.Machine = filterKey
	}

	summary, err, status := h.loadUserSummary(user, interval)
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	vm := v1.NewBadgeDataFrom(summary, filters)
	utils.RespondJSON(w, http.StatusOK, vm)
}

func (h *BadgeHandler) loadUserSummary(user *models.User, interval string) (*models.Summary, error, int) {
	err, from, to := utils.ResolveInterval(interval)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	summaryParams := &models.SummaryParams{
		From: from,
		To:   to,
		User: user,
	}

	summary, err := h.summarySrvc.Construct(summaryParams.From, summaryParams.To, summaryParams.User, summaryParams.Recompute)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return summary, nil, http.StatusOK
}