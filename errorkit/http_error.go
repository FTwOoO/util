package errorkit

type httpCodeErrorImp struct {
	*StructuredError
	httpCode int
}

func NewHttpError(httpCode int) HttpError {
	return &httpCodeErrorImp{
		StructuredError: NewStructuredError().(*StructuredError),
		httpCode:        httpCode,
	}
}

func (this *httpCodeErrorImp) GetHttpCode() int {
	return this.httpCode
}
