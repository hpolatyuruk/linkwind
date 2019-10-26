-- Table: public.resetpasswordtokens

-- DROP TABLE public.resetpasswordtokens;

CREATE TABLE public.resetpasswordtokens
(
    token character varying(20) COLLATE pg_catalog."default" NOT NULL,
    userid integer NOT NULL,
    CONSTRAINT resetpasswordtokens_pkey PRIMARY KEY (token),
    CONSTRAINT userid_fk FOREIGN KEY (userid)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.resetpasswordtokens
    OWNER to postgres;