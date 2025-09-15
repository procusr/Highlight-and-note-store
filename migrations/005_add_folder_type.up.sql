-- Add folder_type column to folders table
ALTER TABLE folders ADD COLUMN folder_type VARCHAR(50) DEFAULT 'normal';

-- Create index for folder_type
CREATE INDEX idx_folders_type ON folders(folder_type);

-- Update existing folders to have normal type
UPDATE folders SET folder_type = 'normal' WHERE folder_type IS NULL;