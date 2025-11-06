package utils

import "context"

const keyRequestID = "requestID"

// SetRequestID sets the request ID in the context.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}