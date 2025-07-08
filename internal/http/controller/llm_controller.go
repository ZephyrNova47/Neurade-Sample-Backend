package controller

import (
	"be/neurade/v2/internal/model/converter"
	"be/neurade/v2/internal/service"
	"be/neurade/v2/internal/util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type LLMController struct {
	LLMService *service.LLMService
	log        *logrus.Logger
}

func NewLLMController(service *service.LLMService, logger *logrus.Logger) *LLMController {
	return &LLMController{
		LLMService: service,
		log:        logger,
	}
}

func (c *LLMController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.log.Println("Failed to parse form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	request := converter.RequestToLLMRequest(r)
	if request.ModelName == "" || request.ModelToken == "" {
		c.log.Println("Missing required fields")
		http.Error(w, "Model name and token are required", http.StatusBadRequest)
		return
	}
	if (util.IsValidLLM(request.ModelName, request.ModelToken)) == false {
		request.Status = "invalid"
	} else {
		request.Status = "active"
	}
	llm, err := c.LLMService.Create(r.Context(), request)
	if err != nil {
		c.log.Println("Failed to create LLM:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(llm)
}

func (c *LLMController) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.log.Println("Invalid ID:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	llm, err := c.LLMService.GetByID(r.Context(), id)
	if err != nil {
		c.log.Println("Failed to get LLM by ID:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(llm)
}
