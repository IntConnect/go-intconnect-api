package entity

type BreakdownResource struct {
	Id          uint64 `gorm:"column:id,primaryKey;autoIncrement"`
	BreakdownId uint64 `gorm:"column:breakdown_id"`
	ImagePath   string `gorm:"column:image_path"`
	VideoPath   string `gorm:"column:video_path"`
}
