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

	"golang.org/x/sync/errgroup"

	"mvdan.cc/xurls/v2"

	"github.com/scylladb/go-set/strset"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

func (lgc *logic) Check(stdin io.Reader, cfgPath string) error {
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
		if lgc.logic.IsIgnoredURL(u) {
			delete(urls, u)
		}
	}

	return lgc.logic.CheckURLs(urls)
}

func (lgc *logic) IsIgnoredURL(uri string) bool {
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
	for _, u := range lgc.cfg.IgnoreURLs {
		if uri == u {
			return true
		}
	}
	for _, h := range lgc.cfg.IgnoreHosts {
		if u.Host == h {
			return true
		}
	}
	return false
}

func (lgc *logic) CheckURLs(urls map[string]*strset.Set) error {
	if len(urls) == 0 {
		return nil
	}
	if lgc.cfg.MaxRequestCount == 0 {
		lgc.cfg.MaxRequestCount = domain.DefaultMaxRequestCount
	}
	semaphore := make(chan struct{}, lgc.cfg.MaxRequestCount)
	resultChan := make(chan error, len(urls))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for u, files := range urls {
		go func(u string, files *strset.Set) {
			semaphore <- struct{}{}
			err := lgc.logic.CheckURL(ctx, u)
			<-semaphore
			if err == nil {
				resultChan <- nil
				return
			}
			resultChan <- fmt.Errorf("failed to check a url %s %s: %w", u, files.String(), err)
		}(u, files)
	}
	endCount := len(urls)
	failedCount := lgc.cfg.MaxFailedRequestCount
	for {
		select {
		case err := <-resultChan:
			endCount--
			if err != nil {
				failedCount--
				fmt.Fprintln(os.Stderr, err)
			}
			if endCount == 0 {
				if failedCount != lgc.cfg.MaxFailedRequestCount {
					return fmt.Errorf("")
				}
				return nil
			}
			if lgc.cfg.MaxFailedRequestCount != -1 && failedCount == -1 {
				return fmt.Errorf("too many urls are dead")
			}
		case <-ctx.Done():
			return fmt.Errorf("context is caceled")
		}
	}
}

func (lgc *logic) CheckURLWithMethod(
	ctx context.Context, u, method string,
) error {
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := lgc.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	// check status code
	if resp.StatusCode/100 != 2 { //nolint:gomnd
		return fmt.Errorf("%s is dead (%d)", u, resp.StatusCode)
	}
	return nil
}

func (lgc *logic) CheckURL(ctx context.Context, u string) error {
	switch lgc.cfg.HTTPMethod {
	case "head,get":
		if err := lgc.logic.CheckURLWithMethod(ctx, u, http.MethodHead); err == nil {
			return nil
		}
		return lgc.logic.CheckURLWithMethod(ctx, u, http.MethodGet)
	case "":
		if err := lgc.logic.CheckURLWithMethod(ctx, u, http.MethodHead); err == nil {
			return nil
		}
		return lgc.logic.CheckURLWithMethod(ctx, u, http.MethodGet)
	case "get":
		return lgc.logic.CheckURLWithMethod(ctx, u, http.MethodGet)
	case "head":
		return lgc.logic.CheckURLWithMethod(ctx, u, http.MethodHead)
	default:
		return fmt.Errorf(`invalid http_method_type: %s`, lgc.cfg.HTTPMethod)
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
		f := f
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
			return urls, fmt.Errorf("failed to read %s: %w", p, err)
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
