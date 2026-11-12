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
}

func NewNode(id int, peers []int, rng *rand.Rand) *Node {
	n := &Node{id: id, peers: peers, rng: rng, votedFor: -1, leader: -1}
	n.timeout = ElectionMin + rng.Intn(ElectionSpread)
	return n
}

func (n *Node) State() (Role, int, int) { return n.role, n.term, n.leader }

func (n *Node) broadcast(kind Kind) []Msg {
	out := []Msg{}
	for _, p := range n.peers {
		out = append(out, Msg{Kind: kind, From: n.id, To: p, Term: n.term})
	}
	return out
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
		if m.Term == n.term && (n.votedFor == -1 || n.votedFor == m.From) {
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
			reply.Success = true
		}
		return []Msg{reply}
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
