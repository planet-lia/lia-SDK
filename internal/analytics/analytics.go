package analytics

import (
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-SDK"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
		"net/url"
	"net/http"
	"github.com/liagame/lia-SDK/internal/config"
	"strings"
	"os"
	)


func isTrackingAllowed() bool {
	return viper.GetBool("analyticsAllow")
}

func isTestingMode() bool {
	return viper.Get("testing") == nil
}

func getPropertyId() string {
	if isTestingMode() {
		return config.PropertyID
	} else {
		return config.TestPropertyID
	}
}

func Log(category string, action string, metadata map[string]string) {
	if !isTrackingAllowed() {
		return
	}

	v := url.Values{
		"v":   {"1"},
		"tid": {getPropertyId()},
		"cid": {viper.GetString("TrackingId")},
		"t":   {"event"},
		"ec":  {category},
		"ea":  {action},
		//"ua":  {r.UserAgent()},
	}
	dataJson, err := json.Marshal(metadata)
	if err != nil {
		os.Exit(lia_SDK.PreparingAnalyticsDataFailed)
	}

	v.Set("el", string(dataJson))

	http.PostForm("https://www.google-analytics.com/collect", v)

}

func ParseStringFlag(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		os.Exit(lia_SDK.PreparingAnalyticsDataFailed)
	}
	return value
}

func ParseIntFlagToString(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetInt(name)
	if err != nil {
		os.Exit(lia_SDK.PreparingAnalyticsDataFailed)
	}
	return fmt.Sprint(value)
}

func ParseBoolFlagToString(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetBool(name)
	if err != nil {
		os.Exit(lia_SDK.PreparingAnalyticsDataFailed)
	}
	return fmt.Sprint(value)
}

func ParseIntSliceFlagToString(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetIntSlice(name)
	if err != nil {
		os.Exit(lia_SDK.PreparingAnalyticsDataFailed)
	}
	return fmt.Sprint(value)
}

func TrimPath(path string) string {
	sep := "/"
	if config.OperatingSystem == "windows" {
		sep = "\\"
	}
	parts := strings.Split(path, sep)
	return parts[len(parts) - 1]
}