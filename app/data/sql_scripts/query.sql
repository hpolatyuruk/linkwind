select f.penalty, f.upvotes,f.downvotes,calculatestoryrank(f.penalty,f.submittedon,
						  f.upvotes,f.downvotes) from 
(select submittedon,upvotes,downvotes, 
calculatestorypenalty(stories.commentcount) as penalty from stories) f  
ORDER BY
	calculatestoryrank Desc;

