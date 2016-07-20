// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package stats

import (
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/text"
	"monkeybird/tr"
	"sync"
	"time"
)

// Stats defines general bot statistics.
type Stats struct {
	m        sync.RWMutex `json:"-"`
	Channels ChannelList  `json:",omitempty"`
}

// Update updates the channel and user statistics, from the given request.
func (s *Stats) Update(w irc.ResponseWriter, r *irc.Request) {
	if !r.FromChannel() {
		return
	}

	s.m.Lock()

	// Update the appropriate channel- and user data.
	cs := s.Channels.Get(w, r.Target)
	us := cs.Users.Get(r.SenderMask)
	us.AddNickname(r.SenderName)
	us.LastSeen = time.Now()
	s.m.Unlock()
}

// FirstOn finds out when a specific user was first seen in
// the channel from whence this command was issued.
func (s *Stats) FirstOn(w irc.ResponseWriter, r *cmd.Request) {
	user, us, ok := s.findUser(w, r)
	if !ok {
		return
	}

	proto.PrivMsg(w, r.Target,
		tr.FirstOnDisplayText,
		r.SenderName,

		text.Bold(user),
		us.FirstSeen.Format(tr.DateFormat),
		us.FirstSeen.Format(tr.TimeFormat),
		time.Since(us.FirstSeen),
	)
}

// LastOn finds out when a specific user was last seen in
// the channel from whence this command was issued.
func (s *Stats) LastOn(w irc.ResponseWriter, r *cmd.Request) {
	user, us, ok := s.findUser(w, r)
	if !ok {
		return
	}

	proto.PrivMsg(w, r.Target,
		tr.LastOnDisplayText,
		r.SenderName,

		text.Bold(user),
		us.LastSeen.Format(tr.DateFormat),
		us.LastSeen.Format(tr.TimeFormat),
		time.Since(us.LastSeen),
	)
}

// findUser finds stats for te user who sent the given request.
// Returns false if it could not be located.
func (s *Stats) findUser(w irc.ResponseWriter, r *cmd.Request) (string, UserStats, bool) {
	if !r.FromChannel() {
		proto.PrivMsg(w, r.SenderName, tr.StatsNotInChannel)
		return "", UserStats{}, false
	}

	user := r.SenderName
	if r.Len() > 0 {
		user = r.String(0)
	}

	s.m.RLock()
	defer s.m.RUnlock()

	cs := s.Channels.Get(w, r.Target)
	us := cs.Users.Find(user)

	if us == nil {
		proto.PrivMsg(w, r.Target, tr.StatsNoSuchUser,
			r.SenderName, text.Bold(user))
		return user, UserStats{}, false
	}

	return user, *us, true
}

// loadStats loads stats data from a file.
func (s *Stats) Load(file string) error {
	s.m.Lock()
	defer s.m.Unlock()
	return mod.Load(file, s, true)
}

// Save saves stats data to a file.
func (s *Stats) Save(file string) error {
	s.m.RLock()
	defer s.m.RUnlock()
	return mod.Save(file, s, true)
}
