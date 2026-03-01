CREATE TYPE authentication_provider AS ENUM (
  'passwordless'
);

CREATE TABLE public.authentication_methods (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    provider authentication_provider NOT NULL,
    sub text NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
	CONSTRAINT pk_authentication_methods PRIMARY KEY (id),
	CONSTRAINT fk_authentication_methods_users FOREIGN KEY (user_id) REFERENCES public.users(id)
);