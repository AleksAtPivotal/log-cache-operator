package logcache

// labelsForLogCache returns the labels for selecting the resources
// belonging to the given Log Cache name.
func labelsForLogCache(name string) map[string]string {
	return map[string]string{"app": "logcache", "logcache_cluster": name}
}
