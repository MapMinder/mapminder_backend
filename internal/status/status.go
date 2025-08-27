package status

import "net/http"

type Status struct {
	Code    int
	Message string
}

/*

全てのStatusを最初の段階で作成するのは不可能のため必要な作成する方針にする

INDEX:
N で始まる Status はSuccess 系のステータスを表す
W で始まる Status はWarn(警告) 系のステータスを表す
E で始まる Status はError(エラー) 系のステータスを表す

*/

// TODO: 現時点で書かれている Status は決まっているものではなくて、これらを別のMessageなどにした方が良い場合は修正が必要

var (
	N0000 = Status{Code: http.StatusOK, Message: "SUCCESS"}
	W2000 = Status{Code: http.StatusBadRequest, Message: "W2000: Values you provided are not valid please fix the values and try again."}
	E5000 = Status{Code: http.StatusInternalServerError, Message: "E5000: Internal server occurred contact the service provider."}
)
