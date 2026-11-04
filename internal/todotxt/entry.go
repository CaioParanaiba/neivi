package todotxt

type Entry struct {
	Title       string
	Priority    string
	Due         string
	EstimateMin int64
	SubjectCode string
	ProjectCode string
	Extras      string
}

func (e Entry) hasClear() bool {
	return e.Priority == "-" || e.Due == "-" ||
		e.EstimateMin == -1 || e.SubjectCode == "-" || e.ProjectCode == "-"
}
