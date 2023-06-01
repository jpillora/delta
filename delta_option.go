package delta

type Option func(*options)

type options struct {
	apply bool
}

// OptionApply sets whether the patch should be applied to the target
// object. By default, the patch is applied.
func OptionApply(apply bool) Option {
	return func(o *options) {
		o.apply = apply
	}
}
