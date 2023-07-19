package types

type Console struct {
	// Name    string // ?
	IP int
	// Region  string
	Content string
	// Pointer int
}

type DiagnosticCode int

const ( // From official Microsoft Documentation
	Error   DiagnosticCode = 1
	Warning                = 2
	Info                   = 3
	Hint                   = 4
)

type Range struct {
	Start int
	End   int
}

type Diagnostic struct {
	Console Console
	Code    DiagnosticCode
	Message string
	Range   Range
}

type CompletionRequest struct {
	Console Console
	Pointer int
}
