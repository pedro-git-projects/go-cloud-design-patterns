package stability

import "context"

// Circuit specifies the signature of the function that’s interacting
// with the database or other upstream service
// It has to have the signature of the function we want to limit
type Circuit func(context.Context) (string, error)
