package ws

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// newTestSubscription registers a fresh subscription on a client that has no
// live connection, so a test can drive the server's reply by feeding raw
// messages to handleMessage directly.
func newTestSubscription(c *Client) *Subscription[AccountResult] {
	genSub := newSubscription(
		newRequest(nil, "accountSubscribe", nil, false),
		func(err error) {},
		"accountUnsubscribe",
		func(msg []byte) (interface{}, error) { return nil, nil },
	)
	c.subscriptionByRequestID[genSub.req.ID] = genSub
	return &Subscription[AccountResult]{
		sub:       genSub,
		closeFunc: func() { genSub.closeFunc(nil) },
	}
}

func newTestClient() *Client {
	return &Client{
		subscriptionByRequestID: map[uint64]*subscription{},
		subscriptionByWSSubID:   map[uint64]*subscription{},
	}
}

func TestWaitConfirmedBlocksUntilServerReply(t *testing.T) {
	c := newTestClient()
	sub := newTestSubscription(c)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, sub.WaitConfirmed(ctx), context.DeadlineExceeded,
		"WaitConfirmed must block while the server has not replied")
}

func TestWaitConfirmedOnServerAck(t *testing.T) {
	c := newTestClient()
	sub := newTestSubscription(c)

	c.handleMessage(fmt.Appendf(nil, `{"jsonrpc":"2.0","result":42,"id":%d}`, sub.sub.req.ID))

	require.NoError(t, sub.WaitConfirmed(context.Background()))
	require.Equal(t, uint64(42), sub.sub.subID, "the ack must register the server-assigned subscription id")
}

func TestWaitConfirmedOnServerRejection(t *testing.T) {
	c := newTestClient()
	sub := newTestSubscription(c)

	c.handleMessage(fmt.Appendf(nil,
		`{"jsonrpc":"2.0","error":{"code":-32601,"message":"method not found"},"id":%d}`, sub.sub.req.ID))

	err := sub.WaitConfirmed(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "method not found", "the server's rejection reason must surface")
}

func TestWaitConfirmedAckSurvivesLaterClose(t *testing.T) {
	c := newTestClient()
	sub := newTestSubscription(c)

	c.handleMessage(fmt.Appendf(nil, `{"jsonrpc":"2.0","result":7,"id":%d}`, sub.sub.req.ID))
	sub.sub.close(fmt.Errorf("connection reset"))

	require.NoError(t, sub.WaitConfirmed(context.Background()),
		"an acknowledged subscription reports acceptance even after it later closes")
}

func TestCloseSubscriptionBeforeAckSkipsUnsubscribe(t *testing.T) {
	c := newTestClient()
	sub := newTestSubscription(c)

	// The test client has no live connection: if closing an unacknowledged
	// subscription attempted the unsubscribe write, it would dereference the
	// nil socket and panic. The subID guard must skip the write entirely.
	c.closeSubscription(sub.sub.req.ID, fmt.Errorf("ack timeout"))

	err := sub.WaitConfirmed(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "ack timeout")
}

func TestWaitConfirmedConnectionFailureBeforeAck(t *testing.T) {
	c := newTestClient()
	sub := newTestSubscription(c)

	c.closeAllSubscription(fmt.Errorf("connection reset"))

	err := sub.WaitConfirmed(context.Background())
	require.Error(t, err, "a connection failure before the ack must surface as the confirmation error")
	require.Contains(t, err.Error(), "connection reset")
}
