package config

import "github.com/spf13/viper"

type View interface {
	IsSet(string) bool
	GetString(string) string
	GetInt(string) int
	GetInt64(string) int64
	GetFloat64(string) float64
	GetStringSlice(string) []string
	GetBool(string) bool
}

// Mutable is a read-write view of the Open Match configuration.
type Mutable interface {
	Set(string, interface{})
	View
}

// Sub returns a subset of configuration filtered by the key.
func Sub(v View, key string) View {
	vcfg, ok := v.(*viper.Viper)
	if ok {
		return vcfg.Sub(key)
	}
	return nil
}
