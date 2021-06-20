package webapi

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/matejgrzinic/portfolio/appcontext"
)

func TestApiTimeline(t *testing.T) {
	type args struct {
		appcontext *appcontext.AppContext
	}
	tests := []struct {
		name string
		args args
		want func(http.ResponseWriter, *http.Request)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApiTimeline(tt.args.appcontext); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApiTimeline() = %v, want %v", got, tt.want)
			}
		})
	}
}
