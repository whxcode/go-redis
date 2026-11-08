package deploy

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status uint8

const (
	StatusPending   Status = iota // 待执行
	StatusRunning                 // 执行中
	StatusExecuted                // 已执行
	StatusCancelled               // 已取消
	StatusFailed                  // 错误
)

func (s Status) Terminal() bool {
	return s == StatusExecuted || s == StatusCancelled || s == StatusFailed
}

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "待执行"
	case StatusRunning:
		return "执行中"
	case StatusExecuted:
		return "已执行"
	case StatusCancelled:
		return "已取消"
	case StatusFailed:
		return "错误"
	default:
		return "未知"
	}
}

type (
	TaskFunc    func() (string, error)
	TaskFuncCtx func(ctx context.Context) (string, error)
)

type TaskSpec struct {
	ID           string
	RunAt        time.Time
	Fn           TaskFunc
	FnCtx        TaskFuncCtx
	Timeout      time.Duration
	MaxAttempts  int
	RetryBackoff time.Duration
	Tags         map[string]string
}

var ErrInvalidSpec = errors.New("invalid spec")

// ErrInvalidSpec 表示任务定义非法（ID 为空、没有任务函数、超时/退避为负）。
func (spec *TaskSpec) validate() error {
	if strings.TrimSpace(spec.ID) == "" {
		return fmt.Errorf("%w: task id must not be empty", ErrInvalidSpec)
	}
	if spec.Fn == nil && spec.FnCtx == nil {
		return fmt.Errorf("%w: task %q has no function", ErrInvalidSpec, spec.ID)
	}
	if spec.Timeout < 0 {
		return fmt.Errorf("%w: task %q has negative timeout", ErrInvalidSpec, spec.ID)
	}
	if spec.RetryBackoff < 0 {
		return fmt.Errorf("%w: task %q has negative retry backoff", ErrInvalidSpec, spec.ID)
	}
	return nil
}
