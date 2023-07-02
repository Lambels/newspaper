package interpreter

type ErrorReporter struct {}

func (r *ErrorReporter) AddErr(err error) {}

func (r *ErrorReporter) Flush() error {}

func (r *ErrorReporter) HadError() bool {}

type Context struct {}
