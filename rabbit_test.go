package rabbit_ch_pool

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRabbit(t *testing.T) {
	t.Run("test rabbit reconnect", func(t *testing.T) {
		opt := &Options{
			Addr:"amqp://guest:guest@localhost:5673",
			ReconnectDelay:1,
		}
		rabbit := NewRabbit(opt)
		require.False(t, rabbit.connection.IsClosed())
		rabbit.connection.Close()
		require.True(t, rabbit.connection.IsClosed())
		time.Sleep(2 * time.Second)
		require.False(t, rabbit.connection.IsClosed())
	})
}
