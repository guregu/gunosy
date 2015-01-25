package daihinmin

import (
	"log"
	"sync"
	"time"
)

var matches = make(map[string]*match)
var matchesMutex = &sync.RWMutex{}

type match struct {
	name    string
	id      string
	size    int
	game    *Game
	users   map[sesh]*client
	players map[sesh]*Player
	readies map[sesh]bool

	join    chan joinReq
	part    chan sesh
	infoplz chan *client
	timeout <-chan time.Time
	die     chan struct{}
}

func NewMatch(name string) *match {
	m := &match{
		name:    name,
		id:      generateID("m:"),
		size:    4,
		game:    NewGame(),
		users:   make(map[sesh]*client),
		players: make(map[sesh]*Player),
		readies: make(map[sesh]bool),

		join:    make(chan joinReq),
		part:    make(chan sesh),
		infoplz: make(chan *client),
		timeout: make(<-chan time.Time),
		die:     make(chan struct{}),
	}
	m.register()
	return m
}

func (m *match) register() {
	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	if _, exists := matches[m.name]; exists {
		panic("Remaking match: " + m.name)
	}
	matches[m.id] = m

	go m.run()
}

func (m *match) unregister() {
	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	log.Printf("match dying: [%s] %s", m.id, m.name)
	if _, exists := matches[m.id]; !exists {
		panic("Deleting non-existent match: " + m.id)
	}
	delete(matches, m.id)
}

func (m match) String() string {
	return m.name + " (" + m.id + ")"
}

func (m *match) run() {
	log.Printf("Running match: %s", m.name)
	defer m.unregister()

	for {
		select {
		case req := <-m.join:
			// TODO reconnect voodoo
			m.users[req.sesh] = req.from
			m.broadcast(UserJoinPartReply{
				X:    "user-join",
				Chan: m.id,
				User: req.from.username(),
			})
			m.broadcast(m.info())
			if req.result != nil {
				req.result <- reqResult{ok: true}
			}
		case s := <-m.part:
			m.goodbye(s)
			// if everyone leaves, die
			if m.usercount() == 0 {
				return
			}
		case c := <-m.infoplz:
			c.send(m.info())
		case <-m.die:
			return
		}
	}
}

func (m *match) broadcast(msg interface{}) {
	for _, c := range m.users {
		c.send(msg)
	}
}

func (m *match) usercount() int {
	return len(m.users)
}

func (m *match) usernames() []string {
	var names []string
	for _, u := range m.users {
		names = append(names, u.username())
	}
	return names
}

func (m *match) find(name string) (s sesh, ok bool) {
	for _, u := range m.users {
		if u.username() == name {
			return u.session, true
		}
	}
	return
}

func (m *match) goodbye(s sesh) bool {
	c, ok := m.users[s]
	if !ok {
		return false
	}

	delete(m.users, s)
	c.match = nil // TODO fix data race?
	// TODO readies players
	m.broadcast(UserJoinPartReply{
		X:    "user-part",
		Chan: m.id,
		User: c.username(),
	})
	m.broadcast(m.info())
	return true
}

func (m *match) info() GameInfo {
	return GameInfo{
		X:     "game-info",
		ID:    m.id,
		Name:  m.name,
		Users: m.usernames(),
	}
}

func matchExists(id string) bool {
	matchesMutex.RLock()
	defer matchesMutex.RUnlock()

	_, exists := matches[id]
	return exists
}

func getMatch(id string) *match {
	matchesMutex.RLock()
	m, exists := matches[id]
	matchesMutex.RUnlock()

	if !exists {
		log.Printf("no match: %s", id)
	}
	return m
}

type joinReq struct {
	sesh
	from     *client
	password string
	result   chan reqResult
}

type matchReq struct {
	from     sesh
	name     string
	password string
	result   chan *match
}
