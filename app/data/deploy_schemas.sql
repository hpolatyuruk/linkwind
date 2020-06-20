-- Table: public.customers

-- DROP TABLE public.customers;

CREATE TABLE public.customers
(
    id integer NOT NULL DEFAULT nextval('customers_id_seq'
    ::regclass),
    email character varying
    (50) COLLATE pg_catalog."default" NOT NULL,
    name character varying
    (25) COLLATE pg_catalog."default" NOT NULL,
    domain character varying
    (50) COLLATE pg_catalog."default",
    registeredon timestamp
    with time zone NOT NULL,
    imglogo bytea,
    title character varying
    (60) COLLATE pg_catalog."default",
    CONSTRAINT id_pkey PRIMARY KEY
    (id),
    CONSTRAINT uc_domain UNIQUE
    (domain)
        INCLUDE
    (domain),
    CONSTRAINT uc_email UNIQUE
    (email)
        INCLUDE
    (email),
    CONSTRAINT uc_name UNIQUE
    (name)
        INCLUDE
    (name)
)

TABLESPACE pg_default;

    ALTER TABLE public.customers
    OWNER to postgres;



    -- Table: public.users

    -- DROP TABLE public.users;

    CREATE TABLE public.users
    (
        fullname character varying(50) COLLATE pg_catalog
        ."default",
    email character varying
        (50) COLLATE pg_catalog."default" NOT NULL,
    password character varying
        (500) COLLATE pg_catalog."default" NOT NULL,
    website character varying
        (50) COLLATE pg_catalog."default",
    about character varying
        (100) COLLATE pg_catalog."default",
    invitecode character varying
        (20) COLLATE pg_catalog."default",
    karma double precision NOT NULL DEFAULT 0,
    username character varying
        (15) COLLATE pg_catalog."default" NOT NULL,
    id serial,
    registeredon timestamp
        with time zone NOT NULL,
    customerid integer,
    CONSTRAINT users_pkey PRIMARY KEY
        (id),
    CONSTRAINT unique_email UNIQUE
        (email)
,
    CONSTRAINT unique_username UNIQUE
        (username)
,
    CONSTRAINT customerid_fk FOREIGN KEY
        (customerid)
        REFERENCES public.customers
        (id) MATCH SIMPLE
        ON
        UPDATE CASCADE
        ON
        DELETE CASCADE
        NOT VALID
)

TABLESPACE
        pg_default;

        ALTER TABLE public.users
    OWNER to postgres;

        -- Index: ix_email

        -- DROP INDEX public.ix_email;

        CREATE INDEX ix_email
    ON public.users USING btree
        (email COLLATE pg_catalog."default")
    TABLESPACE pg_default;

        -- Index: ix_email_password

        -- DROP INDEX public.ix_email_password;

        CREATE INDEX ix_email_password
    ON public.users USING btree
        (email COLLATE pg_catalog."default", password COLLATE pg_catalog."default")
    TABLESPACE pg_default;

        -- Index: ix_username

        -- DROP INDEX public.ix_username;

        CREATE INDEX ix_username
    ON public.users USING btree
        (username COLLATE pg_catalog."default")
    TABLESPACE pg_default;

        -- Index: ix_username_password

        -- DROP INDEX public.ix_username_password;

        CREATE INDEX ix_username_password
    ON public.users USING btree
        (username COLLATE pg_catalog."default", password COLLATE pg_catalog."default")
    TABLESPACE pg_default;



        -- Table: public.stories

        -- DROP TABLE public.stories;

        CREATE TABLE public.stories
        (
            id serial NOT NULL,
            url character varying(500) COLLATE pg_catalog
            ."default",
    title character varying
            (250) COLLATE pg_catalog."default" NOT NULL,
    text text COLLATE pg_catalog."default",
    upvotes integer NOT NULL DEFAULT 0,
    commentcount integer NOT NULL DEFAULT 0,
    userid integer NOT NULL,
    submittedon timestamp
            with time zone NOT NULL,
    tags text[] COLLATE pg_catalog."default",
    downvotes integer NOT NULL DEFAULT 0,
    CONSTRAINT stories_pkey PRIMARY KEY
            (id),
    CONSTRAINT fk_userid FOREIGN KEY
            (userid)
        REFERENCES public.users
            (id) MATCH SIMPLE
        ON
            UPDATE NO ACTION
        ON
            DELETE CASCADE
        NOT VALID
)

TABLESPACE
            pg_default;

            ALTER TABLE public.stories
    OWNER to postgres;

            -- Index: ix_submittedon

            -- DROP INDEX public.ix_submittedon;

            CREATE INDEX ix_submittedon
    ON public.stories USING btree
            (submittedon DESC NULLS LAST)
    TABLESPACE pg_default;
            -- Index: ix_tags

            -- DROP INDEX public.ix_tags;

            CREATE INDEX ix_tags
    ON public.stories USING btree
            (tags COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;




            -- Table: public.storyvotes
            -- DROP TABLE public.storyvotes;

            CREATE TABLE public.storyvotes
            (
                storyid integer NOT NULL,
                userid integer NOT NULL,
                votetype integer NOT NULL,
                CONSTRAINT storyvotes_pk PRIMARY KEY (storyid, userid, votetype),
                CONSTRAINT storyid_fk FOREIGN KEY (storyid)
        REFERENCES public.stories (id)
                MATCH SIMPLE
        ON
                UPDATE NO ACTION
        ON
                DELETE CASCADE,
    CONSTRAINT userid_fk FOREIGN KEY
                (userid)
        REFERENCES public.users
                (id) MATCH SIMPLE
        ON
                UPDATE NO ACTION
        ON
                DELETE CASCADE
)

TABLESPACE pg_default;

                ALTER TABLE public.storyvotes
    OWNER to postgres;




                -- Table: public.saved

                -- DROP TABLE public.saved;

                CREATE TABLE public.saved
                (
                    storyid integer NOT NULL,
                    userid integer NOT NULL,
                    savedon timestamp
                    with time zone NOT NULL,
    CONSTRAINT saved_pkey PRIMARY KEY
                    (userid, storyid),
    CONSTRAINT storyid_fk FOREIGN KEY
                    (storyid)
        REFERENCES public.stories
                    (id) MATCH SIMPLE
        ON
                    UPDATE NO ACTION
        ON
                    DELETE CASCADE
        NOT VALID,
    CONSTRAINT userid_fk
                    FOREIGN KEY
                    (userid)
        REFERENCES public.users
                    (id) MATCH SIMPLE
        ON
                    UPDATE NO ACTION
        ON
                    DELETE CASCADE
)

TABLESPACE pg_default;

                    ALTER TABLE public.saved
    OWNER to postgres;

                    -- Index: ix_savedon

                    -- DROP INDEX public.ix_savedon;

                    CREATE INDEX ix_savedon
    ON public.saved USING btree
                    (savedon DESC)
    TABLESPACE pg_default;




                    -- Table: public.comments

                    -- DROP TABLE public.comments;

                    CREATE TABLE public.comments
                    (
                        comment text COLLATE pg_catalog
                        ."default" NOT NULL,
    upvotes integer NOT NULL,
    storyid integer NOT NULL,
    parentid integer,
    replycount integer NOT NULL DEFAULT 0,
    userid integer NOT NULL,
    commentedon timestamp
                        with time zone NOT NULL,
    id serial,
    downvotes integer NOT NULL,
    CONSTRAINT comments_pkey PRIMARY KEY
                        (id),
    CONSTRAINT "parentId_fk" FOREIGN KEY
                        (parentid)
        REFERENCES public.comments
                        (id) MATCH SIMPLE
        ON
                        UPDATE NO ACTION
        ON
                        DELETE CASCADE,
    CONSTRAINT storyid_fk FOREIGN KEY
                        (storyid)
        REFERENCES public.stories
                        (id) MATCH SIMPLE
        ON
                        UPDATE NO ACTION
        ON
                        DELETE CASCADE,
    CONSTRAINT userid_fk FOREIGN KEY
                        (userid)
        REFERENCES public.users
                        (id) MATCH SIMPLE
        ON
                        UPDATE NO ACTION
        ON
                        DELETE CASCADE
)

TABLESPACE pg_default;

                        ALTER TABLE public.comments
    OWNER to postgres;

                        -- Index: ix_commentedon

                        -- DROP INDEX public.ix_commentedon;

                        CREATE INDEX ix_commentedon
    ON public.comments USING btree
                        (commentedon DESC)
    TABLESPACE pg_default;




                        -- Table: public.commentvotes
                        -- DROP TABLE public.commentvotes;

                        CREATE TABLE public.commentvotes
                        (
                            commentid integer NOT NULL,
                            userid integer NOT NULL,
                            votetype integer NOT NULL,
                            CONSTRAINT commentvotes_pk PRIMARY KEY (commentid, userid, votetype),
                            CONSTRAINT commentid_fk FOREIGN KEY (commentid)
        REFERENCES public.comments (id)
                            MATCH SIMPLE
        ON
                            UPDATE NO ACTION
        ON
                            DELETE CASCADE
        NOT VALID,
    CONSTRAINT userid_fk
                            FOREIGN KEY
                            (userid)
        REFERENCES public.users
                            (id) MATCH SIMPLE
        ON
                            UPDATE NO ACTION
        ON
                            DELETE CASCADE
        NOT VALID
)

TABLESPACE
                            pg_default;

                            ALTER TABLE public.commentvotes
    OWNER to postgres;




                            -- Table: public.invitecodes

                            -- DROP TABLE public.invitecodes;

                            CREATE TABLE public.invitecodes
                            (
                                code character varying(20) COLLATE pg_catalog
                                ."default" NOT NULL,
    inviteruserid integer NOT NULL,
    createdon timestamp
                                with time zone NOT NULL,
    invitedemail character varying
                                (50) COLLATE pg_catalog."default" NOT NULL,
    used boolean NOT NULL DEFAULT false,
    CONSTRAINT invitecodes_pkey PRIMARY KEY
                                (code),
    CONSTRAINT userid_fk FOREIGN KEY
                                (inviteruserid)
        REFERENCES public.users
                                (id) MATCH SIMPLE
        ON
                                UPDATE NO ACTION
        ON
                                DELETE CASCADE
        NOT VALID
)

TABLESPACE
                                pg_default;

                                ALTER TABLE public.invitecodes
    OWNER to postgres;




                                -- Table: public.resetpasswordtokens

                                -- DROP TABLE public.resetpasswordtokens;

                                CREATE TABLE public.resetpasswordtokens
                                (
                                    token character varying(20) COLLATE pg_catalog
                                    ."default" NOT NULL,
    userid integer NOT NULL,
    CONSTRAINT resetpasswordtokens_pkey PRIMARY KEY
                                    (token),
    CONSTRAINT userid_fk FOREIGN KEY
                                    (userid)
        REFERENCES public.users
                                    (id) MATCH SIMPLE
        ON
                                    UPDATE NO ACTION
        ON
                                    DELETE NO ACTION
)

TABLESPACE
                                    pg_default;

                                    ALTER TABLE public.resetpasswordtokens
    OWNER to postgres;






                                    CREATE OR REPLACE FUNCTION public.calculatestorypenalty
                                    (
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
                                    end
                                    if;
    end loop;
                                    return 1;

                                    end ;

$BODY$;

                                    ALTER FUNCTION public.calculatestorypenalty(integer)
    OWNER TO postgres;






                                    CREATE OR REPLACE FUNCTION public.calculatestoryrank
                                    (
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
                                    end
                                    if;
    up :=  POWER
                                    ((votes),0.8);
    timeNow := NOW
                                    ();
    timeDiff :=  EXTRACT
                                    (EPOCH 
                         FROM
                                    (timeNow::timestamp - 
                               $1.submittedon::timestamp));
    down := POWER
                                    (timeDiff+1,0.1);
                                    return (up/down)*calculatestorypenalty($1.
                                    commentcount);
                                    end ;

$BODY$;

                                    ALTER FUNCTION public.calculatestoryrank(stories)
    OWNER TO postgres;





                                    -- Index: calculatestoryrank_idx

                                    -- DROP INDEX public.calculatestoryrank_idx;

                                    CREATE INDEX ix_calculatestoryrank
    ON public.stories USING btree
                                    (calculatestoryrank
                                    (stories.*) DESC NULLS FIRST)
    TABLESPACE pg_default;

