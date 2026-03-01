CREATE TABLE public.topics (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    content text NOT NULL,
    parent_id uuid,
    project_id uuid NOT NULL,
    sort_order int NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
	CONSTRAINT pk_topics PRIMARY KEY (id),
	CONSTRAINT fk_topics_parent FOREIGN KEY (parent_id) REFERENCES public.topics(id),
	CONSTRAINT fk_topics_project FOREIGN KEY (project_id) REFERENCES public.projects(id),
	CONSTRAINT unique_sort_order_per_parent UNIQUE (parent_id, sort_order) INITIALLY DEFERRED
);