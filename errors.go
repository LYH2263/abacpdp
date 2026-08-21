package abacpdp

import "errors"

var (
	// ErrClosed 表示 PDP 已关闭。
	ErrClosed = errors.New("abacpdp: closed")
	// ErrNoPolicy 表示尚未装载可用策略集。
	ErrNoPolicy = errors.New("abacpdp: no policy set")
	// ErrUnknownCombine 表示合并算法名称未知。
	ErrUnknownCombine = errors.New("abacpdp: unknown combine algorithm")
	// ErrInvalidPolicy 表示策略校验失败。
	ErrInvalidPolicy = errors.New("abacpdp: invalid policy")
	// ErrPersist 表示策略快照持久化失败。
	ErrPersist = errors.New("abacpdp: persist failed")
	// ErrCanceled 表示上下文已取消。
	ErrCanceled = errors.New("abacpdp: canceled")
	// ErrRemoteAttr 表示远程属性解析失败。
	ErrRemoteAttr = errors.New("abacpdp: remote attribute")
)
