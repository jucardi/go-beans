package beans

func triggerOnResolve(iInfo *instanceInfo) interface{} {
	if h, ok := iInfo.instance.(IFirstTimeResolveHandler); ok && !iInfo.resolvedFirstTime {
		h.OnFirstTimeResolve()
		iInfo.resolvedFirstTime = true
	}
	if h, ok := iInfo.instance.(IResolveHandler); ok {
		h.OnResolve()
	}
	return iInfo.instance
}
