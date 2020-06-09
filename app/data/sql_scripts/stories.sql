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