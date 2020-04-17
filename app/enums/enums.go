package enums

/*VoteType represents the type of votes (upvote, downvote) for stories and comments.*/
type VoteType int

const (
	/*UpVote represents the positive vote for stories and comments.*/
	UpVote VoteType = 1
	/*DownVote represents the negative vote for stories and comments.*/
	DownVote VoteType = 2
)
