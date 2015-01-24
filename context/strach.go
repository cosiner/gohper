package context

type _serverStrach struct {
	tmplNames    []string
	funcHandlers map[string]*funcHandler
}

var _initStrach = &_serverStrach{
	funcHandlers: make(map[string]*funcHandler),
}

func strachAddTmpl(name string) {
	_initStrach.tmplNames = append(_initStrach.tmplNames, name)
}

func strachTmpls() []string {
	return _initStrach.tmplNames
}

func strachAddFuncHandler(pattern string, handler *funcHandler) {
	_initStrach.funcHandlers[pattern] = handler
}

func strachFuncHandler(pattern string) *funcHandler {
	return _initStrach.funcHandlers[pattern]
}

func strachDestroy() {
	_initStrach = nil
}
