package econ

const rankBucketTicks = 100

// RankMinCommands keeps a nearly idle window from qualifying for efficiency.
// 100 commands per five-minute window allows deliberate waiting and batching.
const RankMinCommands = 100

type RankBucket struct {
	Index    int64              `json:"index"`
	Sold     map[string]float64 `json:"sold,omitempty"`
	Commands int64              `json:"commands,omitempty"`
}

type RankWindow struct {
	Ticks   int64        `json:"ticks"`
	Buckets []RankBucket `json:"buckets,omitempty"`
}

func (r *RankWindow) bucketIndex() int64 { return max(0, r.Ticks-1) / rankBucketTicks }

// Advance counts only running time. The oldest partial bucket is discarded
// early (by at most 99 ticks), so expired sales can never earn a rank.
func (r *RankWindow) Advance(ticks int64) {
	r.Ticks += max(0, ticks)
	oldest := r.bucketIndex() - Window/rankBucketTicks + 1
	for len(r.Buckets) > 0 && r.Buckets[0].Index < oldest {
		r.Buckets = r.Buckets[1:]
	}
}

func (r *RankWindow) bucket() *RankBucket {
	index := r.bucketIndex()
	if len(r.Buckets) == 0 || r.Buckets[len(r.Buckets)-1].Index != index {
		r.Buckets = append(r.Buckets, RankBucket{Index: index})
	}
	return &r.Buckets[len(r.Buckets)-1]
}

func (r *RankWindow) RecordSale(good string, amount float64) {
	b := r.bucket()
	if b.Sold == nil {
		b.Sold = map[string]float64{}
	}
	b.Sold[good] += amount
}

func (r *RankWindow) RecordCommand() { r.bucket().Commands++ }

func (r *RankWindow) totals() (map[string]float64, int64) {
	sold := map[string]float64{}
	var commands int64
	for _, b := range r.Buckets {
		commands += b.Commands
		for good, amount := range b.Sold {
			sold[good] += amount
		}
	}
	return sold, commands
}
