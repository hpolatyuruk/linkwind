-- FUNCTION: public.calculatestorypenalty(integer)

-- DROP FUNCTION public.calculatestorypenalty(integer);

CREATE OR REPLACE FUNCTION public.calculatestorypenalty(
	commentcount integer)
    RETURNS integer
    LANGUAGE 'plpgsql'

    COST 100
    IMMUTABLE 
AS $BODY$declare
    penalty integer := 40;
    i integer := -1;
begin 
    loop
        exit when penalty = i;
				i := i+1;
        if (commentCount = i) then
        return penalty-i;
        end if;
    end loop;
    return 1;

end ;

$BODY$;

ALTER FUNCTION public.calculatestorypenalty(integer)
    OWNER TO postgres;



-- FUNCTION: public.calculatestoryrank(stories)

-- DROP FUNCTION public.calculatestoryrank(stories);

CREATE OR REPLACE FUNCTION public.calculatestoryrank(
	stories)
    RETURNS double precision
    LANGUAGE 'plpgsql'

    COST 100
    IMMUTABLE 
AS $BODY$declare

    timeNow timestamp;

    votes integer;

    up integer;

    down float;

    timeDiff integer;

begin 
    votes := $1.upvotes- $1.downvotes;
    if (votes <= 0) then
        votes := 1;
        end if;
    up :=  POWER((votes),0.8);
    timeNow := NOW();
    timeDiff :=  EXTRACT(EPOCH 
                         FROM (timeNow::timestamp - 
                               $1.submittedon::timestamp));
    down := POWER(timeDiff+1,0.1);
    return (up/down)*calculatestorypenalty($1.commentcount);
end ;

$BODY$;

ALTER FUNCTION public.calculatestoryrank(stories)
    OWNER TO postgres;
