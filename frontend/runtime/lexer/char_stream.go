package lexer

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
)

var (
	OutOfBufferCapacity = errors.New("out of buffer capacity")
)

type CharStream interface {
	NextChar() (rune, error)
	Consume() []rune
	Rollback()
}

type cycleBufferCharStream struct {
	in       *bufio.Reader
	buf      []rune
	cap      int
	size     int
	readPos  int
	writePos int
	seekPos  int
	err      error
}

func (c *cycleBufferCharStream) NextChar() (rune, error) {
	if c.err != nil && c.err != io.EOF {
		return 0, c.err
	} else if c.err == io.EOF && c.writePos == c.readPos {
		return 0, c.err
	}
	c.sync()
	if c.err != nil && c.err != io.EOF {
		return 0, c.err
	}
	ret := c.buf[c.seekPos]
	c.seekPos = c.nextPos(c.seekPos)
	return ret, nil
}

func (c *cycleBufferCharStream) Consume() []rune {
	var ret []rune
	if c.isCycle() {
		ret = make([]rune, c.cap-c.readPos+c.seekPos)
		copy(ret, c.buf[c.readPos:])
		copy(ret[c.cap-c.readPos:], c.buf[:c.seekPos])
	} else {
		ret = make([]rune, c.seekPos-c.readPos)
		copy(ret, c.buf[c.readPos:c.seekPos])
	}
	c.readPos = c.seekPos
	c.size = c.size - len(ret)
	return ret
}

func (c *cycleBufferCharStream) Rollback() {
	if c.seekPos > c.readPos {
		c.seekPos--
	}
}

func (c *cycleBufferCharStream) sync() {
	if c.seekPos < c.size {
		return
	}
	if c.size >= c.cap {
		c.err = errors.Wrapf(OutOfBufferCapacity, "max capacity : %d", c.cap)
		return
	}
	var r rune
	r, _, c.err = c.in.ReadRune()
	if c.err != nil {
		return
	}
	c.buf[c.writePos] = r
	c.writePos = c.nextPos(c.writePos)
	c.size = c.size + 1
}

func (c *cycleBufferCharStream) nextPos(pos int) int {
	return (pos + 1) % c.cap
}

func (c *cycleBufferCharStream) isCycle() bool {
	return c.seekPos == c.readPos && c.size == c.cap || c.seekPos < c.readPos
}

func NewCycleCharStream(cap int, in io.Reader) CharStream {
	return &cycleBufferCharStream{
		in:       bufio.NewReader(in),
		buf:      make([]rune, cap),
		cap:      cap,
		size:     0,
		readPos:  0,
		writePos: 0,
		seekPos:  0,
	}
}
