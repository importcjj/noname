package notify

import (
	"context"
	"log"
)

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
	log.Println(content)
	for _, n := range c.notify {
		err := n.Send(ctx, content)
		if err != nil {
			log.Printf("failed to notify %v", err)
			return err
		}
	}

	return nil
}
