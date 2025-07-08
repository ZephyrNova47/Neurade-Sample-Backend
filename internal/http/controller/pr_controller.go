package controller

import (
	"be/neurade/v2/internal/model/converter"
	"be/neurade/v2/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type PrController struct {
	PrService *service.PrService
	Log       *logrus.Logger
}

func NewPrController(prService *service.PrService, log *logrus.Logger) *PrController {
	return &PrController{
		PrService: prService,
		Log:       log,
	}
}

func (c *PrController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.Log.Println("Failed to parse form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	pr := converter.RequestToPrRequest(r)

	prResponse, err := c.PrService.Create(r.Context(), pr)
	if err != nil {
		c.Log.Println("Failed to create pr:", err)
		http.Error(w, "Failed to create pr:", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prResponse)
}

func (c *PrController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "pr_id"))
	prResponse, err := c.PrService.GetByID(r.Context(), id)
	if err != nil {
		c.Log.Println("Failed to get pr by ID")
		http.Error(w, "Failed to get pr by ID", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prResponse)
}

func (c *PrController) GetAllByCourse(w http.ResponseWriter, r *http.Request) {
	courseID, _ := strconv.Atoi(chi.URLParam(r, "course_id"))
	prResponse, err := c.PrService.GetAllByCourse(r.Context(), courseID)
	if err != nil {
		c.Log.Println("Failed to get all pr by course")
		http.Error(w, "Failed to get all pr by course", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(prResponse)
}
