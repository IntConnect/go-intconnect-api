package entity

type Machine struct {
	Id           uint64   `gorm:"column:id;primaryKey;autoIncrement"`
	FacilityId   uint64   `gorm:"column:facility_id;"`
	Name         string   `gorm:"column:name"`
	Code         string   `gorm:"column:code"`
	Description  string   `gorm:"column:description"`
	ModelPath    string   `json:"model_path"`
	ModelOffsetX float32  `json:"model_offset_x"`
	ModelOffsetY float32  `json:"model_offset_y"`
	ModelOffsetZ float32  `json:"model_offset_z"`
	ModelScale   float32  `json:"model_scale"`
	Facility     Facility `gorm:"foreignKey:FacilityId;references:Id"`
	Auditable
}
