package proton

import (
	"context"

	"github.com/go-resty/resty/v2"
)

func (c *Client) ListShares(ctx context.Context, all bool) ([]ShareMetadata, error) {
	var res struct {
		Shares []ShareMetadata
	}

	if err := c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
		if all {
			r.SetQueryParam("ShowAll", "1")
		}

		return r.SetResult(&res).Get("/drive/shares")
	}); err != nil {
		return nil, err
	}

	return res.Shares, nil
}

// ListSharedWithMeLinks returns the shares other users have shared with this user.
func (c *Client) ListSharedWithMeLinks(ctx context.Context) ([]SharedWithMeLink, error) {
	var links []SharedWithMeLink

	anchorID := ""

	for {
		// the v2 endpoint returns a native JSON bool for More, unlike the
		// 0/1 integers of the v4-era endpoints, so plain bool is used here
		var res struct {
			Links    []SharedWithMeLink
			AnchorID string
			More     bool
		}

		if err := c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
			if anchorID != "" {
				r.SetQueryParam("AnchorID", anchorID)
			}

			return r.SetResult(&res).Get("/drive/v2/sharedwithme")
		}); err != nil {
			return nil, err
		}

		links = append(links, res.Links...)

		if !res.More {
			break
		}

		anchorID = res.AnchorID
	}

	return links, nil
}

func (c *Client) GetShare(ctx context.Context, shareID string) (Share, error) {
	var res struct {
		Share
	}

	if err := c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
		return r.SetResult(&res).Get("/drive/shares/" + shareID)
	}); err != nil {
		return Share{}, err
	}

	return res.Share, nil
}
