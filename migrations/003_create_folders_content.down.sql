-- Drop triggers
DROP TRIGGER IF EXISTS update_content_updated_at ON content;
DROP TRIGGER IF EXISTS update_folders_updated_at ON folders;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_content_user_id;
DROP INDEX IF EXISTS idx_content_folder_id;
DROP INDEX IF EXISTS idx_folders_user_id;
DROP INDEX IF EXISTS idx_folders_parent_id;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS content;
DROP TABLE IF EXISTS folders;
