package superMongo

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	lg "github.com/superwhys/superGo/superLog"
)

const defaultPoolSize = 100
const defaultMode = mgo.Eventual
const defaultTimeout = time.Second * 10

// Client ...
type Client struct {
	session   *mgo.Session
	worker    chan int
	url       string
	addresses []string
	user      string
	password  string

	mode             mgo.Mode
	batchSize        int
	ensureReplicated bool
	poolSize         int
	timeout          time.Duration
	lock             sync.RWMutex
	cntLock          sync.Mutex
	connCounter      map[string]int
}

type Option func(*Client)

func WithAuth(user, password string) Option {
	return func(c *Client) {
		c.user = user
		c.password = password
	}
}

func WithMode(mode mgo.Mode) Option {
	return func(c *Client) {
		c.mode = mode
	}
}

func WithBatch(n int) Option {
	return func(c *Client) {
		c.batchSize = n
	}
}

func WithEnsureReplicated() Option {
	return func(c *Client) {
		c.ensureReplicated = true
	}
}

func WithPoolSize(size int) Option {
	return func(c *Client) {
		c.poolSize = size
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// NewClient returns a `Client` object.
func NewClient(url string, opts ...Option) *Client {
	c := &Client{
		url:         url,
		mode:        defaultMode,
		poolSize:    defaultPoolSize,
		timeout:     defaultTimeout,
		connCounter: make(map[string]int),
	}
	for _, opt := range opts {
		opt(c)
	}
	c.worker = make(chan int, c.poolSize)
	return c
}

// CollectionNames returns all the collection names of a db.
func (c *Client) CollectionNames(db string) ([]string, error) {
	session, err := c.newSession()
	if err != nil {
		lg.Error("Failed to dial mongodb", c.url, db, err)
		return nil, err
	}
	connKey := fmt.Sprintf("collection-names.%s", db)
	c.incConnectionCount(connKey)
	defer c.decConnectionCount(connKey)
	defer session.Close()
	return session.DB(db).CollectionNames()
}

func (c *Client) doDial() error {
	addrURL := c.url
	if len(c.addresses) > 0 || c.user != "" {
		// Construct the dsn instead of given one.
		var u url.URL
		u.Scheme = "mongodb"
		if c.user != "" && c.password != "" {
			u.User = url.UserPassword(c.user, c.password)
		} else if c.user != "" {
			u.User = url.User(c.user)
		}
		u.RawQuery = "connect=direct"
		u.Host = strings.Join(c.addresses, ",")
		addrURL = u.String()
	}
	session, err := mgo.Dial(addrURL)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.session = session
	c.refreshSetting()
	return nil
}

// Open creates a connection to the specified db+collection.
// Deprecated. Use OpenWithContext() instead.
func (c *Client) Open(db, collection string) *Collection {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	return c.OpenWithContext(ctx, db, collection)
}

// OpenWithContext creates a connection to the specified db+collection.
func (c *Client) OpenWithContext(ctx context.Context, db, collection string) *Collection {
	session, err := c.newSession()
	if err != nil {
		lg.Error("Failed to dial mongodb", c.url, db, collection, err)
		return nil
	}
	connKey := fmt.Sprintf("%s.%s", db, collection)
	c.incConnectionCount(connKey)

	colle := session.DB(db).C(collection)
	newColle := &Collection{
		colle:   colle,
		session: session,
		client:  c,
		connKey: connKey,
	}
	newColle.SetContext(ctx)
	return newColle
}

func (c *Client) incConnectionCount(connKey string) {
	c.cntLock.Lock()
	defer c.cntLock.Unlock()

	c.connCounter[connKey]++
}

func (c *Client) decConnectionCount(connKey string) {
	c.cntLock.Lock()
	defer c.cntLock.Unlock()

	c.connCounter[connKey]--
}

func (c *Client) statConnections() map[string]int {
	c.cntLock.Lock()
	defer c.cntLock.Unlock()

	ret := map[string]int{}
	for k, v := range c.connCounter {
		if v > 0 {
			ret[k] = v
		}
	}
	return ret
}

// NewSession returns a session for mongodb operation.
func (c *Client) newSession() (*mgo.Session, error) {
	c.lock.RLock()
	sessionNull := c.session == nil
	c.lock.RUnlock()

	if sessionNull {
		if err := c.doDial(); err != nil {
			return nil, err
		}
	}
	// Limit the total worker.
	select {
	case c.worker <- 1:
	default:
		lg.Error("Mongo worker has exceed mongo pool size (", c.poolSize, "), please check if any mongo connection is leaking.\n", lg.Jsonify(c.statConnections()))
		c.worker <- 1
	}
	return c.session.Copy(), nil
}

// OpenWithLongTimeout creates a connection to the specified db+collection.
// Deprecated. Use OpenWithContext() instead.
func (c *Client) OpenWithLongTimeout(db, collection string) *Collection {
	timeout := time.Hour * 24
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	colle := c.OpenWithContext(ctx, db, collection)
	colle.session.SetSocketTimeout(timeout)
	colle.session.SetCursorTimeout(0)
	return colle
}

func (c *Client) refreshSetting() {
	if c.session == nil {
		return
	}
	c.session.SetMode(c.mode, true)
	if c.batchSize > 0 {
		c.session.SetBatch(c.batchSize)
	}
	if c.ensureReplicated {
		c.session.SetSafe(&mgo.Safe{W: len(c.addresses)})
	} else {
		c.session.SetSafe(&mgo.Safe{W: 1})
	}
}

// SetBatch sets the maximum number of entries return from a single query.
// Helpful to reduce latency if the payload size is small and quantity is large.
func (c *Client) SetBatch(n int) {
	c.batchSize = n
	c.refreshSetting()
}

// SetMode sets the underlying consistency mode in the momgodb client.
func (c *Client) SetMode(mode mgo.Mode) {
	c.mode = mode
	c.refreshSetting()
}

// EnsureReplicated make sure all data written are replicated to all secondary servers.
func (c *Client) EnsureReplicated() {
	c.ensureReplicated = true
	c.refreshSetting()
}
