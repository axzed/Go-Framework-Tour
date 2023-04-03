package demo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	ch := make(chan string, 4)
	go func() {
		str := <-ch
		fmt.Println(str)
	}()
	go func() {
		str := <-ch
		fmt.Println(str)
	}()
	go func() {
		str := <-ch
		fmt.Println(str)
	}()

	ch <- "hello"
	ch <- "hello"
	time.Sleep(time.Second)
}

func TestBroker(t *testing.T) {
	b := &Broker{
		consumers: make([]*Consumer, 0, 10),
	}
	c1 := &Consumer{
		ch: make(chan string, 1),
	}
	c2 := &Consumer{
		ch: make(chan string, 1),
	}
	b.Subscribe(c1)
	b.Subscribe(c2)

	b.Produce("hello")
	fmt.Println(<-c1.ch)
	fmt.Println(<-c2.ch)
}

type Broker struct {
	consumers []*Consumer
}

func (b *Broker) Produce(msg string) {
	for _, c := range b.consumers {
		c.ch <- msg
	}
}

func (b *Broker) Subscribe(c *Consumer) {
	b.consumers = append(b.consumers, c)
}

type Consumer struct {
	ch chan string
}

type Broker1 struct {
	ch        chan string
	consumers []func(s string)
}

func (b *Broker1) Produce(msg string) {
	b.ch <- msg
}

func (b *Broker1) Subscribe(consume func(s string)) {
	b.consumers = append(b.consumers, consume)
}

func (b *Broker1) Start() {
	go func() {
		for {
			s, ok := <-b.ch
			if !ok {
				return
			}
			for _, c := range b.consumers {
				c(s)
			}
		}
	}()
}

func NewBroker1() *Broker1 {
	b := &Broker1{ch: make(chan string, 10), consumers: make([]func(s string), 0, 10)}
	go func() {
		for {
			s, ok := <-b.ch
			if !ok {
				return
			}
			for _, c := range b.consumers {
				c(s)
			}
		}
	}()
	return b
}

func TestBroker1(t *testing.T) {
	b := NewBroker1()
	str1 := ""
	b.Subscribe(func(s string) {
		str1 = str1 + s
	})

	str2 := ""
	b.Subscribe(func(s string) {
		str2 = str2 + s
	})

	b.Produce("hello")
	b.Produce(" ")
	b.Produce("world")

	time.Sleep(time.Second)
	assert.Equal(t, "hello world", str1)
	assert.Equal(t, "hello world", str2)
}
