package forms

type AddLogJobForm struct {
	JobName   string `form:"jobName" binding:"required"`
	Topic     string `form:"topic" binding:"required"`
	IndexType string `form:"indexType" binding:"required"`
	Pipeline  string `form:"pipeline" binding:"required"`
}
