package superMongo

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
)

type Query struct {
	queryProxy *mgo.Query
	ctx        context.Context
}

func (q *Query) All(result interface{}) (err error) {
	if q == nil {
		return ErrDial
	}
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return q.queryProxy.All(result)
}

func (q *Query) For(result interface{}, f func() error) (err error) {
	if q == nil {
		return ErrDial
	}
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return q.queryProxy.For(result, f)
}

func (q *Query) Count() (n int, err error) {
	if q == nil {
		return 0, ErrDial
	}
	if q.ctx.Err() != nil {
		return 0, q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	n, err = q.queryProxy.Count()
	return
}

func (q *Query) Distinct(key string, result interface{}) (err error) {
	if q == nil {
		return ErrDial
	}
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return q.queryProxy.Distinct(key, result)
}

func (q *Query) Apply(change mgo.Change, result interface{}) (info *mgo.ChangeInfo, err error) {
	if q == nil {
		return nil, ErrDial
	}
	if q.ctx.Err() != nil {
		return nil, q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	info, err = q.queryProxy.Apply(change, result)
	return
}

func (q *Query) Batch(n int) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Batch(n)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Collation(collation *mgo.Collation) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Collation(collation)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Comment(comment string) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Comment(comment)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Explain(result interface{}) (err error) {
	if q == nil {
		return ErrDial
	}
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = q.queryProxy.Explain(result)
	return err
}

func (q *Query) Hint(indexKey ...string) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Hint(indexKey...)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Iter() *Iter {
	if q == nil {
		return nil
	}
	tmpIter := q.queryProxy.Iter()
	iter := &Iter{tmpIter, q.ctx}
	return iter
}

func (q *Query) Limit(n int) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Limit(n)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) LogReplay() *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.LogReplay()
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) MapReduce(job *mgo.MapReduce, result interface{}) (info *mgo.MapReduceInfo, err error) {
	if q == nil {
		return nil, ErrDial
	}
	if q.ctx.Err() != nil {
		return nil, q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	info, err = q.queryProxy.MapReduce(job, result)
	return
}

func (q *Query) One(result interface{}) (err error) {
	if q == nil {
		return ErrDial
	}
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = q.queryProxy.One(result)
	return
}

func (q *Query) Prefetch(p float64) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Prefetch(p)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Select(selector interface{}) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Select(selector)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) SetMaxScan(n int) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.SetMaxScan(n)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) SetMaxTime(d time.Duration) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.SetMaxTime(d)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Skip(n int) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Skip(n)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Snapshot() *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Snapshot()
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Sort(fields ...string) *Query {
	if q == nil {
		return nil
	}
	tmpQuery := q.queryProxy.Sort(fields...)
	query := &Query{queryProxy: tmpQuery, ctx: q.ctx}
	return query
}

func (q *Query) Tail(timeout time.Duration) *mgo.Iter {
	if q == nil {
		return nil
	}
	iter := q.queryProxy.Tail(timeout)
	return iter
}
