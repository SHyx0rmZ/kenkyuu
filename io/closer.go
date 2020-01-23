package io

import (
	"github.com/SHyx0rmZ/kenkyuu/errors"
	"io"
)

type multiCloser struct {
	closers []io.Closer
}

func (mc multiCloser) Close() error {
	var me errors.MultiError
	for len(mc.closers) > 0 {
		if len(mc.closers) == 1 {
			if c, ok := mc.closers[0].(multiCloser); ok {
				mc.closers = c.closers
				continue
			}
		}
		e := mc.closers[0].Close()
		if e != nil {
			me.Errors = append(me.Errors, e)
		}
	}
	if len(me.Errors) == 0 {
		return nil
	}
	if len(me.Errors) == 1 {
		return me.Errors[0]
	}
	me.Errors = me.Errors[:len(me.Errors):len(me.Errors)]
	return me
}

func MultiCloser(closers ...io.Closer) io.Closer {
	c := make([]io.Closer, len(closers))
	copy(c, closers)
	return multiCloser{c}
}
