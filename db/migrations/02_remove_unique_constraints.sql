-- Drop unique constraints from tables to allow reusing names after deletion

-- Drop unique constraint from projects table
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_name_key;

-- Drop unique constraint from environments table
ALTER TABLE environments DROP CONSTRAINT IF EXISTS environments_name_key;

-- Drop unique constraint from env_variables table
ALTER TABLE env_variables DROP CONSTRAINT IF EXISTS env_variables_project_id_environment_id_key_key;

-- Note: We'll handle uniqueness check at the application level for active records
