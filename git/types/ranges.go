package types

type Ref struct {
	id     string
	hash   string
	branch string
}

func (r *Ref) Id() string {
	return r.id
}

func (r *Ref) Hash() string {
	return r.hash
}

func (r *Ref) Branch() string {
	return r.branch
}

func NewRef(id, hash, branch string) *Ref {
	r := Ref{
		id:     id,
		hash:   hash,
		branch: branch,
	}
	return &r
}

// Range is used to represent start and endpoints of change-sets
type Range struct {
	from *Ref
	to   *Ref
}

// From returns the starting commit ref
func (r *Range) From() *Ref {
	return r.from
}

// To returns the ending commit ref
func (r *Range) To() *Ref {
	return r.to
}

// NewRange creates a new Range based on a starting and ending commit Ref
func NewRange(from *Ref, to *Ref) *Range {
	r := Range{
		from: from,
		to:   to,
	}
	return &r
}
