package model

import "time"

type CreateBreakdownRequest struct {
	MachineId             uint64 `json:"machine_id"`
	ReportedBy            uint64 `json:"reported_by"`
	VerifiedBy            uint64 `json:"verified_by"`
	ProblemIdentification string `json:"problem_identification"`
	PeopleIssue           string `json:"people_issue"`
	InspectionIssue       string `json:"inspection_issue"`
	ConditionIssue        string `json:"condition_issue"`
	ActionTaken           string `json:"action_taken"`
	PartsIssue            string `json:"parts_issue"`
	AnalysisNotes         string `json:"analysis_notes"`
	CorrectiveAction      string `json:"corrective_action"`
	PreventiveAction      string `json:"preventive_action"`
	Status                string `json:"status"`
	ProblemAt             string `json:"problem_at"`
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
