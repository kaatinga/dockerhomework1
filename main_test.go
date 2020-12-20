package main

import (
	"github.com/julienschmidt/httprouter"
	"reflect"
	"testing"
)

func Test_getPhraseBytes(t *testing.T) {

	tests := []struct {
		name string
		ps httprouter.Params
		want []byte
	}{
		{ "ok", httprouter.Params{httprouter.Param{
			Key:   "phrase",
			Value: "world!",
		}}, []byte("Hello, world!")},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPhraseBytes(tt.ps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPhraseBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
