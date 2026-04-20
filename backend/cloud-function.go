package p

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/package/cloud"
)

func ProcessImage(ctx context.Context, e cloud.GCSEvent) error {
	return cloud.ProcessImage(ctx, e)
}
