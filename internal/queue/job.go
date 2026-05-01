package queue

type JobType string
type JobStatus string

const (
	JobUpload   JobType = "UPLOAD"
	JobDownload JobType = "DOWNLOAD"

	StatusPending    JobStatus = "PENDING"
	StatusProcessing JobStatus = "PROCESSING"
	StatusDone       JobStatus = "DONE"
	StatusFailed     JobStatus = "FAILED"
)

type Job struct {
	ID         string
	Type       JobType
	Status     JobStatus
	Payload    string
	RetryCount int
	CreatedAt  int64
}
