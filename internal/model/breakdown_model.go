package model

import (
	"mime/multipart"
	"time"
)

type CreateBreakdownRequest struct {
	MachineId                 uint64                     `json:"machine_id" validate:"required,number,gte=1,exists=machines;id"`
	ProblemIdentification     string                     `json:"problem_identification" validate:"required,min=2"`
	PeopleIssue               string                     `json:"people_issue" validate:"required,min=2"`
	InspectionIssue           string                     `json:"inspection_issue" validate:"required,min=2"`
	ConditionIssue            string                     `json:"condition_issue" validate:"required,min=2"`
	ActionTaken               string                     `json:"action_taken" validate:"required,min=2"`
	PartsIssue                string                     `json:"parts_issue" validate:"required,min=2"`
	AnalysisNotes             string                     `json:"analysis_notes" validate:"required,min=2"`
	CorrectiveAction          string                     `json:"corrective_action" validate:"required,min=2"`
	PreventiveAction          string                     `json:"preventive_action" validate:"required,min=2"`
	ProblemAt                 string                     `json:"problem_at" validate:"required,datetime"`
	BreakdownResourceRequests []BreakdownResourceRequest `json:"breakdown_resources_requests"`
}

type BreakdownResourceRequest struct {
	ImageFile *multipart.FileHeader `json:"image_file" validate:"omitempty,fileExtension=.png .jpg"`
	VideoFile *multipart.FileHeader `json:"video_file" validate:"omitempty,fileExtension=.mp4"`
}

type UpdateBreakdownRequest struct {
	Id                    uint64    `json:"id"`
	MachineId             uint64    `json:"machine_id"`
	ReportedBy            uint64    `json:"reported_by"`
	VerifiedBy            uint64    `json:"verified_by"`
	ProblemIdentification string    `json:"problem_identification"`
	PeopleIssue           string    `json:"people_issue"`
	InspectionIssue       string    `json:"inspection_issue"`
	ConditionIssue        string    `json:"condition_issue"`
	ActionTaken           string    `json:"action_taken"`
	PartsIssue            string    `json:"parts_issue"`
	AnalysisNotes         string    `json:"analysis_notes"`
	CorrectiveAction      string    `json:"corrective_action"`
	PreventiveAction      string    `json:"preventive_action"`
	Status                string    `json:"status"`
	ProblemAt             time.Time `json:"problem_at"`
}

type BreakdownResponse struct {
	Id                    uint64             `json:"id;primaryKey;autoIncrement"`
	MachineId             uint64             `json:"machine_id"`
	ReportedBy            uint64             `json:"reported_by"`
	VerifiedBy            uint64             `json:"verified_by"`
	ProblemIdentification string             `json:"problem_identification"`
	PeopleIssue           string             `json:"people_issue"`
	InspectionIssue       string             `json:"inspection_issue"`
	ConditionIssue        string             `json:"condition_issue"`
	ActionTaken           string             `json:"action_taken"`
	PartsIssue            string             `json:"parts_issue"`
	AnalysisNotes         string             `json:"analysis_notes"`
	CorrectiveAction      string             `json:"corrective_action"`
	PreventiveAction      string             `json:"preventive_action"`
	Status                string             `json:"status"`
	StartTime             string             `json:"start_time"`
	EndTime               string             `json:"end_time"`
	ProblemAt             string             `json:"problem_at"`
	AuditableResponse     *AuditableResponse `json:"auditable" mapstructure:"auditable"`
}

func (breakdownResponse *BreakdownResponse) GetAuditableResponse() *AuditableResponse {
	return breakdownResponse.AuditableResponse
}

func (breakdownResponse *BreakdownResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	breakdownResponse.AuditableResponse = auditableResponse
}
