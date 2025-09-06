-- Rollback learning platform tables

-- Drop triggers first
DROP TRIGGER IF EXISTS update_user_preferences_updated_at ON user_preferences;
DROP TRIGGER IF EXISTS update_learning_prompts_updated_at ON learning_prompts;
DROP TRIGGER IF EXISTS update_user_learning_progress_updated_at ON user_learning_progress;
DROP TRIGGER IF EXISTS update_learning_methods_updated_at ON learning_methods;
DROP TRIGGER IF EXISTS update_topics_updated_at ON topics;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP INDEX IF EXISTS idx_learning_sessions_started_at;
DROP INDEX IF EXISTS idx_learning_sessions_user_id;
DROP INDEX IF EXISTS idx_learning_prompts_topic_id;
DROP INDEX IF EXISTS idx_learning_prompts_user_id;
DROP INDEX IF EXISTS idx_user_progress_status;
DROP INDEX IF EXISTS idx_user_progress_method_id;
DROP INDEX IF EXISTS idx_user_progress_topic_id;
DROP INDEX IF EXISTS idx_user_progress_user_id;
DROP INDEX IF EXISTS idx_learning_methods_method_key;
DROP INDEX IF EXISTS idx_topics_category;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS learning_sessions;
DROP TABLE IF EXISTS learning_prompts;
DROP TABLE IF EXISTS user_learning_progress;
DROP TABLE IF EXISTS learning_methods;
DROP TABLE IF EXISTS topics;