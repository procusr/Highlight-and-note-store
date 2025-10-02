-- Rollback: Rename learning_goal column back to user_prompt
ALTER TABLE learning_prompts
RENAME COLUMN learning_goal TO user_prompt;
