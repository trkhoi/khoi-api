package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestSubFromViper(t *testing.T) {
	v := viper.New()
	v.Set("a", "a")
	v.Set("b", "b")
	v.Set("c", "c")
	v.Set("a.a", "a.a")
	v.Set("a.b", "a.b")
	av := Sub(v, "a")
	if av == nil {
		t.Fatalf("Sub(%v, 'a') => %v", v, av)
	}

	if av.GetString("a") != "a.a" {
		t.Errorf("av.GetString('a') = %s, expected 'a.a'", av.GetString("a"))
	}
	if av.GetString("a.a") != "" {
		t.Errorf("av.GetString('a.a') = %s, expected ''", av.GetString("a.a"))
	}
	if av.GetString("b") != "a.b" {
		t.Errorf("av.GetString('b') = %s, expected 'a.b'", av.GetString("b"))
	}
	if av.GetString("c") != "" {
		t.Errorf("av.GetString('c') = %s, expected ''", av.GetString(""))
	}
}
