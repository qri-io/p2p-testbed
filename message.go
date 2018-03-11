package main

import (
	"github.com/opentracing/opentracing-go"
	"math/rand"
)

// MsgType indicates the type of message being sent
type MsgType string

// Message is a message to send from one peer to the other
type Message struct {
	ID      string
	Type    MsgType
	Tracing opentracing.TextMapCarrier
	Payload interface{}
}

// NewMessage creates a message
func NewMessage(t MsgType, payload interface{}, span opentracing.Span) Message {
	msg := Message{
		ID:      NewMessageID(),
		Type:    t,
		Payload: payload,
	}

	if span != nil {
		msg.Tracing = opentracing.TextMapCarrier{}
		span.Tracer().Inject(span.Context(), opentracing.TextMap, msg.Tracing)
	}

	return msg
}

var alpharunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// NewMessageID generates a random message identifier
func NewMessageID() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = alpharunes[rand.Intn(len(alpharunes))]
	}
	return string(b)
}
