-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create auth schema and users table
CREATE SCHEMA IF NOT EXISTS auth;

-- Create users table in auth schema
CREATE TABLE IF NOT EXISTS auth.users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email text UNIQUE NOT NULL,
    encrypted_password text NOT NULL,
    email_confirmed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT timezone('utc'::text, now()),
    updated_at timestamptz NOT NULL DEFAULT timezone('utc'::text, now())
);

-- Create public schema tables
CREATE TABLE IF NOT EXISTS public.user_folders (
    user_id uuid NOT NULL,
    name character varying NOT NULL,
    description text,
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_default boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT timezone('utc'::text, now()),
    updated_at timestamp with time zone DEFAULT timezone('utc'::text, now()),
    CONSTRAINT user_folders_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.user_content (
    user_id uuid NOT NULL,
    folder_id uuid,
    title text,
    content text NOT NULL,
    anchor text,
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_archived boolean DEFAULT false,
    is_favorite boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT timezone('utc'::text, now()),
    updated_at timestamp with time zone DEFAULT timezone('utc'::text, now()),
    search_vector tsvector,
    CONSTRAINT user_content_folder_id_fkey FOREIGN KEY (folder_id) REFERENCES public.user_folders(id) ON DELETE SET NULL,
    CONSTRAINT user_content_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_folders_user_id ON public.user_folders(user_id);
CREATE INDEX IF NOT EXISTS idx_user_content_user_id ON public.user_content(user_id);
CREATE INDEX IF NOT EXISTS idx_user_content_folder_id ON public.user_content(folder_id);
CREATE INDEX IF NOT EXISTS idx_user_content_search_vector ON public.user_content USING gin(search_vector);

-- Create triggers to update search_vector when content changes
CREATE OR REPLACE FUNCTION public.update_search_vector()
RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', COALESCE(NEW.title, '') || ' ' || COALESCE(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_search_vector ON public.user_content;
CREATE TRIGGER trigger_update_search_vector
    BEFORE INSERT OR UPDATE ON public.user_content
    FOR EACH ROW EXECUTE FUNCTION public.update_search_vector();
