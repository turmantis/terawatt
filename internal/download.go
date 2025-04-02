package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func DownloadTo(ctx context.Context, url string, dest io.Writer) error {
	fmt.Println("Downloading:", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			panic(err)
		}
	}()
	if res.StatusCode >= 400 {
		return fmt.Errorf("internal: sad status code :( %d", res.StatusCode)
	}
	if _, err = io.Copy(dest, res.Body); err != nil {
		return fmt.Errorf("internal: %w", err)
	}
	return nil
}
