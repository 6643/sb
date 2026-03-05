// sbgen:fingerprint a19af4c1c7bea5279820fd3ac20ee22afd9cd252fe9d06233b0600a98d9a5950
package sb

import (
	"context"
)

func echo(ctx context.Context, env *Envelope) (result *Envelope, errCode RpcErrCode) {
	return nil, RpcRespErr
}
