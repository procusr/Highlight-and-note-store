-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_search_vector ON public.user_content;
DROP FUNCTION IF EXISTS public.update_search_vector();

-- Drop indexes
DROP INDEX IF EXISTS idx_user_content_search_vector;
DROP INDEX IF EXISTS idx_user_content_folder_id;
DROP INDEX IF EXISTS idx_user_content_user_id;
DROP INDEX IF EXISTS idx_user_folders_user_id;

-- Drop tables
DROP TABLE IF EXISTS public.user_content;
DROP TABLE IF EXISTS public.user_folders;
DROP TABLE IF EXISTS auth.users;

-- Drop schema
DROP SCHEMA IF EXISTS auth;
