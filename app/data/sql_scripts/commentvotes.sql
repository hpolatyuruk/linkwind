-- Table: public.commentvotes
-- DROP TABLE public.commentvotes;

CREATE TABLE public.commentvotes
(
    commentid integer NOT NULL,
    userid integer NOT NULL,
    votetype integer NOT NULL,
    CONSTRAINT commentvotes_pk PRIMARY KEY (commentid, userid, votetype),
    CONSTRAINT commentid_fk FOREIGN KEY (commentid)
        REFERENCES public.comments (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT userid_fk FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE public.commentvotes
    OWNER to postgres;