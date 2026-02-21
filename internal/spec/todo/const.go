package todo

// Query and form param names for Todo API.
const (
	ParamIncludeCompleted = "include_completed"
	ParamSort             = "sort"
	ParamLimit            = "limit"
	ParamOffset           = "offset"
)

// Validation and pagination limits.
const (
	LimitDefault = 50
	LimitMax     = 100
)

// Allowed sort values.
const (
	SortDueDateAsc   = "due_date_asc"
	SortDueDateDesc  = "due_date_desc"
	SortCreatedAsc   = "created_at_asc"
	SortCreatedDesc  = "created_at_desc"
)
