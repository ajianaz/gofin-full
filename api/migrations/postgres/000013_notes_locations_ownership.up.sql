-- Add user_id and user_group_id to notes for ownership scoping
ALTER TABLE notes ADD COLUMN IF NOT EXISTS user_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
ALTER TABLE notes ADD COLUMN IF NOT EXISTS user_group_id UUID;

-- Add user_id and user_group_id to locations for ownership scoping
ALTER TABLE locations ADD COLUMN IF NOT EXISTS user_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
ALTER TABLE locations ADD COLUMN IF NOT EXISTS user_group_id UUID;
