package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

func main() {
	var links bool
	var html bool

	flag.BoolVar(&links, "l", false, "Extract links from the webpage")
	flag.BoolVar(&html, "h", false, "Save HTML content of the webpage")
	flag.Parse()

	// Check if no flags were chosen
	if !links && !html {
		fmt.Println("No flags were chosen, setting all flags for https://sibervatan.org.")
		links = true
		html = true
	}

	// Check if a URL was provided; if not, use the default URL
	url := "https://sibervatan.org/"
	if len(flag.Args()) > 0 {
		url = flag.Args()[0]
	}

	// Initialize a new Colly collector
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("An error occurred: ", err)
	})

	// Extract links
	if links {
		linksFile, err := os.Create("links.txt")
		if err != nil {
			log.Fatal("File couldn't be created: ", err)
		}
		defer linksFile.Close()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			if strings.HasPrefix(link, "http") {
				_, err := linksFile.WriteString(link + "\n")
				if err != nil {
					log.Println("Unable to write link:", err)
				} else {
					fmt.Println("Link Detected:", link)
				}
			}
		})
	}

	// Save HTML content
	if html {
		allFile, err := os.Create("HTML.txt")
		if err != nil {
			log.Fatal("File couldn't be created: ", err)
		}
		defer allFile.Close()

		c.OnResponse(func(r *colly.Response) {
			_, err := allFile.WriteString(string(r.Body))
			if err != nil {
				log.Println("Unable to write HTML content:", err)
			} else {
				fmt.Println("HTML content written to HTML.txt")
			}
		})
	}

	// Visit the specified website
	c.Visit(url)

	// Capture screenshot of the given URL
	if err := captureScreenshot(url); err != nil {
		log.Fatal("Error capturing screenshot:", err)
	}

	fmt.Println("Screenshot captured successfully!")
}

// Function to capture screenshot of the given URL
func captureScreenshot(url string) error {

	// Initialize a controllable Chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel() // Release the browser resources when no longer needed

	var screenshotBuffer []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery), // Wait untill the page is fully loaded
		chromedp.FullScreenshot(&screenshotBuffer, 90), // Full screenshot of the entire page
	)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// Write the screenshot to an image file
	err = os.WriteFile("screenshot.png", screenshotBuffer, 0644)
	if err != nil {
		return fmt.Errorf("failed to save screenshot to file: %w", err)
	}

	return nil
}
