package to_auth_service

import "time"

const (
	maxRetries             = 5
	minBackoff             = 200 * time.Millisecond
	jitter                 = 0.2 // 20% jitter
	timeoutWaitChangeState = 5 * time.Second
)
