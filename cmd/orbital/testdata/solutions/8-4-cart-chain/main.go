package main

import "errors"

var ErrCoupling = errors.New("invalid coupling position")

type Cart struct {
	Tag  string
	Next *Cart
}
type Convoy struct{ Head *Cart }

func (c *Convoy) Couple(at int, tag string) error {
	if at < 0 {
		return ErrCoupling
	}
	if at == 0 {
		c.Head = coupledCart(c.Head, tag)
		return nil
	}
	prev := c.Head
	for i := 1; i < at && prev != nil; i++ {
		prev = prev.Next
	}
	if prev == nil {
		return ErrCoupling
	}
	prev.Next = coupledCart(prev.Next, tag)
	return nil
}
func (c *Convoy) Release(tag string) bool {
	if c.Head == nil {
		return false
	}
	if c.Head.Tag == tag {
		c.Head = c.Head.Next
		return true
	}
	for prev := c.Head; prev.Next != nil; prev = prev.Next {
		if prev.Next.Tag == tag {
			prev.Next = prev.Next.Next
			return true
		}
	}
	return false
}
func (c *Convoy) Locate(tag string) int {
	at := 0
	for cart := c.Head; cart != nil; cart = cart.Next {
		if cart.Tag == tag {
			return at
		}
		at++
	}
	return -1
}
func (c *Convoy) Tags() []string {
	var tags []string
	for cart := c.Head; cart != nil; cart = cart.Next {
		tags = append(tags, cart.Tag)
	}
	return tags
}
func main() {}

func coupledCart(next *Cart, tag string) *Cart { return &Cart{Tag: tag, Next: next} }
