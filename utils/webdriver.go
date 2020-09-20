package utils

// Test the web pages automatically
// requires to put chromedriver.exe on root of your application

import (
	"github.com/sclevine/agouti"
	logrus "github.com/sirupsen/logrus"
	_ "os"
)

var webdriver *agouti.WebDriver

// InitializeWebdriver Test All pages loaded when server launches
func InitializeWebdriver(logger *logrus.Entry) {
	webdriver = agouti.ChromeDriver()

	err := webdriver.Start()
	if err != nil {
		logger.Debug(err)
	}
	defer webdriver.Stop()

	browser, err := webdriver.NewPage()
	if err != nil {
		logger.Error("Failed to open page")
		logger.Error(err)
	}

	err = browser.Navigate("https://google.com/")
	if err != nil {
		logger.Error(err)
	}
	html, _ := browser.HTML()
	logger.Debug(html)
}
