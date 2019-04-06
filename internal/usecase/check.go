package usecase

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/sync/errgroup"

	"gopkg.in/yaml.v2"

	"mvdan.cc/xurls"

	"github.com/pkg/errors"
	"github.com/scylladb/go-set/strset"
	"github.com/suzuki-shunsuke/go-cliutil"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

// Check checks whether dead urls are included in files.
func Check(fsys domain.Fsys, stdin io.Reader, cfgPath string) error {
	cfg, err := readCfg(fsys, cfgPath)
	if err != nil {
		return err
	}

	// get file paths
	files, err := getFiles(stdin)
	if err != nil {
		return err
	}
	// urls is a map whose key is url and value is file paths which include the url
	// url -> file paths
	urls, err := extractURLsFromFiles(fsys, files)
	if err != nil {
		return err
	}
	// filter url
	for u := range urls {
		uri, err := url.Parse(u)
		// ignore url if it is failed to parse the url
		if err != nil {
			delete(urls, u)
			continue
		}
		if isIgnoredURL(uri) {
			delete(urls, u)
			continue
		}
	}

	for _, u := range cfg.IgnoreURLs {
		delete(urls, u)
	}
	return checkURLs(urls)
}

func isIgnoredURL(u *url.URL) bool {
	if u.Scheme != "http" && u.Scheme != "https" {
		return true
	}
	return false
}

func findCfg(fsys domain.Fsys) (string, error) {
	wd, err := fsys.Getwd()
	if err != nil {
		return "", err
	}
	return cliutil.FindFile(wd, ".durl.yml", fsys.Exist)
}

func readCfg(fsys domain.Fsys, cfgPath string) (domain.Cfg, error) {
	cfg := domain.Cfg{}
	if cfgPath == "" {
		d, err := findCfg(fsys)
		if err != nil {
			return cfg, err
		}
		cfgPath = d
	}
	rc, err := fsys.Open(cfgPath)
	if err != nil {
		return cfg, err
	}
	defer rc.Close()
	if err := yaml.NewDecoder(rc).Decode(&cfg); err != nil {
		return cfg, err
	}
	return initCfg(cfg)
}

func initCfg(cfg domain.Cfg) (domain.Cfg, error) {
	return cfg, nil
}

func checkURLs(urls map[string]*strset.Set) error {
	eg, ctx := errgroup.WithContext(context.Background())
	client := http.Client{
		Timeout: domain.DefaultTimeout,
	}
	for u, files := range urls {
		// https://golang.org/doc/faq#closures_and_goroutines
		u := u
		files := files
		eg.Go(func() error {
			if err := checkURL(ctx, client, u); err != nil {
				return errors.Wrap(err, files.String())
			}
			return nil
		})
	}
	return eg.Wait()
}

func checkURL(ctx context.Context, client http.Client, u string) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	// check status code
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s is dead (%d)", u, resp.StatusCode)
	}
	return nil
}

func extractURLsFromFiles(fsys domain.Fsys, files *strset.Set) (map[string]*strset.Set, error) {
	// return a map whose key is url and value is file paths which include the url
	size := files.Size()
	if size == 0 {
		return nil, nil
	}
	// extract urls from all files
	type (
		File struct {
			path string
			urls *strset.Set
		}
	)
	urlsChan := make(chan File, size)
	eg, ctx := errgroup.WithContext(context.Background())
	files.Each(func(p string) bool {
		eg.Go(func() error {
			// open a file and extract urls from it
			arr, err := extractURLsFromFile(ctx, fsys, p)
			if err != nil {
				return err
			}
			urlsChan <- File{path: p, urls: arr}
			return nil
		})
		return true
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(urlsChan)
	// url -> file paths
	urls := map[string]*strset.Set{}
	for f := range urlsChan {
		f.urls.Each(func(u string) bool {
			v, ok := urls[u]
			if ok {
				v.Add(f.path)
				return true
			}
			urls[u] = strset.New(f.path)
			return true
		})
	}
	return urls, nil
}

func extractURLsFromFile(ctx context.Context, fsys domain.Fsys, p string) (*strset.Set, error) {
	// open a file and extract urls from it
	urls := strset.New()
	errChan := make(chan error, 1)
	reg := xurls.Strict()
	go func() {
		// open a file
		fi, err := fsys.Open(p)
		if err != nil {
			errChan <- err
			return
		}
		defer fi.Close()
		// read a file per a line
		scanner := bufio.NewScanner(fi)
		for scanner.Scan() {
			// extract urls from a line
			urls.Add(reg.FindAllString(scanner.Text(), -1)...)
		}
		errChan <- scanner.Err()
	}()
	select {
	case <-ctx.Done():
		return nil, nil
	case err := <-errChan:
		if err != nil {
			return urls, errors.Wrapf(err, "failed to read %s", p)
		}
		return urls, nil
	}
}

func getFiles(stdin io.Reader) (*strset.Set, error) {
	files := strset.New()
	if stdin != nil {
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			files.Add(scanner.Text())
		}
		return files, scanner.Err()
	}
	return files, nil
}
