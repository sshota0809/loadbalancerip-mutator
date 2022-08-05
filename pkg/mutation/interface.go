package mutation

import "net/http"

type Mutator interface {
	GenerateHandler() (http.Handler, error)
}
