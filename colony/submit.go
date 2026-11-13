package colony

// Submit checks an answer against the active archive contract's tracked truth.
// The answer must be JSON encodable. A wrong answer fails the attempt.
func Submit(key string, v any) bool { return world.Submit(key, v) }
