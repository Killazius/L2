package downloader

import (
	"16/internal/config"
	"16/internal/parser"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

type Downloader struct {
	client *http.Client
	cfg    *config.Config
	parser *parser.Parser
}

func New(cfg *config.Config, parser *parser.Parser) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: cfg.TimeoutSeconds(),
		},
		cfg:    cfg,
		parser: parser,
	}
}

func (d *Downloader) download(u *url.URL) (*parser.ParseResult, error) {
	urlString := u.String()

	var resp *http.Response
	var err error
	for try := uint(0); try < d.cfg.Retries; try++ {
		slog.Info("Attempting to download", "url", urlString, "try", try+1)
		resp, err = d.client.Get(urlString)
		if err == nil {
			slog.Info("Downloaded", "url", urlString, "status", resp.StatusCode)
			break
		}
		slog.Warn("Download attempt failed", "url", urlString, "error", err)
	}
	defer resp.Body.Close()
	if err != nil {
		slog.Error("All download attempts failed", "url", urlString, "error", err)
		return nil, err
	}
	result, err := d.parser.Parse(resp, u)
	if err != nil {
		slog.Error("Failed to parse response", "url", urlString, "error", err)
		return nil, err
	}
	if err := d.save(result.Data, result.Filename); err != nil {
		slog.Error("Failed to save data", "filename", result.Filename, "error", err)
		return nil, err
	}
	slog.Info("Saved file", "filename", result.Filename)
	return result, nil
}

func (d *Downloader) save(data []byte, filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (d *Downloader) Start(u *url.URL) error {
	type job struct {
		u     *url.URL
		depth uint
	}

	jobs := make(chan job, d.cfg.NumWorkers*2)
	var wg sync.WaitGroup
	var mu sync.Mutex
	visited := make(map[string]struct{})

	enqueue := func(ju *url.URL, depth uint) {
		key := ju.String()
		mu.Lock()
		if _, ok := visited[key]; ok {
			mu.Unlock()
			return
		}
		visited[key] = struct{}{}
		mu.Unlock()
		wg.Add(1)
		jobs <- job{u: ju, depth: depth}
	}

	enqueue(u, d.cfg.Depth)

	go func() {
		wg.Wait()
		close(jobs)
	}()

	workerCount := int(d.cfg.NumWorkers)
	var workersWG sync.WaitGroup
	workersWG.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer workersWG.Done()
			for jb := range jobs {
				res, err := d.download(jb.u)
				if err != nil {
					wg.Done()
					continue
				}
				if jb.depth > 0 {
					for _, link := range res.SubLinks {
						parsed, perr := url.Parse(link)
						if perr != nil {
							slog.Debug("Failed to parse sublink", "link", link, "error", perr)
							continue
						}
						enqueue(parsed, jb.depth-1)
					}
				}
				wg.Done()
			}
		}()
	}

	workersWG.Wait()
	return nil
}
