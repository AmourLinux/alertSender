package cmd

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var cfgFile string
var endpoint string
var apikey string
var secret string
var businessName string
var message string

func sign(form string) string {
	str := "alarmMessage=" + message + "&apiKey=" + apikey + "&businessName=" + businessName + "&secret=" + secret
	sum := md5.Sum([]byte(str))
	s := fmt.Sprintf("%x", sum)
	return s
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alert-sender",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		message, _ = strconv.Unquote("\"" + message + "\"")
		apikey = viper.GetString("apikey")
		secret = viper.GetString("secret")
		if apikey == "" {
			fmt.Println("apikey must present")
			os.Exit(1)
		}
		if secret == "" {
			fmt.Println("secret must present")
			os.Exit(1)
		}
		if businessName == "" {
			fmt.Println("business-name must present")
			os.Exit(1)
		}

		var httpClient = &http.Client{
			Timeout: time.Second * 10,
		}

		form := url.Values{}
		form.Add("alarmMessage", message)
		form.Add("apiKey", apikey)
		form.Add("businessName", businessName)
		form.Add("sign", sign(form.Encode()))
		req, _ := http.NewRequest("POST", endpoint+"/alarm/send-message", bytes.NewBufferString(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := httpClient.Do(req)
		if err != nil {
			os.Exit(1)
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		if !strings.HasPrefix(resp.Status, "2") {
			os.Exit(1)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "http://alarm.bafang.com", "alarm service endpoint")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alert-sender.yaml)")
	_ = viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))

	rootCmd.PersistentFlags().StringVar(&apikey, "apikey", "", "apikey")
	_ = viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))

	rootCmd.PersistentFlags().StringVar(&secret, "secret", "", "secret")
	_ = viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))

	rootCmd.PersistentFlags().StringVar(&businessName, "business-name", "", "business-name")
	rootCmd.PersistentFlags().StringVar(&message, "message", "", "message to send")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".alter-sender" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".alter-sender")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
