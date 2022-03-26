package config

import "github.com/spf13/viper"

type appConfig struct {
	// Listening port
	Port string
	// App secret. Used to encrypt authentication data
	Secret string
	// App user creating by default
	DefaultAdmin appUser
}

type appUser struct {
	Username string
	Password string
	Email    string
}

var appcfg appConfig

func appConfInit() {
	port := ""
	if viper.IsSet("PORT") {
		port = viper.GetString("PORT") // port from envvar
	} else {
		port = viper.GetString("service.port") // port from config file
	}
	if port == "" {
		port = "8088" // default app port
	}

	secret := ""
	if viper.IsSet("APP_SECRET") {
		secret = viper.GetString("APP_SECRET") // app secret from envvar
	} else if viper.GetString("service.secret") != "" {
		secret = viper.GetString("service.secret") // app secret from config file
	}

	defAdmin := appUser{
		Username: viper.GetString("service.default_admin.username"),
		Password: viper.GetString("service.default_admin.password"),
		Email:    viper.GetString("service.default_admin.email"),
	}

	if viper.IsSet("ADMIN_PASSWORD") {
		// redefine the Password from the environment variables
		defAdmin.Password = viper.GetString("ADMIN_PASSWORD")
	}

	appcfg = appConfig{
		Port:         port,
		Secret:       secret,
		DefaultAdmin: defAdmin,
	}
}

func GetAppPort() string {
	return appcfg.Port
}

func GetAppSecret() string {
	return appcfg.Secret
}

func GetDefaultAdmin() appUser {
	return appcfg.DefaultAdmin
}
