package result

import "net/http"

type Result struct {
	isOk   bool
	status int
	ok     any
	err    error
}

func Ok(ok any) Result {
	return Result{
		isOk:   true,
		status: http.StatusOK,
		ok:     ok,
		err:    nil,
	}
}

func Err(status int, err error) Result {
	return Result{
		isOk:   false,
		status: status,
		ok:     nil,
		err:    err,
	}
}

func (r Result) Status() int {
	return r.status
}

func (r Result) Ok() (any, bool) {
	return r.ok, r.isOk
}

func (r Result) Err() (error, bool) {
	return r.err, !r.isOk
}
