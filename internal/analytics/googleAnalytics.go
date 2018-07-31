package analytics

import (
	"encoding/json"
	"fmt"
	"github.com/liagame/lia-cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
)

const propertyID = "UA-122844498-1"

func isTrackingAllowed() bool {
	return viper.GetBool("analyticsAllow")
}

func Log(category string, action string, metadata map[string]string) {
	if !isTrackingAllowed() {
		return
	}

	v := url.Values{
		"v":   {"1"},
		"tid": {propertyID},
		"cid": {viper.GetString("TrackingId")},
		"t":   {"event"},
		"ec":  {category},
		"ea":  {action},
		//"ua":  {r.UserAgent()},
	}
	dataJson, err := json.Marshal(metadata)
	if err != nil {
		os.Exit(lia_cli.PreparingAnalyticsDataFailed)
	}

	v.Set("el", string(dataJson))

	http.PostForm("https://www.google-analytics.com/collect", v)
}

func ParseStringFlag(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		os.Exit(lia_cli.PreparingAnalyticsDataFailed)
	}
	return value
}

func ParseIntFlagToString(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetInt(name)
	if err != nil {
		os.Exit(lia_cli.PreparingAnalyticsDataFailed)
	}
	return fmt.Sprint(value)
}

func ParseBoolFlagToString(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetBool(name)
	if err != nil {
		os.Exit(lia_cli.PreparingAnalyticsDataFailed)
	}
	return fmt.Sprint(value)
}

func ParseIntSliceFlagToString(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetIntSlice(name)
	if err != nil {
		os.Exit(lia_cli.PreparingAnalyticsDataFailed)
	}
	return fmt.Sprint(value)
}
