-- Enhance learning_sessions table for AI-powered learning platform

-- Add new columns to learning_sessions table
ALTER TABLE learning_sessions ADD COLUMN IF NOT EXISTS current_step INTEGER DEFAULT 1;
ALTER TABLE learning_sessions ADD COLUMN IF NOT EXISTS concepts_completed TEXT[] DEFAULT '{}';
ALTER TABLE learning_sessions ADD COLUMN IF NOT EXISTS session_summary TEXT;
ALTER TABLE learning_sessions ADD COLUMN IF NOT EXISTS total_assessments INTEGER DEFAULT 0;
ALTER TABLE learning_sessions ADD COLUMN IF NOT EXISTS passed_assessments INTEGER DEFAULT 0;
ALTER TABLE learning_sessions ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'active';

-- Create index for status queries
CREATE INDEX IF NOT EXISTS idx_learning_sessions_status ON learning_sessions(status);

-- Add new columns to learning_prompts table for session tracking
ALTER TABLE learning_prompts ADD COLUMN IF NOT EXISTS session_id UUID REFERENCES learning_sessions(id) ON DELETE CASCADE;
ALTER TABLE learning_prompts ADD COLUMN IF NOT EXISTS step_number INTEGER;
ALTER TABLE learning_prompts ADD COLUMN IF NOT EXISTS interaction_type VARCHAR(50);
ALTER TABLE learning_prompts ADD COLUMN IF NOT EXISTS assessment_passed BOOLEAN;

-- Create indexes for learning_prompts
CREATE INDEX IF NOT EXISTS idx_learning_prompts_session_id ON learning_prompts(session_id);
CREATE INDEX IF NOT EXISTS idx_learning_prompts_step_number ON learning_prompts(step_number);
CREATE INDEX IF NOT EXISTS idx_learning_prompts_interaction_type ON learning_prompts(interaction_type);

-- Update existing learning_prompts to set interaction_type based on existing data
UPDATE learning_prompts 
SET interaction_type = 'initial_content' 
WHERE interaction_type IS NULL AND ai_response IS NOT NULL;

UPDATE learning_prompts 
SET interaction_type = 'user_question' 
WHERE interaction_type IS NULL AND ai_response IS NULL;