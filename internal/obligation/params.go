package obligation

// CopyParams 返回义务参数的可写副本。
// 始终返回非 nil map：调用方（如 evalRule 盖 rule 戳）会直接写入返回值，
// nil 入参（JSON 里 parameters 为 null）不得导致后续写 panic。
func CopyParams(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
