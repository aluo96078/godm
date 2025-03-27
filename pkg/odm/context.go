package odm

import "context"

// WithContext 設置 ODM 操作的自定義上下文。
// WithContext sets a custom context for the GODM operations.
func (o *GODM) WithContext(ctx context.Context) *GODM {
	o.Ctx = ctx
	return o
}

// getContext 返回設置的上下文，如果未設置則默認返回 context.TODO()。
// getContext returns the set context or defaults to context.TODO().
func (o *GODM) getContext() context.Context {
	if o.Ctx != nil {
		return o.Ctx
	}
	return context.TODO()
}
