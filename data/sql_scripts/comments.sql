-- Table: public.comments

-- DROP TABLE public.comments;

CREATE TABLE public.comments
(
    id integer NOT NULL,
    comment text COLLATE pg_catalog."default" NOT NULL,
    commentedon time with time zone NOT NULL,
    upvotes integer NOT NULL,
    storyid integer NOT NULL,
    parentid integer,
    downvotes integer NOT NULL DEFAULT 0,
    replycount integer NOT NULL DEFAULT 0,
    userid integer NOT NULL,
    CONSTRAINT comments_pkey PRIMARY KEY (id),
    CONSTRAINT parentid_fk FOREIGN KEY (parentid)
        REFERENCES public.comments (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT storyid_fk FOREIGN KEY (storyid)
        REFERENCES public.stories (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT userid_fk FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE public.comments
    OWNER to postgres;