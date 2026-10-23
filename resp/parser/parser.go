// parser 包用于解析 RESP 协议数据
package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

type readState struct {
	// * 表示多条数据
	// $ 表示单条数据
	// + 表示错误
	// : 表示整数
	msgType           byte
	readingMultiLine  bool
	expectedArgsCount int
	args              [][]byte
	bulkLen           int64
}

func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && s.expectedArgsCount == len(s.args)
}

type Payload struct {
	Data resp.Reply
	Err  error
}

func ParserTream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)

	go parse0(reader, ch)

	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {}

// readLine 读取一行数据
// 1. 没有读取到 $7 按照 \r\n 切割
// 2. 读取到了 $7 严格读取字符个数
// *1\r\n$3\r\nSET\r\n
func readLine(bufReader *bufio.Reader, state *readState) (msg []byte, ok bool, err error) {
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}

		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, fmt.Errorf("protocol error: %s", string(msg))
		}

	}

	// 还需要读取后面的 \r\n
	msg, err = bufReader.Peek(int(state.bulkLen) + 2)
	if err != nil {
		return nil, true, err
	}

	if len(msg) == 0 ||

		bytes.HasSuffix(msg, []byte("\r\n")) {
		return nil, false, fmt.Errorf("protocol error: %s", string(msg))
	}

	state.bulkLen = 0

	return msg, false, nil
}

// parseMutiBulkHeader 解析多条数据头部
func parseMultiBulkHeader(msg []byte, state *readState) (err error) {
	expectedLine, err := strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return fmt.Errorf("protocol error: %s", string(msg))
	}

	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.expectedArgsCount = int(expectedLine)
		state.readingMultiLine = true
		state.args = make([][]byte, 0, state.expectedArgsCount)
	} else {
		return fmt.Errorf("protocol error: %s,%d", string(msg), expectedLine)
	}

	return nil
}

// +OK\r\n,-err\r\n,:100\r\n
func parserSingleLineReply(msg []byte) (recev resp.Reply, err error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	switch msg[0] {
	case '+':
		recev = reply.MakeStatusReply(str[1:])
	case '-':
		recev = reply.MakeErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("protocol error: %s", str)
		}
		recev = reply.MakeIntReply(int(val))
	default:
		err = fmt.Errorf("protocol error: %s", str)

	}
	return
}

// $3\r\nSET\r\n$3....
// PONG\r\n
func parserBody(msg []byte, state *readState) (err error) {
	line := msg[0 : len(msg)-2]

	// $3

	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return fmt.Errorf("protocol error: %s", string(msg))
		}

		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		} else {
			state.args = append(state.args, line)
		}

	}

	return nil
}
