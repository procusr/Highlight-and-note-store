-- Rename user_prompt column to learning_goal in learning_prompts table
ALTER TABLE learning_prompts
RENAME COLUMN user_prompt TO learning_goal;
