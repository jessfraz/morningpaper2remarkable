package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

func GetFiles(cfg *AppConfig) error {
	// Iterate over however many pages we want to display.
	for i := 1; i <= cfg.MaxPages; i++ {
		if err := downloadFilesFromFeed(cfg, i); err != nil {
			return err
		}
	}
	return nil
}

// downloadFilesFromFeed parses the RSS feed for the morning paper and downloads the links for
// papers.
func downloadFilesFromFeed(cfg *AppConfig, page int) error {
	// Parse the feed.
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(fmt.Sprintf(cfg.MorningPaperRSSFeedURL, page))
	if err != nil {
		return fmt.Errorf("parsing the url %s failed: %v", cfg.MorningPaperRSSFeedURL, err)
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

		// Create a name for the resulting file from the title.
		name := getNameForPaperFile(item.Title, item.PublishedParsed)
		// Use the item title here because Adrian uses better titles than
		// what is usually in the link for the paper.
		file := filepath.Join(cfg.DataDir, name)

		if paperAlreadySynced(file) {
			logrus.WithFields(logrus.Fields{
				"paper": paper.Text(),
				"link":  paperLink,
			}).Info("skipping paper (already synced)")

			continue
		}

		if err := downloadPaper(paperLink, file); err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"paper": paper.Text(),
			"link":  paperLink,
			"file":  file,
		}).Info("downloaded paper to file")

		// Sync the file with remarkable cloud.
		if err := cfg.RemarkableAPI.SyncFileAndRename(file, fmt.Sprintf("%s (%s)", strings.TrimSpace(item.Title), item.PublishedParsed.Format("2006-01-02"))); err != nil {
			return err
		}

		// Continue here.
		continue
	}

	return nil
}

func paperAlreadySynced(file string) bool {
	// check if file exists
	_, err := os.Stat(file)
	return err == nil
}

func downloadPaper(link, file string) error {
	// Open the file.
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("opening file %s failed: %v", file, err)
	}
	defer f.Close()

	// Get the file contents.
	resp, err := http.Get(link)
	if err != nil {
		f.Close()
		os.Remove(file)
		return fmt.Errorf("getting %s failed: %v", link, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Delete the file.
		f.Close()
		os.Remove(file)
		return fmt.Errorf("status code for getting %s error: %d %s", link, resp.StatusCode, resp.Status)
	}

	// Copy the contents.
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(file)
		return fmt.Errorf("writing file %s failed: %v", file, err)
	}

	return nil
}

func getNameForPaperFile(title string, published *time.Time) string {
	parts := strings.Split(title, "http")
	title = parts[0]

	name := strings.Replace(strings.Replace(strings.ToLower(title), " ", "-", -1), ":", "", -1)

	parts = strings.Split(name, "/")

	// Return the last part.
	return fmt.Sprintf("%s-%s.pdf", published.Format("2006-01-02"), parts[len(parts)-1])
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
