-- Table: public.saved

-- DROP TABLE public.saved;

CREATE TABLE public.saved
(
    storyid integer NOT NULL,
    userid integer NOT NULL,
    CONSTRAINT saved_pkey PRIMARY KEY (userid, storyid),
    CONSTRAINT storyid_fk FOREIGN KEY (storyid)
        REFERENCES public.stories (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT userid_fk FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE public.saved
    OWNER to postgres;