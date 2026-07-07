// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ws

import (
	"context"
	"sync"
)

type result interface{}

type decoderFunc func([]byte) (interface{}, error)

type subscription struct {
	req               *request
	subID             uint64
	stream            chan result
	done              chan struct{}
	once              sync.Once
	ack               chan struct{}
	ackOnce           sync.Once
	err               error
	closeFunc         func(err error)
	unsubscribeMethod string
	decoderFunc       decoderFunc
}

func newSubscription(
	req *request,
	closeFunc func(err error),
	unsubscribeMethod string,
	decoderFunc decoderFunc,
) *subscription {
	return &subscription{
		req:               req,
		stream:            make(chan result, 256),
		done:              make(chan struct{}),
		ack:               make(chan struct{}),
		closeFunc:         closeFunc,
		unsubscribeMethod: unsubscribeMethod,
		decoderFunc:       decoderFunc,
	}
}

// confirm records the server's acknowledgement of the subscription request.
// Idempotent: a duplicate acknowledgement for the same request is a no-op.
func (s *subscription) confirm() {
	s.ackOnce.Do(func() {
		close(s.ack)
	})
}

func (s *subscription) close(err error) {
	s.once.Do(func() {
		s.err = err
		close(s.done)
	})
}

func (s *subscription) isDone() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}

// Subscription is a type-safe generic subscription.
type Subscription[T any] struct {
	sub       *subscription
	closeFunc func()
}

func (s *Subscription[T]) Recv(ctx context.Context) (*T, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.sub.done:
		if s.sub.err != nil {
			return nil, s.sub.err
		}
		return nil, ErrSubscriptionClosed
	case d := <-s.sub.stream:
		return d.(*T), nil
	}
}

func (s *Subscription[T]) Unsubscribe() {
	s.closeFunc()
}

// WaitConfirmed blocks until the server responds to the subscription request.
// The subscribe call itself only writes the request to the socket; the server's
// reply — a subscription id on acceptance, or an error object on rejection —
// arrives asynchronously on the read loop. WaitConfirmed exposes that reply:
// it returns nil once the server acknowledges the subscription, the server's
// error if it rejects it (or the connection fails first), or ctx.Err() when
// ctx is done. A subscription closed after being acknowledged still reports
// nil — acceptance, once observed, is not retracted.
func (s *Subscription[T]) WaitConfirmed(ctx context.Context) error {
	select {
	case <-s.sub.ack:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-s.sub.done:
		// The acknowledgement and a later close may both have landed before this
		// select ran; acceptance wins over the subsequent close.
		select {
		case <-s.sub.ack:
			return nil
		default:
		}
		if s.sub.err != nil {
			return s.sub.err
		}
		return ErrSubscriptionClosed
	}
}
