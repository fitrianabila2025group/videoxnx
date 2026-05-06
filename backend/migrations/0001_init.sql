# Empty - GORM AutoMigrate handles schema. Optional manual SQL kept here.

-- Indexes that may help when traffic grows:
-- CREATE INDEX IF NOT EXISTS posts_published_idx ON posts (status, safety_status, COALESCE(published_at, scraped_at) DESC);
-- CREATE INDEX IF NOT EXISTS posts_view_idx ON posts (view_count DESC);
