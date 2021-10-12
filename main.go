package main

import (
	"GO/nomad/scrapper"
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"
)

const File_Name string = "jobs.csv"

func handleHome(c echo.Context) error {
	// return c.String(http.StatusOK, "Kimbs Hello, World!")
	return c.File("home.html")
}
func handleScrape(c echo.Context) error {
	// return c.String(http.StatusOK, "Kimbs Hello, World!")
	defer os.Remove(File_Name)
	fmt.Println("term:", c.FormValue("term"))
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(File_Name, File_Name)
}

func main() {
	// scrapper.Scrape("python")
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
