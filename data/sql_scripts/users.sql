-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE public.users
(
    fullname character varying(50) COLLATE pg_catalog."default",
    email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password character varying(50) COLLATE pg_catalog."default" NOT NULL,
    website character varying(50) COLLATE pg_catalog."default",
    about character varying(500) COLLATE pg_catalog."default",
    invitedby character varying(15) COLLATE pg_catalog."default",
    invitecode character varying(15) COLLATE pg_catalog."default",
    karma double precision NOT NULL DEFAULT 0,
    username character varying(15) COLLATE pg_catalog."default" NOT NULL,
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    registeredon timestamp with time zone NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT unique_email UNIQUE (email)

)

TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to postgres;