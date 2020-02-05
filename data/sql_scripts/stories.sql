-- Table: public.stories

-- DROP TABLE public.stories;

CREATE TABLE public.stories
(
    id serial,
    url character varying(500) COLLATE pg_catalog."default",
    title character varying(150) COLLATE pg_catalog."default" NOT NULL,
    text text COLLATE pg_catalog."default",
    upvotes integer NOT NULL DEFAULT 0,
    commentcount integer NOT NULL DEFAULT 0,
    userid integer NOT NULL,
    submittedon timestamp with time zone NOT NULL,
    tags text[] COLLATE pg_catalog."default",
    downvotes integer NOT NULL DEFAULT 0,
    CONSTRAINT stories_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.stories
    OWNER to postgres;

-- Index: ix_submittedon

-- DROP INDEX public.ix_submittedon;

CREATE INDEX ix_submittedon
    ON public.stories USING btree
    (submittedon DESC)
    TABLESPACE pg_default;

-- Index: ix_tags

-- DROP INDEX public.ix_tags;

CREATE INDEX ix_tags
    ON public.stories USING btree
    (tags COLLATE pg_catalog."default")
    TABLESPACE pg_default;