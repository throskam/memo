package handlers

import "net/http"

// Problem represents an RFC 9457 problem detail.
type Problem struct {
	Type     string
	Title    string
	Status   int
	Detail   string
	Instance string
	Err      error
}

func (p Problem) Error() string {
	if p.Err == nil {
		return ""
	}

	return p.Err.Error()
}

const (
	TypeAboutBlank = "about:blank"
	TypeInternal   = "/problems/internal"
	TypeForbidden  = "/problems/forbidden"
	TypeNotFound   = "/problems/not-found"
)

type ProblemOption func(*Problem)

func NewProblem(err error, opts ...ProblemOption) Problem {
	p := Problem{
		Type:   TypeAboutBlank,
		Status: http.StatusInternalServerError,
		Err:    err,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&p)
		}
	}

	if p.Type == TypeAboutBlank {
		p.Type = defaultTypeForStatus(p.Status)
	}

	return p
}

func defaultTypeForStatus(status int) string {
	switch status {
	case http.StatusInternalServerError:
		return TypeInternal
	case http.StatusForbidden:
		return TypeForbidden
	case http.StatusNotFound:
		return TypeNotFound
	default:
		return TypeAboutBlank
	}
}

func WithStatus(status int) ProblemOption {
	return func(p *Problem) {
		p.Status = status
	}
}

func WithType(value string) ProblemOption {
	return func(p *Problem) {
		p.Type = value
	}
}

func WithInternalType() ProblemOption {
	return WithType(TypeInternal)
}

func WithForbiddenType() ProblemOption {
	return WithType(TypeForbidden)
}

func WithNotFoundType() ProblemOption {
	return WithType(TypeNotFound)
}

func WithTitle(title string) ProblemOption {
	return func(p *Problem) {
		p.Title = title
	}
}

func WithDetail(detail string) ProblemOption {
	return func(p *Problem) {
		p.Detail = detail
	}
}

func WithInstance(instance string) ProblemOption {
	return func(p *Problem) {
		p.Instance = instance
	}
}
