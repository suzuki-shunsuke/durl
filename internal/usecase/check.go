package usecase

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/sync/errgroup"

	"mvdan.cc/xurls"

	"github.com/scylladb/go-set/strset"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

func Check(fsys domain.Fsys, stdin io.Reader) error {
	files, err := getFiles(stdin)
	if err != nil {
		return err
	}
	urls, err := extractURLsFromFiles(fsys, files)
	if err != nil {
		return err
	}
	return checkURLs(urls)
}

func checkURLs(urls map[string]*strset.Set) error {
	eg, ctx := errgroup.WithContext(context.Background())
	client := http.Client{
		Timeout: domain.DefaultTimeout,
	}
	for u := range urls {
		eg.Go(func() error {
			// GET url
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
				return fmt.Errorf("%s is broken (%d)", u, resp.StatusCode)
			}
			return nil
		})
	}
	return eg.Wait()
}

func extractURLsFromFiles(fsys domain.Fsys, files *strset.Set) (map[string]*strset.Set, error) {
	size := files.Size()
	if size == 0 {
		return nil, nil
	}
	// extract urls
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
	urls := strset.New()
	errChan := make(chan error, 1)
	reg := xurls.Strict()
	go func() {
		fi, err := fsys.Open(p)
		if err != nil {
			errChan <- err
			return
		}
		defer fi.Close()
		scanner := bufio.NewScanner(fi)
		for scanner.Scan() {
			urls.Add(reg.FindAllString(scanner.Text(), -1)...)
		}
		errChan <- scanner.Err()
	}()
	select {
	case <-ctx.Done():
		return nil, nil
	case err := <-errChan:
		return urls, err
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
