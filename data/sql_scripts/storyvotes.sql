-- Table: public.storyvotes
-- DROP TABLE public.storyvotes;

CREATE TABLE public.storyvotes
(
    storyid integer NOT NULL,
    userid integer NOT NULL,
    votetype integer NOT NULL,
    CONSTRAINT storyvotes_pk PRIMARY KEY (storyid, userid, votetype),
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

ALTER TABLE public.storyvotes
    OWNER to postgres;