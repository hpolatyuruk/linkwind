-- Table: public.comments

-- DROP TABLE public.comments;

CREATE TABLE public.comments
(
    comment text COLLATE pg_catalog."default" NOT NULL,
    upvotes integer NOT NULL,
    storyid integer NOT NULL,
    parentid integer,
    replycount integer NOT NULL DEFAULT 0,
    userid integer NOT NULL,
    commentedon timestamp with time zone NOT NULL,
    id serial NOT NULL,
    downvotes integer NOT NULL,
    CONSTRAINT comments_pkey PRIMARY KEY (id),
    CONSTRAINT "parentId_fk" FOREIGN KEY (parentid)
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

-- Index: ix_commentedon

-- DROP INDEX public.ix_commentedon;

CREATE INDEX ix_commentedon
    ON public.comments USING btree
    (commentedon DESC)
    TABLESPACE pg_default;