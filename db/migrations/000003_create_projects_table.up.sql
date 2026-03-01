CREATE TABLE public.projects (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    owner_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
	CONSTRAINT pk_projects PRIMARY KEY (id),
	CONSTRAINT fk_projects_owner FOREIGN KEY (owner_id) REFERENCES public.users(id)
);