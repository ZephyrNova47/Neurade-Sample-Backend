package model

import "time"

type PrCreateRequest struct {
	CourseID      int       `json:"course_id"`
	AssignmentID  int       `json:"assignment_id"`
	PrName        string    `json:"pr_name"`
	PrDescription string    `json:"pr_description"`
	Status        string    `json:"status"`
	PrNumber      int       `json:"pr_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PrResponse struct {
	ID            int       `json:"id"`
	CourseID      int       `json:"course_id"`
	AssignmentID  int       `json:"assignment_id"`
	PrName        string    `json:"pr_name"`
	PrDescription string    `json:"pr_description"`
	Status        string    `json:"status"`
	PrNumber      int       `json:"pr_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
