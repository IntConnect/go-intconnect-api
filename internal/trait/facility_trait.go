package trait

type FacilityStatus string

const (
	FacilityStatusActive      FacilityStatus = "Active"
	FacilityStatusMaintenance FacilityStatus = "Maintenance"
	FacilityStatusArchived    FacilityStatus = "Archived"
)
