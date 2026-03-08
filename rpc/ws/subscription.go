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
		closeFunc:         closeFunc,
		unsubscribeMethod: unsubscribeMethod,
		decoderFunc:       decoderFunc,
	}
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
