package rabbit_ch_pool

import (
	"errors"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second
)

var (
	ErrorNotConnected = errors.New("failed: not connected")
)

// Rabbit service struct
type Rabbit struct {
	opt        *Options
	connection *amqp.Connection
	chPool     Pooler
	sigs       chan os.Signal
	isReady    bool
}

// NewRabbit create new svc
func NewRabbit(opt *Options) *Rabbit {
	r := &Rabbit{
		opt:  opt,
		sigs: make(chan os.Signal, 1),
	}

	if opt.ReconnectDelay == 0 {
		opt.ReconnectDelay = reconnectDelay
	}

	conn, _ := r.connect()
	r.connection = conn
	signal.Notify(r.sigs, syscall.SIGINT, syscall.SIGTERM)
	go r.handleReconnect()

	return r
}

func (r *Rabbit) handleReconnect() {
	for {
		if r.connection.IsClosed() {
			r.isReady = false

			_, err := r.connect()

			if err != nil {
				select {
				case <-r.sigs:
					return
				case <-time.After(r.opt.ReconnectDelay):
				}
				continue
			}
		}
		time.Sleep(r.opt.ReconnectDelay)
	}
}

func (r *Rabbit) connect() (*amqp.Connection, error) {
	r.isReady = false
	conn, err := amqp.Dial(r.opt.Addr)

	if err != nil {
		return nil, err
	}

	r.changeConnection(conn)
	r.isReady = true
	return conn, nil
}

func (r *Rabbit) changeConnection(connection *amqp.Connection) {
	r.connection = connection
	r.chPool = NewChannelPool(r.opt, r.connection)
}

// PublishMessage publish message to queue
func (r *Rabbit) PublishMessage(msg amqp.Publishing, routingKey string) error {
	if !r.isReady {
		return ErrorNotConnected
	}
	ch, err := r.chPool.Get()
	if err != nil {
		return err
	}

	err = ch.Publish(
		r.opt.ExchangeName,
		routingKey,
		false,
		false,
		msg,
	)
	if err != nil {
		return err
	}

	err = r.chPool.Put(ch)
	if err != nil {
		return err
	}

	return nil
}
