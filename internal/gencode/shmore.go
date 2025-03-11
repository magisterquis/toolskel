package gencode

/*
 * shmore.go
 * "Generate" the latest shmore
 * By J. Stuart McMurray
 * Created 20250310
 * Last Modified 20250310
 */

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/net/context/ctxhttp"
)

// ShmoreURL is the URL for the latest version of Shmore
const ShmoreURL = "https://raw.githubusercontent.com/magisterquis/shmore/refs/heads/master/shmore.subr"

// Shmore downloads the latest Shmore.  The context is there not only to leave
// room for future cancellation but also to make this method not a Generator.
func (p Params) Shmore(ctx context.Context) ([]byte, error) {
	/* Request Shmore. */
	res, err := ctxhttp.Get(ctx, nil, ShmoreURL)
	if nil != err {
		return nil, fmt.Errorf(
			"requesting Shmore from %s: %s",
			ShmoreURL,
			err,
		)
	}
	defer res.Body.Close()

	/* Read the library from the response. */
	b, err := io.ReadAll(res.Body)
	if nil != err {
		return nil, fmt.Errorf("reading HTTP response body: %w", err)
	}

	return b, nil
}
