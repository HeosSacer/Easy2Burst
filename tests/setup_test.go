package tests

import (
	"testing"
	"github.com/HeosSacer/Easy2Burst/internal"
	"reflect"
)

func TestCheckTools(t *testing.T){
	internal.CheckTools()
}

func TestNeedsJava(t *testing.T){
	result := internal.NeedsJava("1.8.0")
	AssertEqual(t, result, false)
	result = internal.NeedsJava("1.9.0")
	AssertEqual(t, result, true)
}

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
