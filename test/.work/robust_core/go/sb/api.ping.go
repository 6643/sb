// sbgen:fingerprint cf6ae706dfdf3c3f0a69469b999b809483ce86a5ad03d3a05424ebc447853448
package sb

import (
	"context"
)

func ping(ctx context.Context) (errCode RpcErrCode) {
	return RpcRespErr
}
