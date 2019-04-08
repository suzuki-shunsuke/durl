package usecase

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"gopkg.in/yaml.v2"

	"mvdan.cc/xurls"

	"github.com/pkg/errors"
	"github.com/scylladb/go-set/strset"
	"github.com/suzuki-shunsuke/go-cliutil"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

func (lgc *logic) Check(stdin io.Reader, cfgPath string) error {
	cfg, err := lgc.logic.ReadCfg(cfgPath)
	if err != nil {
		return err
	}

	// get file paths
	files, err := lgc.logic.GetFiles(stdin)
	if err != nil {
		return err
	}
	// urls is a map whose key is url and value is file paths which include the url
	// url -> file paths
	urls, err := lgc.logic.ExtractURLsFromFiles(files)
	if err != nil {
		return err
	}
	// filter url
	for u := range urls {
		if lgc.logic.IsIgnoredURL(u, cfg) {
			delete(urls, u)
		}
	}

	return lgc.logic.CheckURLs(cfg, urls)
}

func (lgc *logic) IsIgnoredURL(uri string, cfg domain.Cfg) bool {
	u, err := url.Parse(uri)
	if err != nil {
		// ignore url if it is failed to parse the url
		return true
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return true
	}
	for _, ignoreHost := range domain.IgnoreHosts {
		if u.Host == ignoreHost || strings.HasPrefix(u.Host, fmt.Sprintf("%s:", ignoreHost)) {
			return true
		}
	}
	for _, u := range cfg.IgnoreURLs {
		if uri == u {
			return true
		}
	}
	for _, h := range cfg.IgnoreHosts {
		if u.Host == h {
			return true
		}
	}
	return false
}

func (lgc *logic) FindCfg() (string, error) {
	wd, err := lgc.fsys.Getwd()
	if err != nil {
		return "", err
	}
	return cliutil.FindFile(wd, ".durl.yml", lgc.fsys.Exist)
}

func (lgc *logic) ReadCfg(cfgPath string) (domain.Cfg, error) {
	cfg := domain.Cfg{
		HTTPMethod: "head,get",
	}
	if cfgPath == "" {
		d, err := lgc.logic.FindCfg()
		if err != nil {
			return cfg, err
		}
		cfgPath = d
	}
	rc, err := lgc.fsys.Open(cfgPath)
	if err != nil {
		return cfg, err
	}
	defer rc.Close()
	if err := yaml.NewDecoder(rc).Decode(&cfg); err != nil {
		return cfg, err
	}
	return lgc.logic.InitCfg(cfg)
}

func (lgc *logic) InitCfg(cfg domain.Cfg) (domain.Cfg, error) {
	methods := map[string]struct{}{
		"get":      {},
		"head":     {},
		"head,get": {},
	}
	if _, ok := methods[cfg.HTTPMethod]; !ok {
		return cfg, fmt.Errorf(`invalid http_method_type: %s`, cfg.HTTPMethod)
	}
	return cfg, nil
}

func (lgc *logic) CheckURLs(cfg domain.Cfg, urls map[string]*strset.Set) error {
	if len(urls) == 0 {
		return nil
	}
	if cfg.HTTPRequestTimeout == 0 {
		cfg.HTTPRequestTimeout = domain.DefaultTimeout
	}
	client := http.Client{
		Timeout: domain.DefaultTimeout * time.Second,
	}
	if cfg.MaxRequestCount == 0 {
		cfg.MaxRequestCount = domain.DefaultMaxRequestCount
	}
	semaphore := make(chan struct{}, cfg.MaxRequestCount)
	resultChan := make(chan error, len(urls))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for u, files := range urls {
		go func(u string, files *strset.Set) {
			semaphore <- struct{}{}
			err := lgc.logic.CheckURL(ctx, cfg, client, u)
			<-semaphore
			if err == nil {
				resultChan <- nil
				return
			}
			resultChan <- errors.Wrapf(err, "failed to check a url: %s %s", u, files.String())
		}(u, files)
	}
	endCount := len(urls)
	failedCount := cfg.MaxFailedRequestCount
	for {
		select {
		case err := <-resultChan:
			endCount--
			if err != nil {
				failedCount--
				fmt.Fprintln(os.Stderr, err)
			}
			if endCount == 0 {
				if failedCount != cfg.MaxFailedRequestCount {
					return fmt.Errorf("")
				}
				return nil
			}
			if failedCount == -1 {
				return fmt.Errorf("too many urls are dead")
			}
		}
	}
}

func (lgc *logic) CheckURLWithMethod(
	ctx context.Context, client http.Client, u, method string,
) error {
	req, err := http.NewRequest(method, u, nil)
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
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("%s is dead (%d)", u, resp.StatusCode)
	}
	return nil
}

func (lgc *logic) CheckURL(
	ctx context.Context, cfg domain.Cfg, client http.Client, u string,
) error {
	switch cfg.HTTPMethod {
	case "head,get":
		if err := lgc.logic.CheckURLWithMethod(ctx, client, u, http.MethodHead); err == nil {
			return nil
		}
		return lgc.logic.CheckURLWithMethod(ctx, client, u, http.MethodGet)
	case "":
		if err := lgc.logic.CheckURLWithMethod(ctx, client, u, http.MethodHead); err == nil {
			return nil
		}
		return lgc.logic.CheckURLWithMethod(ctx, client, u, http.MethodGet)
	case "get":
		return lgc.logic.CheckURLWithMethod(ctx, client, u, http.MethodGet)
	case "head":
		return lgc.logic.CheckURLWithMethod(ctx, client, u, http.MethodHead)
	default:
		return fmt.Errorf(`invalid http_method_type: %s`, cfg.HTTPMethod)
	}
}

func (lgc *logic) ExtractURLsFromFiles(files *strset.Set) (map[string]*strset.Set, error) {
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
			arr, err := lgc.logic.ExtractURLsFromFile(ctx, p)
			if err != nil {
				// https://github.com/suzuki-shunsuke/durl/issues/27
				fmt.Fprintf(os.Stderr, "failed to extract urls from a file %s: %s\n", p, err)
				return nil
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

func (lgc *logic) ExtractURLsFromFile(ctx context.Context, p string) (*strset.Set, error) {
	// open a file and extract urls from it
	urls := strset.New()
	errChan := make(chan error, 1)
	reg := xurls.Strict()
	go func() {
		// open a file
		fi, err := lgc.fsys.Open(p)
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

func (lgc *logic) GetFiles(stdin io.Reader) (*strset.Set, error) {
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
