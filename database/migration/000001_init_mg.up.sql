CREATE TABLE public.users
(
    id         uuid NOT NULL,
    name       varchar NULL,
    age        int NULL,
    created_at timestamp with time zone NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
);
