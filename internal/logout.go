package internal

import (
	"fmt"
	"github.com/liagame/lia-SDK/internal/config"
	"github.com/spf13/viper"
)

func Logout() {
	config.UserToken = ""
	config.UserToken = ""

	viper.Set("username", "")
	viper.Set("token", "")
	viper.WriteConfig()

	fmt.Println("Logged out successfully.")
}