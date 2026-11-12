package main

import "math/rand"

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

type Kind int

const (
	Vote Kind = iota
	VoteReply
	Append
	AppendReply
)

type Msg struct {
	Kind    Kind
	From    int
	To      int
	Term    int
	Granted bool
	Success bool
	// log replication
	LastIndex int
	LastTerm  int
	PrevIndex int
	PrevTerm  int
	Entries   []Entry
	Commit    int
	Match     int
}

type Entry struct {
	Term int
	Cmd  string
}

const (
	ElectionMin    = 10
	ElectionSpread = 10
	HeartbeatEvery = 3
)

type Node struct {
	id       int
	peers    []int
	rng      *rand.Rand
	role     Role
	term     int
	votedFor int
	leader   int
	votes    map[int]bool
	elapsed  int
	timeout  int
	log      []Entry // index i+1 lives at log[i]
	commit   int
	next     map[int]int
	match    map[int]int
}

func NewNode(id int, peers []int, rng *rand.Rand) *Node {
	n := &Node{id: id, peers: peers, rng: rng, votedFor: -1, leader: -1}
	n.timeout = ElectionMin + rng.Intn(ElectionSpread)
	return n
}

func (n *Node) State() (Role, int, int) { return n.role, n.term, n.leader }

func (n *Node) lastIndexTerm() (int, int) {
	if len(n.log) == 0 {
		return 0, 0
	}
	return len(n.log), n.log[len(n.log)-1].Term
}

func (n *Node) termAt(index int) int {
	if index == 0 {
		return 0
	}
	return n.log[index-1].Term
}

func (n *Node) broadcast(kind Kind) []Msg {
	out := []Msg{}
	for _, p := range n.peers {
		m := Msg{Kind: kind, From: n.id, To: p, Term: n.term}
		if kind == Vote {
			m.LastIndex, m.LastTerm = n.lastIndexTerm()
		} else {
			m.PrevIndex = n.next[p] - 1
			m.PrevTerm = n.termAt(m.PrevIndex)
			m.Entries = append([]Entry{}, n.log[m.PrevIndex:]...)
			m.Commit = n.commit
		}
		out = append(out, m)
	}
	return out
}

func (n *Node) Propose(cmd string) (int, bool) {
	if n.role != Leader {
		return 0, false
	}
	n.log = append(n.log, Entry{Term: n.term, Cmd: cmd})
	n.advanceCommit()
	return len(n.log), true
}

func (n *Node) Committed() []string {
	out := []string{}
	for _, e := range n.log[:n.commit] {
		out = append(out, e.Cmd)
	}
	return out
}

// advanceCommit commits the highest current-term index stored on a majority.
func (n *Node) advanceCommit() {
	for index := len(n.log); index > n.commit; index-- {
		if n.log[index-1].Term != n.term {
			break
		}
		count := 1
		for _, p := range n.peers {
			if n.match[p] >= index {
				count++
			}
		}
		if count > (len(n.peers)+1)/2 {
			n.commit = index
			return
		}
	}
}

func (n *Node) Tick() []Msg {
	n.elapsed++
	if n.role == Leader {
		if n.elapsed >= HeartbeatEvery {
			n.elapsed = 0
			return n.broadcast(Append)
		}
		return []Msg{}
	}
	if n.elapsed < n.timeout {
		return []Msg{}
	}
	n.role = Candidate
	n.term++
	n.votedFor = n.id
	n.votes = map[int]bool{n.id: true}
	n.leader = -1
	n.elapsed = 0
	n.timeout = ElectionMin + n.rng.Intn(ElectionSpread)
	if n.won() {
		return n.becomeLeader()
	}
	return n.broadcast(Vote)
}

func (n *Node) won() bool { return len(n.votes) > (len(n.peers)+1)/2 }

func (n *Node) becomeLeader() []Msg {
	n.role = Leader
	n.leader = n.id
	n.elapsed = 0
	n.next = map[int]int{}
	n.match = map[int]int{}
	for _, p := range n.peers {
		n.next[p] = len(n.log) + 1
	}
	return n.broadcast(Append)
}

func (n *Node) Step(m Msg) []Msg {
	if m.Term > n.term {
		n.term = m.Term
		n.role = Follower
		n.votedFor = -1
		n.leader = -1
	}
	reply := Msg{From: n.id, To: m.From, Term: n.term}
	switch m.Kind {
	case Vote:
		reply.Kind = VoteReply
		lastIndex, lastTerm := n.lastIndexTerm()
		upToDate := m.LastTerm > lastTerm || (m.LastTerm == lastTerm && m.LastIndex >= lastIndex)
		if m.Term == n.term && upToDate && (n.votedFor == -1 || n.votedFor == m.From) {
			n.votedFor = m.From
			n.elapsed = 0
			reply.Granted = true
		}
		return []Msg{reply}
	case Append:
		reply.Kind = AppendReply
		if m.Term == n.term {
			n.role = Follower
			n.leader = m.From
			n.elapsed = 0
			if m.PrevIndex <= len(n.log) && n.termAt(m.PrevIndex) == m.PrevTerm {
				for i, e := range m.Entries {
					index := m.PrevIndex + 1 + i
					if index <= len(n.log) && n.log[index-1].Term != e.Term {
						n.log = n.log[:index-1]
					}
					if index > len(n.log) {
						n.log = append(n.log, e)
					}
				}
				reply.Match = m.PrevIndex + len(m.Entries)
				if c := min(m.Commit, reply.Match); c > n.commit {
					n.commit = c
				}
				reply.Success = true
			}
		}
		return []Msg{reply}
	case AppendReply:
		if n.role == Leader && m.Term == n.term {
			if m.Success {
				n.match[m.From] = max(n.match[m.From], m.Match)
				n.next[m.From] = n.match[m.From] + 1
				n.advanceCommit()
			} else if n.next[m.From] > 1 {
				n.next[m.From]--
			}
		}
	case VoteReply:
		if n.role == Candidate && m.Term == n.term && m.Granted {
			n.votes[m.From] = true
			if n.won() {
				return n.becomeLeader()
			}
		}
	}
	return []Msg{}
}

func main() {}
