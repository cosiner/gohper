package server

import (
	"strconv"
	"sync"
	"time"
)

type (
	// SessionStore is common interface of session store
	SessionStore interface {
		// Init init store with given config
		Init(conf string) error
		// IsExist check whether session with given id is exist
		IsExist(id string) bool
		// Save save values with given id and expire time
		Save(id string, values Values, expire uint64)
		// Get return values of given id
		Get(id string) Values
		// Rename move values exist in old id to new id
		Rename(oldId, newId string)
	}
	// Session represent a server session
	Session struct {
		id string
		*AttrContainer
	}
	// sessionNode is a session node record request count reference this session
	sessionNode struct {
		refs int
		sess *Session
	}
	// sessionManager is a session manager for server that is responsible for
	// get session from session store, generate session id, create new session
	// and store exist session to session store
	sessionManager struct {
		store    SessionStore
		expire   uint64
		sessions map[string]*sessionNode
		lock     *sync.Mutex
	}
)

// newSession create a new session with given id
func newSession(id string) *Session {
	return &Session{
		id:            id,
		AttrContainer: NewAttrContainer(),
	}
}

// newSessionWith create a new session with given id an initial attributes
func newSessionWith(id string, values Values) *Session {
	return &Session{
		id:            id,
		AttrContainer: NewAttrContainerVals(values),
	}
}

// sessionId return session's id
func (sess *Session) sessionId() string {
	return sess.id
}

// newSessionManager create a new session manager with given store and expire
func newSessionManager(store SessionStore, expire uint64) *sessionManager {
	return &sessionManager{
		store:    store,
		expire:   expire,
		sessions: make(map[string]*sessionNode),
		lock:     new(sync.Mutex),
	}
}

// newSessionId generate new session id different from exists
func (sm *sessionManager) newSessionId() string {
	sm.lock.Lock()
	for {
		id := strconv.Itoa(int(time.Now().UnixNano() / 10))
		if _, has := sm.sessions[id]; !has && !sm.store.IsExist(id) {
			sm.lock.Unlock()
			return id
		}
	}
}

// session return exist session
func (sm *sessionManager) session(id string) (sess *Session) {
	sm.lock.Lock()
	sn := sm.sessions[id]
	if sn != nil { // exist in session manager
		sn.refs++
		sess = sn.sess
		sm.lock.Unlock()
	} else {
		sm.lock.Unlock()
		if values := sm.store.Get(id); values != nil { // get from session store
			sess = sm.addSession(newSessionWith(id, values))
		}
	}
	return
}

// addSession add a session to manager
func (sm *sessionManager) addSession(sess *Session) *Session {
	sn := &sessionNode{refs: 1, sess: sess}
	sm.lock.Lock()
	sm.sessions[sess.id] = sn
	sm.lock.Unlock()
	return sess
}

// newSession create a new session from strach
func (sm *sessionManager) newSession() *Session {
	return sm.addSession(newSession(sm.newSessionId()))
}

// storeSession store session to session store, given session must exist in manager
func (sm *sessionManager) storeSession(sess *Session) {
	id, lock := sess.id, sm.lock
	lock.Lock()
	sn := sm.sessions[id]
	if sn == nil {
		panic("Unexpected: session haven't stored in server.sessions")
	}
	sn.refs--
	if sn.refs == 0 {
		delete(sm.sessions, id)
		lock.Unlock()
		sn.sess.AccessAllAttrs(func(values Values) {
			sm.store.Save(id, values, sm.expire)
		})
	} else {
		lock.Unlock()
	}
}
