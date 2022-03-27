package superFlags

import (
	"github.com/superwhys/superGo/superSlices"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	lg "github.com/superwhys/superGo/superLog"
	// Import remote config
	_ "github.com/spf13/viper/remote"
)

var (
	allKeys     []string
	requiredKey []string
	config      *string
	debug       *bool

	v = viper.New()
)

func initFlags() {
	v.AddConfigPath(".")
	v.AddConfigPath("./etc/")
	err := v.BindPFlags(pflag.CommandLine)
	if err != nil {
		lg.Fatal("BindPFlags Error!")
	}
	config = pflag.StringP("config", "f", "", "Specify config file to parse. Support json, yaml, toml etc.")
	debug = pflag.Bool("debug", false, "Set true to enable debug mode")

	allKeys = append(allKeys, "debug", "owner")
}

// Parse has to called after main() before any application code.
func Parse() {
	initFlags()
	pflag.Parse()
	if *debug {
		lg.EnableDebug()
	}

	for _, k := range requiredKey {
		if isZero(v.Get(k)) {
			lg.Fatal("Missing", k)
		}
	}
	expectedKeys := superSlices.NewStringSet(nil)
	for _, k := range allKeys {
		if err := expectedKeys.Add(strings.ToLower(k)); err != nil {
			lg.Fatalf("Add Key Error: --%s", k)
		}
	}

	for _, k := range v.AllKeys() {
		if strings.Contains(k, ".") {
			// Ignore nested key
			continue
		}
		if !expectedKeys.Contains(k) {
			lg.Fatalf("Unknown flag in config: --%s", k)
		}
	}

	if config != nil && *config != "" {
		v.SetConfigFile(*config)
		if err := v.ReadInConfig(); err != nil {
			lg.Errorf("Failed to read on local file, ", err)
		} else {
			lg.Info("Read config from local file!")
		}
	}

	if v.GetBool("debug") {
		lg.EnableDebug()
	}
}

func isZero(i interface{}) bool {
	switch i.(type) {
	case bool:
		// It's trivial to check a bool, since it makes the flag no sense(always true).
		return !i.(bool)
	case string:
		return i.(string) == ""
	case time.Duration:
		return i.(time.Duration) == 0
	case float64:
		return i.(float64) == 0
	case int:
		return i.(int) == 0
	case []string:
		return len(i.([]string)) == 0
	case []interface{}:
		return len(i.([]interface{})) == 0
	default:
		return true
	}
}

func String(key, defaultValue, usage string) func() string {
	pflag.String(key, defaultValue, usage)
	v.SetDefault(key, defaultValue)
	err := v.BindPFlag(key, pflag.Lookup(key))
	if err != nil {
		lg.Fatalf("BindPFlag err, Key: --%s", key)
	}
	allKeys = append(allKeys, key)
	return func() string {
		return v.GetString(key)
	}
}

func StringRequired(key, usage string) func() string {
	requiredKey = append(requiredKey, key)
	allKeys = append(allKeys, key)
	return String(key, "", usage)
}

func Bool(key string, defaultValue bool, usage string) func() bool {
	pflag.Bool(key, defaultValue, usage)
	v.SetDefault(key, defaultValue)
	err := v.BindPFlag(key, pflag.Lookup(key))
	if err != nil {
		lg.Fatalf("BindPFlag err, Key: --%s", key)
	}
	allKeys = append(allKeys, key)
	return func() bool {
		return v.GetBool(key)
	}
}

func BoolRequired(key, usage string) func() bool {
	requiredKey = append(requiredKey, key)
	allKeys = append(allKeys, key)
	return Bool(key, false, usage)
}

func Int(key string, defaultValue int, usage string) func() int {
	pflag.Int(key, defaultValue, usage)
	v.SetDefault(key, defaultValue)
	err := v.BindPFlag(key, pflag.Lookup(key))
	if err != nil {
		lg.Fatalf("BindPFlag err, Key: --%s", key)
	}
	allKeys = append(allKeys, key)
	return func() int {
		return v.GetInt(key)
	}
}

func IntRequired(key, usage string) func() int {
	requiredKey = append(requiredKey, key)
	allKeys = append(allKeys, key)
	return Int(key, 0, usage)
}

func Slice(key string, defaultValue []string, usage string) func() []string {
	pflag.StringSlice(key, defaultValue, usage)
	v.SetDefault(key, defaultValue)
	err := v.BindPFlag(key, pflag.Lookup(key))
	if err != nil {
		lg.Fatalf("BindPFlag err, Key: --%s", key)
	}
	allKeys = append(allKeys, key)
	return func() []string {
		return v.GetStringSlice(key)
	}
}

func Float64(key string, defaultValue float64, usage string) func() float64 {
	pflag.Float64(key, defaultValue, usage)
	v.SetDefault(key, defaultValue)
	err := v.BindPFlag(key, pflag.Lookup(key))
	if err != nil {
		lg.Fatalf("BindPFlag err, Key: --%s", key)
	}
	allKeys = append(allKeys, key)
	return func() float64 {
		return v.GetFloat64(key)
	}
}

func Float64Required(key, usage string) func() float64 {
	requiredKey = append(requiredKey, key)
	allKeys = append(allKeys, key)
	return Float64(key, 0, usage)
}

func Duration(key string, defaultValue time.Duration, usage string) func() time.Duration {
	pflag.Duration(key, defaultValue, usage)
	v.SetDefault(key, defaultValue)
	err := v.BindPFlag(key, pflag.Lookup(key))
	if err != nil {
		lg.Fatalf("BindPFlag err, Key: --%s", key)
	}
	allKeys = append(allKeys, key)
	return func() time.Duration {
		return v.GetDuration(key)
	}
}

func DurationRequired(key, usage string) func() time.Duration {
	requiredKey = append(requiredKey, key)
	allKeys = append(allKeys, key)
	return Duration(key, 0, usage)
}
