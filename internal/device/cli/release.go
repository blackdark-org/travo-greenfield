package cli

import "context"

func (d *Device) VerifyRelease(ctx context.Context, manifest, signature, key string) error {
	_, err := d.run(ctx, TierRead, "usign", "-V", "-q", "-m", manifest, "-x", signature, "-p", key)
	return err
}
