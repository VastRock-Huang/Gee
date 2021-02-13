package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func handlerName (ctx *Context) {

}

func handlerB(c *Context){

}

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name",handlerName)
	r.addRoute("GET", "/hello/:file", handlerB)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func Test_parsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

//func Test_parsePattern(t *testing.T) {
//	type args struct {
//		pattern string
//	}
//	tests := []struct {
//		name string
//		args args
//		want []string
//	}{
//		// TODO: Add test cases.
//		{
//			"1",
//			args{"/p/:name/*"},
//			[]string{"p", ":name", "*"},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := parsePattern(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("parsePattern() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestGetRouter(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/geektutu")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	fmt.Printf("%v", n)
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}

