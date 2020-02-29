-- Table: public.customers

-- DROP TABLE public.customers;

CREATE TABLE public.customers
(
    id serial,
    email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    domain character varying(50) COLLATE pg_catalog."default" NOT NULL,
    registeredon timestamp with time zone NOT NULL,
    imglogo bytea,
    CONSTRAINT id_pkey PRIMARY KEY (id),
    CONSTRAINT uc_domain UNIQUE (domain)
        INCLUDE(domain),
    CONSTRAINT uc_email UNIQUE (email)
        INCLUDE(email),
    CONSTRAINT uc_name UNIQUE (name)
        INCLUDE(name)
)

TABLESPACE pg_default;

ALTER TABLE public.customers
    OWNER to postgres;