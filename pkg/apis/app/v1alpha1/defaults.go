package v1alpha1

const (
	defaultBaseImage          = "loggregator/log-cache"
	defaultVersion            = "latest"
	defaultNozzleBaseImage    = "loggregator/log-cache-nozzle"
	defaultNozzleVersion      = "latest"
	defaultSchedulerBaseImage = "loggregator/log-cache-scheduler"
	defaultSchedulerVersion   = "latest"
	defaultGatewayBaseImage   = "loggregator/log-cache-gateway"
	defaultGatewayVersion     = "latest"
)

// SetDefaults Sets the default values for LogCache Operator and returns true if the spec was changed
func (l *LogCache) SetDefaults() bool {
	changed := false
	co := &l.Spec

	// Setup defaults for LogCache
	if co.LogCachePod.Nodes == 0 {
		co.LogCachePod.Nodes = 1
		changed = true
	}
	if len(co.LogCachePod.BaseImage) == 0 {
		co.LogCachePod.BaseImage = defaultBaseImage
		changed = true
	}
	if len(co.LogCachePod.Version) == 0 {
		co.LogCachePod.Version = defaultVersion
		changed = true
	}

	// Setup defaults for LogCache Gateway

	if len(co.LogCachePod.GatewayBaseImage) == 0 {
		co.LogCachePod.GatewayBaseImage = defaultGatewayBaseImage
		changed = true
	}
	if len(co.LogCachePod.GatewayVersion) == 0 {
		co.LogCachePod.GatewayVersion = defaultGatewayVersion
		changed = true
	}

	// Setup defaults for LogCacheNozzle
	if co.LogCacheNozzle.Nodes == 0 {
		co.LogCacheNozzle.Nodes = 1
		changed = true
	}

	if len(co.LogCacheNozzle.BaseImage) == 0 {
		co.LogCacheNozzle.BaseImage = defaultNozzleBaseImage
		changed = true
	}
	if len(co.LogCacheNozzle.Version) == 0 {
		co.LogCacheNozzle.Version = defaultNozzleVersion
		changed = true
	}

	// Setup defaults for LogCacheScheduler
	if co.LogCacheScheduler.Nodes == 0 {
		co.LogCacheScheduler.Nodes = 1
		changed = true
	}

	if len(co.LogCacheScheduler.BaseImage) == 0 {
		co.LogCacheScheduler.BaseImage = defaultSchedulerBaseImage
		changed = true
	}
	if len(co.LogCacheScheduler.Version) == 0 {
		co.LogCacheScheduler.Version = defaultSchedulerVersion
		changed = true
	}

	return changed
}
