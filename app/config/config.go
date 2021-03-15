// Create package in v.1.0.0
// config package contains App global variable with config value using in app(main) package from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define appConfig type which is type of App
// appConfig dose not implement any interface, so use by explicitly importing in app(main) package
// Also, App implement various config interface each of package in syscheck domain by declaring method

package config

// appConfig having config value and return that value with method. Not implement interface
type appConfig struct {
	// esAddress represent host address of elasticsearch server
	esAddress *string
}
