package webdriver

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
)

type Session interface {
	Close(ctx context.Context) error
	NavigateTo(ctx context.Context, url string) error
	GetCurrentURL(ctx context.Context) (url string, err error)
	Back(ctx context.Context) error
	Forward(ctx context.Context) error
	Refresh(ctx context.Context) error
	GetTitle(ctx context.Context) (title string, err error)
	GetWindowHandle(ctx context.Context) (handle string, err error)
	CloseWindow(ctx context.Context) error
	SwitchToWindow(ctx context.Context, handle string) error
	GetWindowHandles(ctx context.Context) (handles []string, err error)
	NewWindow(ctx context.Context) (handle, typ string, err error)
	SwitchToParentFrame(ctx context.Context) error
	GetWindowRect(ctx context.Context) (rectangle image.Rectangle, err error)
	SetWindowRect(ctx context.Context, rectangle image.Rectangle) error
	MaximizeWindow(ctx context.Context) error
	MinimizeWindow(ctx context.Context) error
	FindElement(ctx context.Context, strategy LocationStrategy, selector string) (element Element, err error)
	FindElements(ctx context.Context, strategy LocationStrategy, selector string) (elements []Element, err error)
	Screenshot(ctx context.Context) (image image.Image, err error)
	Actions(ctx context.Context, actions io.Reader) error
	AcceptAlert(ctx context.Context) error
}

type FirefoxSession interface {
	Session
	GetContext(ctx context.Context) (value string, err error)
	SetContext(ctx context.Context, value string) error
	InstallAddon(ctx context.Context, path string) (id string, err error)
	UninstallAddon(ctx context.Context, id string) error
}

type session struct {
	ID string
	Capabilities map[string]any

	client *Client
}

func (s *session) url() string {
	return s.client.url() + "/session/" + s.ID
}

func (s *session) Close(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodDelete, s.url(), nil))(func(r io.Reader) error {
		var body struct {
			Value struct {
				Error      string
				Message    string
				StackTrace string
			}
		}
		return json.NewDecoder(r).Decode(&body)
	})
}

func (s *session) NavigateTo(ctx context.Context, url string) error {
	body, err := json.Marshal(struct {
		URL string `json:"url"`
	}{
		URL: url,
	})
	if err != nil {
		return err
	}
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/url", bytes.NewReader(body)))(nil)
}

func (s *session) GetCurrentURL(ctx context.Context) (url string, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodGet, s.url()+"/url", nil))(func(r io.Reader) error {
		var body struct {
			Value string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		url = body.Value
		return nil
	})
	return url, err
}

func (s *session) Back(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/back", bytes.NewReader([]byte("{}"))))(nil)
}

func (s *session) Forward(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/forward", bytes.NewReader([]byte("{}"))))(nil)
}

func (s *session) Refresh(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/refresh", bytes.NewReader([]byte("{}"))))(nil)
}

func (s *session) GetTitle(ctx context.Context) (title string, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodGet, s.url()+"/title", nil))(func(r io.Reader) error {
		var body struct {
			Value string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		title = body.Value
		return nil
	})
	return title, err
}

func (s *session) GetWindowHandle(ctx context.Context) (handle string, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodGet, s.url()+"/window", nil))(func(r io.Reader) error {
		var body struct {
			Value string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		handle = body.Value
		return nil
	})
	return handle, err
}

func (s *session) CloseWindow(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodDelete, s.url()+"/window", bytes.NewReader([]byte("{}"))))(nil)
}

func (s *session) SwitchToWindow(ctx context.Context, handle string) error {
	body, err := json.Marshal(struct {
		Handle string `json:"handle"`
	}{
		Handle: handle,
	})
	if err != nil {
		return err
	}
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/window", bytes.NewReader(body)))(nil)
}

func (s *session) GetWindowHandles(ctx context.Context) (handles []string, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodGet, s.url()+"/window/handles", nil))(func(r io.Reader) error {
		var body struct {
			Value []string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		handles = body.Value
		return nil
	})
	return handles, err
}

func (s *session) NewWindow(ctx context.Context) (handle, typ string, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/window/new", bytes.NewReader([]byte("{}"))))(func(r io.Reader) error {
		var body struct {
			Value struct {
				Handle string `json:"handle"`
				Type   string `json:"type"`
			} `json:"value"`
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		handle = body.Value.Handle
		typ = body.Value.Type
		return nil
	})
	return handle, typ, err
}

func (s *session) SwitchToFrame(ctx context.Context, id string) error {
	body, err := json.Marshal(struct {
		ID struct {
			ID string `json:"element-6066-11e4-a52e-4f735466cecf"`
		} `json:"id"`
	}{
		ID: struct {
			ID string `json:"element-6066-11e4-a52e-4f735466cecf"`
		}{ID: id},
	})
	if err != nil {
		return nil
	}
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/frame", bytes.NewReader(body)))(nil)
}

func (s *session) SwitchToParentFrame(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/frame/parent", nil))(nil)
}

func (s *session) GetWindowRect(ctx context.Context) (rectangle image.Rectangle, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodGet, s.url()+"/window/rect", nil))(func(r io.Reader) error {
		var body struct {
			Value struct {
				Width  int
				Height int
				X      int
				Y      int
			}
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		rectangle = image.Rectangle{
			Min: image.Point{
				X: body.Value.X,
				Y: body.Value.Y,
			},
			Max: image.Point{
				X: body.Value.X + body.Value.Width,
				Y: body.Value.Y + body.Value.Height,
			},
		}
		return nil
	})
	return rectangle, err
}

func (s *session) SetWindowRect(ctx context.Context, rectangle image.Rectangle) error {
	body, err := json.Marshal(struct {
		Width  int `json:"width"`
		Height int `json:"height"`
		X      int `json:"x"`
		Y      int `json:"y"`
	}{
		Width:  rectangle.Dx(),
		Height: rectangle.Dy(),
		X:      rectangle.Min.X,
		Y:      rectangle.Min.Y,
	})
	if err != nil {
		return err
	}
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/window/rect", bytes.NewReader(body)))(nil)
}

func (s *session) MaximizeWindow(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/window/maximize", bytes.NewReader([]byte("{}"))))(nil)
}

func (s *session) MinimizeWindow(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/window/minimize", bytes.NewReader([]byte("{}"))))(nil)
}

type LocationStrategy string

const (
	CSSSelector             LocationStrategy = "css selector"
	LinkTextSelector        LocationStrategy = "link text"
	PartialLinkTextSelector LocationStrategy = "partial link text"
	TagName                 LocationStrategy = "tag name"
	XPathSelector           LocationStrategy = "xpath"
)

func (s *session) FindElement(ctx context.Context, strategy LocationStrategy, selector string) (e Element, err error) {
	body, err := json.Marshal(struct {
		Using string `json:"using"`
		Value string `json:"value"`
	}{
		Using: string(strategy),
		Value: selector,
	})
	if err != nil {
		return nil, err
	}
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/element", bytes.NewReader(body)))(func(r io.Reader) error {
		var body struct {
			Value map[string]string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		if len(body.Value) != 1 {
			return fmt.Errorf("webdriver: malformed result")
		}
		for _, v := range body.Value {
			e = &element{
				ID:      v,
				session: s,
			}
		}
		return nil
	})
	return e, err
}

func (s *session) FindElements(ctx context.Context, strategy LocationStrategy, selector string) (e []Element, err error) {
	body, err := json.Marshal(struct {
		Using string `json:"using"`
		Value string `json:"value"`
	}{
		Using: string(strategy),
		Value: selector,
	})
	if err != nil {
		return nil, err
	}
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/elements", bytes.NewReader(body)))(func(r io.Reader) error {
		var body struct {
			Value []map[string]string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		for i := range body.Value {
			if len(body.Value[i]) != 1 {
				e = nil
				return fmt.Errorf("webdriver: malformed result")
			}
			for _, v := range body.Value[i] {
				e = append(e, &element{
					ID:      v,
					session: s,
				})
			}
		}
		return nil
	})
	return e, err
}

func (s *session) Screenshot(ctx context.Context) (img image.Image, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodGet, s.url()+"/screenshot", nil))(func(r io.Reader) error {
		var body struct {
			Value string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		bs, err := base64.StdEncoding.DecodeString(body.Value)
		if err != nil {
			return err
		}
		img, _, err = image.Decode(bytes.NewReader(bs))
		if err != nil {
			return err
		}
		return nil
	})
	return img, err
}

func (s *session) Actions(ctx context.Context, actions io.Reader) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/actions", actions))(nil)
}

func (s *session) AcceptAlert(ctx context.Context) error {
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/alert/accept", bytes.NewReader([]byte("{}"))))(nil)
}

func (s *session) GetContext(ctx context.Context) (value string, err error) {
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodGet, s.url()+"/moz/context", nil))(func(r io.Reader) error {
		var body struct {
			Value string `json:"value"`
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		value = body.Value
		return nil
	})
	return value, err
}

func (s *session) SetContext(ctx context.Context, value string) error {
	body, err := json.Marshal(struct {
		Value string `json:"context"`
	}{
		Value: value,
	})
	if err != nil {
		return err
	}
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/moz/context", bytes.NewReader(body)))(nil)
}

func (s *session) InstallAddon(ctx context.Context, path string) (id string, err error) {
	body, err := json.Marshal(struct {
		Path      string `json:"path,omitempty"`
		Addon     string `json:"addon,omitempty"`
		Temporary bool   `json:"temporary"`
	}{
		Path: path,
	})
	if err != nil {
		return "", err
	}
	err = s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/moz/addon/install", bytes.NewReader(body)))(func(r io.Reader) error {
		var body struct {
			Value string `json:"value"`
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		id = body.Value
		return nil
	})
	return id, err
}

func (s *session) UninstallAddon(ctx context.Context, id string) error {
	body, err := json.Marshal(struct {
		ID string `json:"id"`
	}{
		ID: id,
	})
	if err != nil {
		return err
	}
	return s.client.do(http.NewRequestWithContext(ctx, http.MethodPost, s.url()+"/moz/addon/uninstall", bytes.NewReader(body)))(nil)
}
