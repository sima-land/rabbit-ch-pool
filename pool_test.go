package rabbit_ch_pool

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChannelPool(t *testing.T) {
	t.Run("Success spawn pool", func(t *testing.T) {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5673")
		require.NoError(t, err)

		opt := &Options{
			PoolSize:5,
			PoolTimeout:1,
		}

		pool := NewChannelPool(opt, conn)
		require.Equal(t, 5, len(pool.chPool))

		ch, err := pool.Get()
		require.NoError(t, err)
		require.Equal(t, 4, len(pool.chPool))

		err = pool.Put(ch)
		require.NoError(t, err)
		require.Equal(t, 5, len(pool.chPool))
	})
	t.Run("Get channel, error after timeout", func(t *testing.T) {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5673")
		require.NoError(t, err)

		opt := &Options{
			PoolSize:5,
			PoolTimeout:1,
		}

		pool := NewChannelPool(opt, conn)
		for i := 0; i < 5; i++ {
			pool.Get()
		}
		require.Equal(t, 0, len(pool.chPool))
		_, err = pool.Get()
		require.Error(t, err)
		require.Equal(t, "Exception (504) Reason: \"channel/connection is not open\"", err.Error())
	})
}
