package webdriver

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"io"
	"net/http"
)

type Element interface {
	GetElementAttribute(ctx context.Context, attribute string) (value string, err error)
	GetElementProperty(ctx context.Context, property string) (value string, err error)
	GetElementText(ctx context.Context) (text string, err error)
	Click(ctx context.Context) error
	SendKeys(ctx context.Context, keys string) error
	MoveTo(ctx context.Context) error
}

type element struct {
	ID string

	session *session
}

func (e *element) url() string {
	return e.session.url() + "/element/" + e.ID
}

func (e *element) GetElementAttribute(ctx context.Context, attribute string) (value string, err error) {
	err = e.session.client.do(http.NewRequestWithContext(ctx, http.MethodGet, e.url()+"/attribute/"+attribute, nil))(func(r io.Reader) error {
		var body struct {
			Value string
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

func (e *element) GetElementProperty(ctx context.Context, property string) (value string, err error) {
	err = e.session.client.do(http.NewRequestWithContext(ctx, http.MethodGet, e.url()+"/property/"+property, nil))(func(r io.Reader) error {
		var body struct {
			Value string
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

func (e *element) GetElementText(ctx context.Context) (text string, err error) {
	err = e.session.client.do(http.NewRequestWithContext(ctx, http.MethodGet, e.url()+"/text", nil))(func(r io.Reader) error {
		var body struct {
			Value string
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		text = body.Value
		return nil
	})
	return text, err
}

func (e *element) GetElementRect(ctx context.Context) (rectangle image.Rectangle, err error) {
	err = e.session.client.do(http.NewRequestWithContext(ctx, http.MethodGet, e.url()+"/rect", nil))(func(r io.Reader) error {
		var body struct {
			Value struct {
				Width  float32
				Height float32
				X      float32
				Y      float32
			}
		}
		err := json.NewDecoder(r).Decode(&body)
		if err != nil {
			return err
		}
		rectangle = image.Rectangle{
			Min: image.Point{
				X: int(body.Value.X),
				Y: int(body.Value.Y),
			},
			Max: image.Point{
				X: int(body.Value.X + body.Value.Width),
				Y: int(body.Value.Y + body.Value.Height),
			},
		}
		return nil
	})
	return rectangle, err
}

func (e *element) Click(ctx context.Context) error {
	return e.session.client.do(http.NewRequestWithContext(ctx, http.MethodPost, e.url()+"/click", bytes.NewReader([]byte("{}"))))(nil)
}

func (e *element) SendKeys(ctx context.Context, keys string) error {
	body, err := json.Marshal(struct {
		Text string `json:"text"`
	}{
		Text: keys,
	})
	if err != nil {
		return err
	}
	return e.session.client.do(http.NewRequestWithContext(ctx, http.MethodPost, e.url()+"/value", bytes.NewReader(body)))(nil)
}

func (e *element) MoveTo(ctx context.Context) error {
	rect, err := e.GetElementRect(ctx)
	if err != nil {
		return err
	}

	type _action struct {
		Type     string `json:"type"`
		Duration uint   `json:"duration"`
		X        int    `json:"x"`
		Y        int    `json:"y"`
	}
	type action struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Parameters struct {
			PointerType string `json:"pointerType"`
		} `json:"parameters"`
		Actions []_action `json:"actions"`
	}
	body, err := json.Marshal(struct {
		Actions []action `json:"actions"`
	}{
		Actions: []action{
			{
				Type: "pointer",
				Parameters: struct {
					PointerType string `json:"pointerType"`
				}{PointerType: "mouse"},
				Actions: []_action{
					{
						Type:     "pointerMove",
						Duration: 1000,
						X:        rect.Min.X,
						Y:        rect.Min.Y,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return e.session.Actions(ctx, bytes.NewReader(body))
}

func (e *element) String() string {
	return e.ID
}
