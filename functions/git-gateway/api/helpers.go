package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/netlify/git-gateway/conf"
	"github.com/netlify/git-gateway/models"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

func addRequestID(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	id := uuid.NewRandom().String()
	ctx := r.Context()
	ctx = withRequestID(ctx, id)
	return ctx, nil
}

func sanitizeOutput(obj interface{}) interface{} {
	switch v := obj.(type) {
	case InstanceResponse:
		v.Instance.BaseConfig.GitHub.AccessToken = ""
	case *InstanceResponse:
		v.Instance.BaseConfig.GitHub.AccessToken = ""
	case models.Instance:
		v.BaseConfig.GitHub.AccessToken = ""
	case *models.Instance:
		v.BaseConfig.GitHub.AccessToken = ""
	case *conf.Configuration:
		v.GitHub.AccessToken = ""
	case conf.Configuration:
		v.GitHub.AccessToken = ""
		// must return here because v != obj due to value copying
		return v
	default:
	}
	return obj
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	obj = sanitizeOutput(obj)

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

// From https://golang.org/src/net/http/httputil/reverseproxy.go?s=2298:2359#L72
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
