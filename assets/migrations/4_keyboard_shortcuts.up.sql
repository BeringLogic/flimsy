
ALTER TABLE item ADD COLUMN shortcut TEXT;
UPDATE item SET shortcut = "";
