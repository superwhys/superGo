package superMongo

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Iter struct {
	iter *mgo.Iter
	ctx  context.Context
}

func (iter *Iter) All(result interface{}) (err error) {
	if iter == nil {
		return ErrDial
	}
	if iter.ctx.Err() != nil {
		return iter.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	err = iter.iter.All(result)
	return err
}

func (iter *Iter) Close() error {
	if iter == nil {
		return ErrDial
	}
	err := iter.iter.Close()
	return err
}

func (iter *Iter) Done() bool {
	if iter == nil {
		return true
	}
	return iter.iter.Done()
}

func (iter *Iter) Err() error {
	if iter == nil {
		return ErrDial
	}
	return iter.iter.Err()
}

func (iter *Iter) For(result interface{}, f func() error) (err error) {
	if iter == nil {
		return ErrDial
	}
	if iter.ctx.Err() != nil {
		return iter.ctx.Err()
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return iter.iter.For(result, f)
}

func (iter *Iter) Next(result interface{}) bool {
	if iter == nil {
		return false
	}
	if iter.ctx.Err() != nil {
		return false
	}
	ret := iter.iter.Next(result)
	return ret
}

func (iter *Iter) State() (int64, []bson.Raw) {
	if iter == nil {
		return 0, nil
	}
	return iter.iter.State()
}

func (iter *Iter) Timeout() bool {
	if iter == nil {
		return false
	}
	return iter.iter.Timeout()
}
