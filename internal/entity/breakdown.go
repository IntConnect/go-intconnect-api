package entity

import "time"

type Breakdown struct {
	Id                    uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId             uint64    `gorm:"column:machine_id"`
	ReportedBy            uint64    `gorm:"column:reported_by"`
	VerifiedBy            uint64    `gorm:"column:verified_by"`
	ProblemIdentification string    `gorm:"column:problem_identification"`
	PeopleIssue           string    `gorm:"column:people_issue"`
	InspectionIssue       string    `gorm:"column:inspection_issue"`
	ConditionIssue        string    `gorm:"column:condition_issue"`
	ActionTaken           string    `gorm:"column:action_taken"`
	PartsIssue            string    `gorm:"column:parts_issue"`
	AnalysisNotes         string    `gorm:"column:analysis_notes"`
	CorrectiveAction      string    `gorm:"column:corrective_action"`
	PreventiveAction      string    `gorm:"column:preventive_action"`
	Status                string    `gorm:"column:status"`
	StartTime             time.Time `gorm:"column:start_time"`
	EndTime               time.Time `gorm:"column:end_time"`
	ProblemAt             time.Time `gorm:"column:problem_at"`
	Machine               *Machine  `gorm:"foreignKey:MachineId;references:Id"`
	ReportedByUser        *User     `gorm:"foreignKey:ReportedBy;references:Id"`
	VerifiedByUser        *User     `gorm:"foreignKey:VerifiedBy;references:Id"`
	Auditable             Auditable `gorm:"embedded"`
}

func (breakdownEntity *Breakdown) GetAuditable() *Auditable {
	return &breakdownEntity.Auditable
}
