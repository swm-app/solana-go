// Copyright 2021 github.com/gagliardetto
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

type RootResult uint64

// RootSubscribe subscribes to receive notification
// anytime a new root is set by the validator.
func (cl *Client) RootSubscribe() (*Subscription[RootResult], error) {
	genSub, err := cl.subscribe(
		nil,
		nil,
		"rootSubscribe",
		"rootUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res RootResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &Subscription[RootResult]{
		sub:       genSub,
		closeFunc: func() { genSub.closeFunc(nil) },
	}, nil
}
