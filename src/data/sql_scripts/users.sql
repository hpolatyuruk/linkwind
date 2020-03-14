-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE public.users
(
    fullname character varying(50) COLLATE pg_catalog."default",
    email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password character varying(500) COLLATE pg_catalog."default" NOT NULL,
    website character varying(50) COLLATE pg_catalog."default",
    about character varying(100) COLLATE pg_catalog."default",
    invitecode character varying(20) COLLATE pg_catalog."default",
    karma double precision NOT NULL DEFAULT 0,
    username character varying(15) COLLATE pg_catalog."default" NOT NULL,
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    registeredon timestamp with time zone NOT NULL,
    customerid integer,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT unique_email UNIQUE (email)
,
    CONSTRAINT unique_username UNIQUE (username)
,
    CONSTRAINT customerid_fk FOREIGN KEY (customerid)
        REFERENCES public.customers (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to postgres;

-- Index: ix_email

-- DROP INDEX public.ix_email;

CREATE INDEX ix_email
    ON public.users USING btree
    (email COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_email_password

-- DROP INDEX public.ix_email_password;

CREATE INDEX ix_email_password
    ON public.users USING btree
    (email COLLATE pg_catalog."default", password COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_username

-- DROP INDEX public.ix_username;

CREATE INDEX ix_username
    ON public.users USING btree
    (username COLLATE pg_catalog."default")
    TABLESPACE pg_default;

-- Index: ix_username_password

-- DROP INDEX public.ix_username_password;

CREATE INDEX ix_username_password
    ON public.users USING btree
    (username COLLATE pg_catalog."default", password COLLATE pg_catalog."default")
    TABLESPACE pg_default;