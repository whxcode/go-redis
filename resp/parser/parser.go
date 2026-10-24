// parser 包用于解析 RESP 协议数据
// $3\r\nwhx\r\n
// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
package parser

import (
	"bufio"
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
	// fmt.Printf("readState.finished() called: expected=%d, actual=%d\n", s.expectedArgsCount, len(s.args))
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

func parse0(reader io.Reader, ch chan<- *Payload) {
	bufReader := bufio.NewReader(reader)
	state := &readState{}
	for {
		// input, _ := bufReader.Peek(100) // 预读一个字节，避免阻塞
		// log.Printf("parse0: input=%q\n", input)
		msg, ioErr, err := readLine(bufReader, state)

		if err == io.EOF {
			if ioErr {
				ch <- &Payload{
					Data: nil,
					Err:  err,
				}

				close(ch)
				return
			}

			ch <- &Payload{
				Data: nil,
				Err:  fmt.Errorf("protocol error: %s", string(msg)),
			}

			continue
		}

		if len(msg) == 0 {
			continue
		}

		if !state.readingMultiLine {
			if msg[0] == '*' {
				err := parseMultiBulkHeader(msg, state)
				//
				if err != nil {
					ch <- &Payload{
						Err: fmt.Errorf("protocol error: %s", string(msg)),
					}

					state = &readState{}
					continue
				}
				//
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: reply.MakeEmptyBulkReply(),
					}
					state = &readState{}
					continue
				}
				// $3\r\nxxx\r\n
			} else if msg[0] == '$' {
				err := parseBulkHeader(msg, state)
				if err != nil {
					ch <- &Payload{
						Err: fmt.Errorf("protocol error: %s", string(msg)),
					}

					state = &readState{}
					continue
				}

				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: reply.MakeEmptyBulkReply(),
					}
					state = &readState{}
					continue
				}
			} else { // ":"
				result, err := parserSingleLineReply(msg)
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = &readState{}
				continue
			}
		} else {
			err := parserBody(msg, state)
			if err != nil {
				ch <- &Payload{
					Err: fmt.Errorf("parserBody protocol error: %s", string(msg)),
				}
				continue
			}

			if state.finished() {
				var result resp.Reply
				switch state.msgType {
				case '*':
					result = reply.MakeMultiBulkReply(state.args)
					// case '$':
				default:
					// WARNING: 这儿不是应该返回一个 BulkReply 吗？为什么要返回 MultiBulkReply？
					///因为这儿强制断言了 home/whx/study/go-redis/resp/handler/handler.go:86:3
					// r, ok := payload.Data.(*reply.MultiBulkReply)
					// result = reply.MakeBulkReply(state.args[0])
					result = reply.MakeMultiBulkReply(state.args)
				}

				ch <- &Payload{
					Data: result,
				}

				state = &readState{}

			}

		}

	}
}

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

		return msg, false, nil

	}

	// 正确的做法：严格按照长度读取
	// $3\r\n 代表后面有 3 个字节 + 2 个字节的 \r\n
	msg = make([]byte, state.bulkLen+2)
	_, err = io.ReadFull(bufReader, msg)
	if err != nil {
		return nil, true, err
	}

	// 校验一下后缀对不对
	if msg[state.bulkLen] != '\r' || msg[state.bulkLen+1] != '\n' {
		return nil, false, fmt.Errorf("protocol error: missing CRLF")
	}

	state.bulkLen = 0 // 读完后必须归零，不然下次循环又跑这里来了

	return msg, false, nil
}

// parseMultiBulkHeader 解析多条数据头部
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

func parseBulkHeader(msg []byte, state *readState) (err error) {
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return fmt.Errorf("parseBulkHeader protocol error: %s", string(msg))
	}

	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, state.expectedArgsCount)
	} else {
		return fmt.Errorf("parseBulkHeader protocol error: %s", string(msg))
	}

	return
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
	// 去掉尾部的 \r\n
	line := msg[0 : len(msg)-2]

	if len(line) == 0 {
		return nil
	}

	if line[0] == '$' {
		// 说明这是一个新的 Bulk String 头部，例如 $3
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return fmt.Errorf("protocol error: %s", string(msg))
		}
		// 注意：这里不要 append，真正的值在下一行读进来
	} else {
		// 说明这是真正的参数数据，例如 SET
		state.args = append(state.args, line)
	}

	return nil
}
