package apploader

import (
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/tnborg/panel/pkg/types"
)

var apps sync.Map

type Loader struct{}

func (r *Loader) Add(app ...types.App) {
	for item := range slices.Values(app) {
		slug := getSlug(item)
		apps.Store(slug, item)
	}
}

func (r *Loader) Register(mux chi.Router) {
	/*for slug, item := range r.Apps {
		mux.Route("/"+slug, item.Route)
	}*/

	apps.Range(func(key, value any) bool {
		app := value.(types.App)
		mux.Route("/"+key.(string), app.Route)
		return true
	})
}

func Slugs() []string {
	var slugs []string
	apps.Range(func(key, value any) bool {
		slugs = append(slugs, key.(string))
		return true
	})
	return slugs
}

func getSlug(app types.App) string {
	if app == nil {
		return ""
	}

	t := reflect.TypeOf(app)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	pkgPath := t.PkgPath()
	if pkgPath == "" {
		return ""
	}

	parts := strings.Split(pkgPath, "/")
	return parts[len(parts)-1]
}
