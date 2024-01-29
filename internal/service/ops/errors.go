package ops

import "fmt"

// S3ClientFailedError represents an error when trying to start a S3 client.
type S3ClientFailedError struct {
	Err error
}

// Allow S3ClientFailedError to satisfy error interface.
func (e *S3ClientFailedError) Error() string {
	return fmt.Sprintf("Failed to start a s3 client: %q", e.Err)
}
