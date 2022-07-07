package util

type Report struct {
	OnError func(v interface{})
	OnWarn func(v interface{})
	OnSuccess func(v interface{})
}

func (r *Report) ReportError(err error) {
	r.OnError(err)
}
