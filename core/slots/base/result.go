/**
 * @description:
 *
 * @author: helloworld
 * @date:2019-07-11
 */
package base

type SlotResultStatus int8

const (
	ResultStatusPass = iota
	ResultStatusBlocked
	ResultStatusWait
	ResultStatusError
)

type TokenResult struct {
	Status        SlotResultStatus
	BlockedReason string
	WaitMs        uint64
	ErrorMsg      string
}

func NewResultPass() *TokenResult {
	return &TokenResult{Status: ResultStatusPass}
}

func NewResultBlock(blockedReason string) *TokenResult {
	return &TokenResult{Status: ResultStatusBlocked, BlockedReason: blockedReason}
}

func NewResultWait(waitMs uint64) *TokenResult {
	return &TokenResult{Status: ResultStatusWait, WaitMs: waitMs}
}

func NewResultError(errorMsg string) *TokenResult {
	return &TokenResult{Status: ResultStatusError, ErrorMsg: errorMsg}
}
