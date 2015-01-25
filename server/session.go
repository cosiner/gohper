package server

import (
	"strconv"
	"sync"
	"time"
)

//==============================================================================
//                           Session
//==============================================================================
const _COOKIE_SESSION = "session"

type SessionStore interface {
	Init(conf string) error
	IsExist(id string) bool
	Save(id string, values map[string]interface{}, expire uint64)
	Get(id string) map[string]interface{}
	Rename(oldId, newId string)
}

type Session struct {
	id string
	*AttrContainer
}

func newSession(id string) *Session {
	return &Session{
		id:            id,
		AttrContainer: NewAttrContainer(),
	}
}

func newSessionWith(id string, values map[string]interface{}) *Session {
	return &Session{
		id:            id,
		AttrContainer: NewAttrContainerVals(values),
	}
}

func (sess *Session) sessionId() string {
	return sess.id
}

//==============================================================================
//                           Server Session
//==============================================================================
type sessionNode struct {
	refs int
	sess *Session
}

type serverSession struct {
	store    SessionStore
	expire   uint64
	sessions map[string]*sessionNode
	lock     *sync.Mutex
}

func newServerSession(store SessionStore, expire uint64) *serverSession {
	return &serverSession{
		store:    store,
		expire:   expire,
		sessions: make(map[string]*sessionNode),
		lock:     new(sync.Mutex),
	}
}

func (srvSess *serverSession) initStore(conf string) error {
	return srvSess.store.Init(conf)
}

func (srvSess *serverSession) newSessionId() string {
	srvSess.lock.Lock()
	for {
		id := strconv.Itoa(int(time.Now().UnixNano() / 10))
		if _, has := srvSess.sessions[id]; !has && !srvSess.store.IsExist(id) {
			srvSess.lock.Unlock()
			return id
		}
	}
}

func (srvSess *serverSession) session(id string) (sess *Session) {
	srvSess.lock.Lock()
	sn := srvSess.sessions[id]
	if sn != nil {
		sn.refs++
		sess = sn.sess
		srvSess.lock.Unlock()
	} else {
		srvSess.lock.Unlock()
		if values := srvSess.store.Get(id); values != nil {
			sess = srvSess.addSession(newSessionWith(id, values))
		}
	}
	return
}

func (srvSess *serverSession) addSession(sess *Session) *Session {
	sn := &sessionNode{refs: 1, sess: sess}
	srvSess.lock.Lock()
	srvSess.sessions[sess.id] = sn
	srvSess.lock.Unlock()
	return sess
}

func (srvSess *serverSession) newSession() *Session {
	return srvSess.addSession(newSession(srvSess.newSessionId()))
}

func (srvSess *serverSession) storeSession(sess *Session) {
	id, lock := sess.id, srvSess.lock
	lock.Lock()
	sn := srvSess.sessions[id]
	if sn == nil {
		panic("Unexpected: session haven't stored in server.sessions")
	}
	sn.refs--
	if sn.refs == 0 {
		delete(srvSess.sessions, id)
		lock.Unlock()
		sn.sess.AccessAllAttrs(func(values map[string]interface{}) {
			srvSess.store.Save(id, values, srvSess.expire)
		})
	} else {
		lock.Unlock()
	}
}
