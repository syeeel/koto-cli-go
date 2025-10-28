-- Migration: Add work_duration column for Pomodoro timer feature
-- This column stores the cumulative work time in minutes for each todo item
-- Default value is 0, which is safe for existing rows

ALTER TABLE todos ADD COLUMN work_duration INTEGER NOT NULL DEFAULT 0;
