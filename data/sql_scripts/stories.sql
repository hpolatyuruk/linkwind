-- Table: public.stories

-- DROP TABLE public.stories;

CREATE TABLE public.stories
(
    id integer NOT NULL DEFAULT nextval('stories_id_seq'::regclass),
    url character varying(500) COLLATE pg_catalog."default",
    title character varying(150) COLLATE pg_catalog."default" NOT NULL,
    text text COLLATE pg_catalog."default",
    tags character varying[] COLLATE pg_catalog."default",
    upvotes integer NOT NULL DEFAULT 0,
    commentcount integer NOT NULL DEFAULT 0,
    downvotes integer NOT NULL DEFAULT 0,
    userid integer NOT NULL,
    submittedon timestamp with time zone NOT NULL,
    CONSTRAINT stories_pkey PRIMARY KEY (id),
    CONSTRAINT userid_fk FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE public.stories
    OWNER to postgres;