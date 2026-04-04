package wgetm

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	StartURL  string
	MaxDepth  int
	Workers   int
	OutputDir string
	Timeout   time.Duration
	Verbose   bool
}

type Stats struct {
	Downloaded int
	Skipped    int
	Errors     int
}

type Crawler struct {
	cfg    Config
	base   *url.URL
	client *http.Client
}

type crawlJob struct {
	url   string
	depth int
}

type crawlResult struct {
	job        crawlJob
	localPath  string
	isHTML     bool
	children   []crawlJob
	downloaded bool
	err        error
}

type htmlRecord struct {
	rawURL    string
	localPath string
}

func NewCrawler(cfg Config) (*Crawler, error) {
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}

	base, err := url.Parse(cfg.StartURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %q: %w", cfg.StartURL, err)
	}
	if base.Scheme != "http" && base.Scheme != "https" {
		return nil, fmt.Errorf("URL scheme must be http or https, got %q", base.Scheme)
	}

	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return nil, fmt.Errorf("make output dir error: %w", err)
	}

	return &Crawler{
		cfg:    cfg,
		base:   base,
		client: &http.Client{Timeout: cfg.Timeout},
	}, nil
}

func (c *Crawler) Run() Stats {
	jobCh := make(chan crawlJob)
	resCh := make(chan crawlResult, c.cfg.Workers)

	for i := 0; i < c.cfg.Workers; i++ {
		go func() {
			for job := range jobCh {
				resCh <- c.processJob(job)
			}
		}()
	}

	visited := map[string]bool{c.cfg.StartURL: true}
	queue := []crawlJob{{c.cfg.StartURL, 0}}
	inFlight := 0
	var stats Stats
	paths := make(map[string]string)
	var htmlFiles []htmlRecord

	for len(queue) > 0 || inFlight > 0 {

		var sendCh chan crawlJob
		var next crawlJob
		if len(queue) > 0 {
			sendCh = jobCh
			next = queue[0]
		}

		select {
		case sendCh <- next:
			queue = queue[1:]
			inFlight++

		case res := <-resCh:
			inFlight--
			switch {
			case res.err != nil:
				stats.Errors++
				c.logf("[error] %s: %v", res.job.url, res.err)
			case !res.downloaded:
				stats.Skipped++
			default:
				stats.Downloaded++
				paths[res.job.url] = res.localPath
				if res.isHTML {
					htmlFiles = append(htmlFiles, htmlRecord{res.job.url, res.localPath})
				}
			}

			for _, child := range res.children {
				if !visited[child.url] {
					visited[child.url] = true
					queue = append(queue, child)
				}
			}
		}
	}

	close(jobCh)

	c.rewriteAll(htmlFiles, paths)

	return stats
}

func (c *Crawler) processJob(job crawlJob) crawlResult {

	u, err := url.Parse(job.url)
	if err != nil || u.Host != c.base.Host {
		return crawlResult{job: job, downloaded: false}
	}

	c.logf("[fetch] depth=%d %s", job.depth, job.url)

	fr, err := fetch(c.client, job.url)
	if err != nil {
		return crawlResult{job: job, err: err}
	}
	if fr.statusCode < 200 || fr.statusCode >= 400 {
		c.logf("[skip] %s → HTTP %d", job.url, fr.statusCode)
		return crawlResult{job: job, downloaded: false}
	}

	localPath := urlToLocalPath(job.url, fr.contentType)
	if err := saveFile(c.cfg.OutputDir, localPath, fr.body); err != nil {
		return crawlResult{job: job, err: fmt.Errorf("save %s: %w", localPath, err)}
	}
	c.logf("[saved] %s", localPath)

	isHTML := strings.Contains(strings.ToLower(fr.contentType), "text/html")

	var children []crawlJob
	if isHTML && job.depth < c.cfg.MaxDepth {
		base, _ := url.Parse(job.url)
		for _, lnk := range extractLinks(fr.body, base) {

			child, err := url.Parse(lnk.url)
			if err != nil || child.Host != c.base.Host {
				continue
			}
			childDepth := job.depth
			if lnk.isPage {
				childDepth++
			}
			if childDepth <= c.cfg.MaxDepth {
				children = append(children, crawlJob{lnk.url, childDepth})
			}
		}
	}

	return crawlResult{
		job:        job,
		localPath:  localPath,
		isHTML:     isHTML,
		children:   children,
		downloaded: true,
	}
}

func (c *Crawler) rewriteAll(files []htmlRecord, paths map[string]string) {
	for _, f := range files {
		abs := filepath.Join(c.cfg.OutputDir, f.localPath)
		data, err := os.ReadFile(abs)
		if err != nil {
			c.logf("[rewrite error] read %s: %v", f.localPath, err)
			continue
		}

		base, _ := url.Parse(f.rawURL)
		dir := filepath.Dir(f.localPath)

		relMap := make(map[string]string, len(paths))
		for rawURL, target := range paths {
			rel, err := filepath.Rel(dir, target)
			if err != nil {
				rel = "/" + filepath.ToSlash(target)
			} else {
				rel = filepath.ToSlash(rel)
			}
			relMap[rawURL] = rel
		}

		rewritten := rewriteLinks(data, base, relMap)
		if err := os.WriteFile(abs, rewritten, 0o644); err != nil {
			c.logf("[rewrite error] write %s: %v", f.localPath, err)
		}
	}
}

func (c *Crawler) logf(format string, args ...any) {
	if c.cfg.Verbose {
		fmt.Printf(format+"\n", args...)
	}
}
