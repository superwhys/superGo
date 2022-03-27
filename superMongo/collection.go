package superMongo

import (
	"context"
	"fmt"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// ErrDial returns dialing error.
var ErrDial = errors.New("Failed to dial to mongodb")

// Collection is the mongodb collection to operate on.
//
// Origin mgo library will throw panic error on sometimes, such as query on closed session.
// In the library, we convert the mgo panic error into normal error to avoid program crash.
// Maybe there are some problems, so you should handle mongo err properly by your situation.
// Each collection should invoke `.Close()` as soon as the op has finished.
type Collection struct {
	ctx       context.Context
	colle     *mgo.Collection
	session   *mgo.Session
	client    *Client
	closeOnce sync.Once
	connKey   string
}

// Close session
func (colle *Collection) Close() {
	if colle == nil {
		return
	}
	colle.closeOnce.Do(func() {
		colle.session.Close()
		<-colle.client.worker
		colle.client.decConnectionCount(colle.connKey)
		// Remove reference for GC
		colle.client = nil
		colle.session = nil
		colle.colle = nil
	})
}

// WithContext horrors input context, and close session when context is done.
// Deprecated. Use SetContext() instead.
func (colle *Collection) WithContext(ctx context.Context) *Collection {
	if colle == nil {
		return nil
	}
	return colle.SetContext(ctx)
}

// SetContext horors input context, and close session when context is done.
func (colle *Collection) SetContext(ctx context.Context) *Collection {
	if colle == nil {
		return nil
	}
	colle.ctx = ctx
	return colle
}

// SetBatch sets the maximum number of entries return from a single query.
// Helpful to reduce latency if the payload size is small and quantity is large.
func (colle *Collection) SetBatch(n int) *Collection {
	if colle == nil {
		return nil
	}
	colle.session.SetBatch(n)
	return colle
}

// SetMode sets the underlying consistency mode in the momgodb collection.
func (colle *Collection) SetMode(mode mgo.Mode) *Collection {
	if colle == nil {
		return nil
	}
	colle.session.SetMode(mode, true)
	return colle
}

// EnsureReplicated make sure all data written are replicated to all secondary servers.
func (colle *Collection) EnsureReplicated() *Collection {
	if colle == nil {
		return nil
	}
	colle.session.SetSafe(&mgo.Safe{W: len(colle.client.addresses)})
	return colle
}

func (colle *Collection) Count() (n int, err error) {
	if colle == nil {
		return 0, ErrDial
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	n, err = colle.colle.Count()
	return
}

func (colle *Collection) Create(info *mgo.CollectionInfo) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.Create(info)
	return
}

func (colle *Collection) DropAllIndexes() (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	err = colle.colle.DropAllIndexes()
	return
}

func (colle *Collection) DropCollection() (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	err = colle.colle.DropCollection()
	return
}

func (colle *Collection) DropIndex(key ...string) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	err = colle.colle.DropIndex(key...)
	return
}

func (colle *Collection) DropIndexName(name string) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	err = colle.colle.DropIndexName(name)
	return
}

func (colle *Collection) EnsureIndex(index mgo.Index) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.EnsureIndex(index)
	return
}

func (colle *Collection) EnsureIndexKey(key ...string) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.EnsureIndexKey(key...)
	return
}

func (colle *Collection) Find(query interface{}) *Query {
	if colle == nil {
		return nil
	}
	tmpQuery := colle.colle.Find(query)
	return &Query{tmpQuery, colle.ctx}
}

func (colle *Collection) FindId(id interface{}) *Query {
	if colle == nil {
		return nil
	}
	tmpQuery := colle.colle.FindId(id)
	return &Query{tmpQuery, colle.ctx}
}

func (colle *Collection) Indexes() (indexes []mgo.Index, err error) {
	if colle == nil {
		return nil, ErrDial
	}
	if colle.ctx.Err() != nil {
		return nil, colle.ctx.Err()
	}
	return colle.colle.Indexes()
}

func (colle *Collection) Insert(docs ...interface{}) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.Insert(docs...)
	return
}

func (colle *Collection) NewIter(session *mgo.Session, firstBatch []bson.Raw, cursorId int64, err error) (iter *Iter) {
	if colle == nil {
		return nil
	}
	tmpIter := colle.colle.NewIter(session, firstBatch, cursorId, err)
	iter = &Iter{tmpIter, colle.ctx}
	return
}

func (colle *Collection) Pipe(pipeline interface{}) (pipe *mgo.Pipe) {
	if colle == nil {
		return nil
	}
	pipe = colle.colle.Pipe(pipeline)
	return
}

func (colle *Collection) Remove(selector interface{}) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.Remove(selector)
	return
}

func (colle *Collection) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	if colle == nil {
		return nil, ErrDial
	}
	if colle.ctx.Err() != nil {
		return nil, colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	info, err = colle.colle.RemoveAll(selector)
	return
}

func (colle *Collection) RemoveId(id interface{}) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.RemoveId(id)
	return
}

func (colle *Collection) Repair() (iter *Iter) {
	if colle == nil {
		return nil
	}
	tmpIter := colle.colle.Repair()
	iter = &Iter{tmpIter, colle.ctx}
	return
}

func (colle *Collection) Update(selector interface{}, update interface{}) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.Update(selector, update)
	return
}

func (colle *Collection) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	if colle == nil {
		return nil, ErrDial
	}
	if colle.ctx.Err() != nil {
		return nil, colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	info, err = colle.colle.UpdateAll(selector, update)
	return
}

func (colle *Collection) UpdateId(id interface{}, update interface{}) (err error) {
	if colle == nil {
		return ErrDial
	}
	if colle.ctx.Err() != nil {
		return colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = colle.colle.UpdateId(id, update)
	return
}

func (colle *Collection) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	if colle == nil {
		return nil, ErrDial
	}
	if colle.ctx.Err() != nil {
		return nil, colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	info, err = colle.colle.Upsert(selector, update)
	return
}

func (colle *Collection) UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	if colle == nil {
		return nil, ErrDial
	}
	if colle.ctx.Err() != nil {
		return nil, colle.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	info, err = colle.colle.UpsertId(id, update)
	return
}

func (colle *Collection) With(s *mgo.Session) (conn *mgo.Collection) {
	if colle == nil {
		return nil
	}
	conn = colle.colle.With(s)
	return
}

func (colle *Collection) Bulk() *mgo.Bulk {
	if colle == nil {
		return nil
	}
	return colle.colle.Bulk()
}
