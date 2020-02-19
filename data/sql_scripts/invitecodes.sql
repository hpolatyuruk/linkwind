-- Table: public.invitecodes

-- DROP TABLE public.invitecodes;

CREATE TABLE public.invitecodes
(
    code character varying(20) COLLATE pg_catalog."default" NOT NULL,
    inviteruserid integer NOT NULL,
    createdon timestamp with time zone NOT NULL,
    invitedemail character varying(50) COLLATE pg_catalog."default" NOT NULL,
    used boolean NOT NULL DEFAULT false,
    CONSTRAINT invitecodes_pkey PRIMARY KEY (code),
    CONSTRAINT userid_fk FOREIGN KEY (inviteruserid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE public.invitecodes
    OWNER to postgres;