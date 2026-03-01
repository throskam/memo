CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
	username text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
	CONSTRAINT pk_users PRIMARY KEY (id)
);