package stability

import "context"

// Effector defines the signature of thje function
// that interacts with a service in the Retry pattern
type Effector func(context.Context) (string, error)
