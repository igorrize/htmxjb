package domain


type JobSource int

const (
	Indeed JobSource = iota
	LinkedIn
	Csv
)

func (js JobSource) String() string {
	switch js {
	case Indeed:
		return "indeed"
	case LinkedIn:
		return "linkedin"
	case Csv:
		return "csv"
	default:
		return "unknown"
	}
}

type JobType int

const (
	Remote JobType = iota
	Onsite
	Hybrid
)

func (jt JobType) String() string {
	switch jt {
	case Remote:
		return "remote"
	case Onsite:
		return "onsite"
	case Hybrid:
		return "hybrid"
	default:
		return "unknown"
	}
}

type Job struct {
	ID          int64
	ExternalID  string
	Title       string
	Description string
	Type        JobType
	Source      JobSource
}

// ID          int       `json:"id"`
//
// Title       string    `json:"title"`
// Description string    `json:"description,omitempty"`
// CreatedAt   time.Time `json:"created_at,omitempty"`
// UpdatedAt   time.Time `json:"updated_at,omitempty"`
// Type        string    `json:"type"`
// Location    string    `json:"location"`
// Salary      string    `json:"salary"`
// IsNew       bool      `json:"is_new"`
// Company     string    `json:"company"`
// Tags        []string  `json:"tags,omitempty"`
