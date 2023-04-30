package ads

import "time"

type Ad struct {
	ID        int64
	Title     string
	Text      string
	AuthorID  int64
	Published bool
	Created   time.Time
	Modified  time.Time
}
