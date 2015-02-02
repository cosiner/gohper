package server

import (
	"strconv"
	"sync"
	"time"

	. "github.com/cosiner/golib/errors"
)

const (
	_SESSION_DISABLE = "Session has been disabled!"
)

type (
	// Session represent a server session
	Session struct {
		id string
		AttrContainer
	}

	// SessionStore is common interface of session store
	// it's responsible for store session values by id, also need to manage
	// values' lifetime
	SessionStore interface {
		// Init init store with given config
		Init(conf string) error
		// Destroy destroy the session store, release all resources it owned
		Destroy()
		// IsExist check whether session with given id is exist
		IsExist(id string) bool
		// Save save values with given id and lifetime time
		// lifetime :<0 means never expired, 0 means delete right now, >0 means
		// set up lifetime
		Save(id string, values Values, lifetime int64)
		// Get return values of given id
		Get(id string) Values
		// Rename move values exist in old id to new id
		Rename(oldId, newId string)
	}

	// SessionManager is a manager of session responsible for get session and
	// create new session, generate new session id, and store session,
	SessionManager interface {
		// SetStore set up session store and session lifetime for manager
		Init(store SessionStore, lifetime int64) error
		// Destroy destroy session manager, also responsible for destroy session store
		Destroy()
		// Session acquire a exist session by id, if not exist, nil is retuened
		Session(id string) *Session
		// NewSession create a new session with new id
		NewSession() *Session
		// StoreSession store an session
		StoreSession(*Session)
	}

	// manageNode is a session node record request count reference this session
	// for session manager
	manageNode struct {
		refs int
		sess *Session
	}
	// sessionManager is a session manager for server that is responsible for
	// get session from session store, generate session id, create new session
	// and store exist session to session store
	sessionManager struct {
		store    SessionStore
		lifetime int64
		sessions map[string]*manageNode
		lock     *sync.Mutex
	}

	// panicSessionManager is a empty session manager, if use it, any operation
	// to get session from it will cause panic
	panicSessionManager struct{}
)

// newSession create a new session with given id
func NewSession(id string) *Session {
	return &Session{
		id:            id,
		AttrContainer: NewLockedAttrContainer(),
	}
}

// NewSessionWith create a new session with given id an initial attributes
func NewSessionWith(id string, values Values) *Session {
	return &Session{
		id:            id,
		AttrContainer: NewLockedAttrContainerWith(values),
	}
}

// Id return session's id
func (sess *Session) Id() string {
	return sess.id
}

// panic session manager

func newPanicSessionManager() SessionManager               { return panicSessionManager{} }
func (panicSessionManager) Init(SessionStore, int64) error { return nil }
func (panicSessionManager) Session(string) *Session        { PanicServer(_SESSION_DISABLE); return nil }
func (panicSessionManager) NewSession() *Session           { PanicServer(_SESSION_DISABLE); return nil }
func (panicSessionManager) StoreSession(*Session)          { PanicServer(_SESSION_DISABLE) }
func (panicSessionManager) Destroy()                       {}

// NewSessionManager create a new session manager
func NewSessionManager() SessionManager {
	return &sessionManager{
		sessions: make(map[string]*manageNode),
		lock:     new(sync.Mutex),
	}
}

// Init set up store and lifetime for session manager
func (sm *sessionManager) Init(store SessionStore, lifetime int64) error {
	if store == nil {
		return Err("Empty session store")
	}
	if lifetime == 0 {
		return Err("Session lifetime is zero 0 no session will be stored")
	}
	sm.store = store
	sm.lifetime = lifetime
	return nil
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
func (sm *sessionManager) Session(id string) (sess *Session) {
	sm.lock.Lock()
	sn := sm.sessions[id]
	if sn != nil { // exist in session manager
		sn.refs++
		sess = sn.sess
		sm.lock.Unlock()
	} else {
		sm.lock.Unlock()
		if values := sm.store.Get(id); values != nil { // get from session store
			sess = sm.addSession(NewSessionWith(id, values))
		}
	}
	return
}

// addSession add a session to manager
func (sm *sessionManager) addSession(sess *Session) *Session {
	sn := &manageNode{refs: 1, sess: sess}
	sm.lock.Lock()
	sm.sessions[sess.Id()] = sn
	sm.lock.Unlock()
	return sess
}

// NewSession create a new session from strach
func (sm *sessionManager) NewSession() *Session {
	return sm.addSession(NewSession(sm.newSessionId()))
}

// StoreSession store session to session store, given session must exist in manager
func (sm *sessionManager) StoreSession(sess *Session) {
	id, lock := sess.Id(), sm.lock
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
			sm.store.Save(id, values, sm.lifetime)
		})
	} else {
		lock.Unlock()
	}
}

// Destroy destroy destroy session manager and store
func (sm *sessionManager) Destroy() {
	sm.lock.Lock()
	sm.sessions = nil
	sm.lock.Unlock()
	sm.lock = nil
	sm.store.Destroy()
}
