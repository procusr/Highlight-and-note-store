-- Rollback enhancements to learning_sessions table

-- Remove indexes
DROP INDEX IF EXISTS idx_learning_sessions_status;
DROP INDEX IF EXISTS idx_learning_prompts_session_id;
DROP INDEX IF EXISTS idx_learning_prompts_step_number;
DROP INDEX IF EXISTS idx_learning_prompts_interaction_type;

-- Remove columns from learning_prompts
ALTER TABLE learning_prompts DROP COLUMN IF EXISTS session_id;
ALTER TABLE learning_prompts DROP COLUMN IF EXISTS step_number;
ALTER TABLE learning_prompts DROP COLUMN IF EXISTS interaction_type;
ALTER TABLE learning_prompts DROP COLUMN IF EXISTS assessment_passed;

-- Remove columns from learning_sessions
ALTER TABLE learning_sessions DROP COLUMN IF EXISTS current_step;
ALTER TABLE learning_sessions DROP COLUMN IF EXISTS concepts_completed;
ALTER TABLE learning_sessions DROP COLUMN IF EXISTS session_summary;
ALTER TABLE learning_sessions DROP COLUMN IF EXISTS total_assessments;
ALTER TABLE learning_sessions DROP COLUMN IF EXISTS passed_assessments;
ALTER TABLE learning_sessions DROP COLUMN IF EXISTS status;