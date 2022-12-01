package webdriver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Endpoint string
	Client   *http.Client
}

func (c *Client) do(req *http.Request, err error) func(f func(io.Reader) error) error {
	return func(f func(io.Reader) error) error {
		if err != nil {
			return err
		}
		hc := c.Client
		if hc == nil {
			hc = http.DefaultClient
		}
		resp, err := hc.Do(req)
		if err != nil {
			return err
		}
		defer func() {
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}()
		if resp.StatusCode != http.StatusOK {
			var e struct {
				Value struct {
					Error      string
					Message    string
					StackTrace string
				}
			}
			err = json.NewDecoder(resp.Body).Decode(&e)
			if err != nil {
				return err
			}
			return &Error{
				ErrorCode:  ErrorCode(e.Value.Error),
				Message:    e.Value.Message,
				StackTrace: strings.Split(e.Value.StackTrace, "\n"),
			}
		}
		if f == nil {
			return nil
		}
		return f(resp.Body)
	}
}

func (c *Client) url() string {
	return c.Endpoint
}

func (c *Client) Status(ctx context.Context) error {
	return c.do(http.NewRequestWithContext(ctx, http.MethodGet, c.url()+"/status", nil))(func(r io.Reader) error {
		var body struct {
			Value struct {
				Message string
				Ready   bool
			}
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		if body.Value.Ready {
			return nil
		}
		return fmt.Errorf("webdriver: %s", body.Value.Message)
	})
}

func (c *Client) Session(ctx context.Context) (s *session, err error) {
	err = c.do(http.NewRequestWithContext(ctx, http.MethodPost, c.url()+"/session", bytes.NewReader([]byte(`{"capabilities":{"firstMatch":[{"browserName":"firefox"},{"browserName":"chrome"}]}}`))))(func(r io.Reader) error {
		var body struct {
			Value struct {
				SessionID    string
				Capabilities map[string]interface{}
			}
		}
		err = json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		s = &session{
			ID:     body.Value.SessionID,
			client: c,
		}
		return nil
	})
	return s, err
}

func (c *Client) SessionWithCapabilities(ctx context.Context, body io.Reader) (s *session, err error) {
	err = c.do(http.NewRequestWithContext(ctx, http.MethodPost, c.url()+"/session", body))(func(r io.Reader) error {
		var body struct {
			Value struct {
				SessionID    string
				Capabilities map[string]interface{}
			}
		}
		err = json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		s = &session{
			ID:     body.Value.SessionID,
			client: c,
		}
		return nil
	})
	return s, err
}
