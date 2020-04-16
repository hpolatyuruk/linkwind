-- FUNCTION: public.calculatestorypenalty(integer)

-- DROP FUNCTION public.calculatestorypenalty(integer);

CREATE OR REPLACE FUNCTION public.calculatestorypenalty(
	commentcount integer)
    RETURNS integer
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE 
    
AS $BODY$declare
	penalty integer := 40;
	i integer := 0;
begin 
	loop
		exit when penalty = i;
		i := i+1;
		if (commentCount = i) then
		return penalty-i;
		end if;
	end loop;
	
	return 0;
end ;
$BODY$;

ALTER FUNCTION public.calculatestorypenalty(integer)
    OWNER TO postgres;



-- FUNCTION: public.calculatestoryrank(integer, timestamp with time zone, integer, integer)

-- DROP FUNCTION public.calculatestoryrank(integer, timestamp with time zone, integer, integer);

CREATE OR REPLACE FUNCTION public.calculatestoryrank(
	penalty integer,
	submittedon timestamp with time zone,
	upvotes integer,
	downvotes integer)
    RETURNS double precision
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE 
    
AS $BODY$declare
	score float;
	timeNow timestamp;
	votes integer;
	up integer;
	down float;
	timeDiff integer;
begin 
	votes := upvotes- downvotes;
	if (votes <= 0) then
		votes := 1;
		end if;
	up :=  POWER((votes-1),0.8);
	
	timeNow := NOW();
	timeDiff :=  EXTRACT(EPOCH 
						 FROM (timeNow::timestamp - 
							   submittedon::timestamp));
	down := POWER(timeDiff+2,1.8);
	score := (up/down)*penalty;
	return score;
end ;
$BODY$;

ALTER FUNCTION public.calculatestoryrank(integer, timestamp with time zone, integer, integer)
    OWNER TO postgres;
