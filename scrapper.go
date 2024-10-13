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
	var screenshot bool

	// Define flags for extracting links, saving HTML, and taking a screenshot
	flag.BoolVar(&links, "links", false, "Extract links from the webpage")
	flag.BoolVar(&html, "html", false, "Save HTML content of the webpage")
	flag.BoolVar(&screenshot, "screenshot", false, "Take a screenshot of the webpage")

	// Handle the help flag to display usage instructions
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of the tool")
		fmt.Println("  -html\n        Save HTML content of the webpage")
		fmt.Println("  -links\n        Extract links from the webpage")
		fmt.Println("  -screenshot\n        Take a screenshot of the webpage")
		fmt.Println("\nExample:")
		fmt.Printf("go run scrapper.go  -html -links -screenshot https://example.com\n")
	}

	// Parse flags
	flag.Parse()

	// Check if a URL was provided; if not, print an error and exit
	if len(flag.Args()) == 0 {
		log.Fatal("No URL provided. Please input a URL to scrape")
	}

	// The first argument is the URL
	url := flag.Args()[0]

	// Check if no flags were chosen
	if !links && !html && !screenshot {
		fmt.Println("No flags were chosen, setting all flags for", url)
		links = true
		html = true
		screenshot = true
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

	c.Visit(url)

	// Capture screenshot if the flag is set
	if screenshot {
		if err := captureScreenshot(url); err != nil {
			log.Fatal("Error capturing screenshot:", err)
		}
		fmt.Println("Screenshot captured successfully!")
	}
}

// Screenshot func
func captureScreenshot(url string) error {
	// Chrome init
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel() // Release the browser resources when no longer needed

	var screenshotBuffer []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery), // Wait until the page is fully loaded
		chromedp.FullScreenshot(&screenshotBuffer, 90), // Take screenshot
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
