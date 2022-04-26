package notify

import "context"

type Notify interface {
	Send(ctx context.Context, content string) error
}

type combine struct {
	notify []Notify
}

func Combine(notify ...Notify) Notify {
	return &combine{
		notify: notify,
	}
}

func (c *combine) Send(ctx context.Context, content string) error {
	for _, n := range c.notify {
		err := n.Send(ctx, content)
		if err != nil {
			return err
		}
	}

	return nil
}
