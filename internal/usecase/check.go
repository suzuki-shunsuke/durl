package usecase

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"gopkg.in/yaml.v2"

	"mvdan.cc/xurls"

	"github.com/pkg/errors"
	"github.com/scylladb/go-set/strset"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

// Check checks whether broken urls are included in files.
func Check(fsys domain.Fsys, stdin io.Reader, cfgPath string) error {
	cfg, err := readCfg(fsys, cfgPath)
	if err != nil {
		return err
	}

	files, err := getFiles(stdin)
	if err != nil {
		return err
	}
	urls, err := extractURLsFromFiles(fsys, files)
	if err != nil {
		return err
	}
	for _, u := range cfg.IgnoreURLs {
		delete(urls, u)
	}
	return checkURLs(urls)
}

func findCfg(fsys domain.Fsys) (string, error) {
	wd, err := fsys.Getwd()
	if err != nil {
		return "", err
	}
	for {
		p := filepath.Join(wd, ".durl.yml")
		if fsys.Exist(p) {
			return p, nil
		}
		if wd == "/" || wd == "" {
			return "", fmt.Errorf(".durl.yml is not found")
		}
		wd = filepath.Dir(wd)
	}
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
		return fmt.Errorf("%s is broken (%d)", u, resp.StatusCode)
	}
	return nil
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
