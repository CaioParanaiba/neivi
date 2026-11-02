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
