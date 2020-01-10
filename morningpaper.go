package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

func getFiles() error {
	// Iterate over however many pages we want to display.
	for i := 1; i <= maxPages; i++ {
		if err := downloadFilesFromFeed(i); err != nil {
			return err
		}
	}
	return nil
}

// downloadFilesFromFeed parses the RSS feed for the morning paper and downloads the links for
// papers.
func downloadFilesFromFeed(page int) error {
	// Parse the feed.
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(fmt.Sprintf(morningPaperRSSFeedURL, page))
	if err != nil {
		return fmt.Errorf("parsing the url %s failed: %v", morningPaperRSSFeedURL, err)
	}

	// Iterate over the items.
	for _, item := range feed.Items {
		if item == nil || item.PublishedParsed == nil {
			// Continue early.
			continue
		}

		// Ignore End of Term
		if strings.HasPrefix("end of term", strings.ToLower(item.Title)) {
			continue
		}

		// Ignore The Year Ahead
		if strings.HasPrefix("the year ahead", strings.ToLower(item.Title)) {
			continue
		}

		logrus.WithFields(logrus.Fields{
			"title":     item.Title,
			"published": item.PublishedParsed.String(),
		}).Debug("parsing article")

		// Try to get the first link to the paper from the content.
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(item.Content))
		if err != nil {
			return fmt.Errorf("parsing article %q content as HTML failed: %v", item.Title, err)
		}
		paper := doc.Find("a")
		paperLink, ok := paper.Attr("href")
		if !ok {
			return fmt.Errorf("paper link for article %q does not have an href", item.Title)
		}

		logrus.WithFields(logrus.Fields{
			"title": item.Title,
			"paper": paper.Text(),
			"link":  paperLink,
		}).Debug("found paper link")

		for !strings.HasSuffix(paperLink, ".pdf") && !strings.HasSuffix(paperLink, "/REF") {
			// Handle arxiv papers.
			if strings.HasPrefix(paperLink, "https://arxiv.org") {
				// Get the pdf link for arxiv.org.
				parts := strings.Split(strings.Trim(paperLink, "/"), "/")
				paperLink = fmt.Sprintf("https://arxiv.org/pdf/%s.pdf", parts[len(parts)-1])
				continue
			}

			// Try to see if we have a link for it in our known papers.
			pl, ok := knownPapersDownloadLinks[paperLink]
			if ok {
				paperLink = pl
				break
			}

			// Try to find the link on the next page.
			paperLink, err = tryToFindPDFLink(paperLink)
			if err != nil {
				return err
			}

			// Bail.
			break
		}

		if len(paperLink) < 1 {
			// Maybe throw an error?
			logrus.WithFields(logrus.Fields{
				"title": item.Title,
				"paper": paper.Text(),
				"link":  paperLink,
			}).Warn("could not find PDF to download")

			continue
		}

		// Download the pdf.
		logrus.WithFields(logrus.Fields{
			"link": paperLink,
		}).Debug("downloading paper")

		file, err := downloadPaper(paperLink)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"paper": paper.Text(),
			"link":  paperLink,
		}).Info("downloaded paper")

		// Sync the file with remarkable cloud.
		if err := rmAPI.SyncFileAndRename(dataDir, file, fmt.Sprintf("%s (%s)", strings.TrimSpace(item.Title), item.PublishedParsed.Format("2006-01-02"))); err != nil {
			return err
		}

		// Continue here.
		continue
	}

	return nil
}

func downloadPaper(link string) ([]byte, error) {
	// Get the file contents.
	resp, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("getting %s failed: %v", link, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code for getting %s error: %d %s", link, resp.StatusCode, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body failed: %v", err)
	}

	return b, nil
}

func tryToFindPDFLink(link string) (string, error) {
	// Request the HTML page.
	resp, err := http.Get(link)
	if err != nil {
		return "", fmt.Errorf("getting %s failed: %v", link, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code for getting %s error: %d %s", link, resp.StatusCode, resp.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parsing link %s failed: %v", link, err)
	}

	// Iterate over all the links.
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, ok := s.Attr("href")
		if !ok {
			return true
		}
		text := s.Text()

		if text == "PDF" && strings.HasPrefix(link, "https://dl.acm.org") {
			// Return false to break.
			// Cannot download from ACM.
			return false
		}

		// Found a link to a pdf.
		if strings.HasPrefix(href, ".pdf") {
			link = href

			// Return false to break.
			return false
		}

		return true
	})

	return "", nil
}
