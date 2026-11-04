package todotxt

type Entry struct {
	Title       string
	Description string
	Priority    string
	Due         string
	EstimateMin int64
	SubjectCode string
	ProjectCode string
	Extras      string
}

func (e Entry) hasClear() bool {
	return e.Priority == "-" || e.Due == "-" ||
		e.EstimateMin == -1 || e.SubjectCode == "-" || e.ProjectCode == "-" || e.Description == "-"
}

func (e Entry) isEmpty() bool {
	return e.Title == "" && e.Priority == "" && e.Due == "" &&
		e.EstimateMin == 0 && e.SubjectCode == "" && e.ProjectCode == "" && e.Extras == "" && e.Description == ""
}
