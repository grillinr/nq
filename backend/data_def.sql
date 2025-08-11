-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================
-- User
-- =====================
CREATE TABLE "User" (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    auth_provider TEXT
);

-- =====================
-- MediaType
-- =====================
CREATE TABLE MediaType (
    type_id SERIAL PRIMARY KEY,
    type_name TEXT NOT NULL UNIQUE
);

-- =====================
-- MediaItem
-- =====================
CREATE TABLE MediaItem (
    media_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    type_id INT NOT NULL REFERENCES MediaType(type_id) ON DELETE RESTRICT,
    release_date DATE,
    description TEXT,
    cover_url TEXT
);

-- =====================
-- CreatorRole
-- =====================
CREATE TABLE CreatorRole (
    role_id SERIAL PRIMARY KEY,
    role_name TEXT NOT NULL UNIQUE
);

-- =====================
-- Creator
-- =====================
CREATE TABLE Creator (
    creator_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    role_id INT NOT NULL REFERENCES CreatorRole(role_id) ON DELETE RESTRICT
);

-- =====================
-- MediaCreator (M:N between MediaItem and Creator)
-- =====================
CREATE TABLE MediaCreator (
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    creator_id UUID NOT NULL REFERENCES Creator(creator_id) ON DELETE CASCADE,
    PRIMARY KEY (media_id, creator_id)
);

-- =====================
-- Platform
-- =====================
CREATE TABLE Platform (
    platform_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    base_url TEXT
);

-- =====================
-- MediaPlatform (M:N between MediaItem and Platform)
-- =====================
CREATE TABLE MediaPlatform (
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    platform_id UUID NOT NULL REFERENCES Platform(platform_id) ON DELETE CASCADE,
    external_id TEXT,
    PRIMARY KEY (media_id, platform_id)
);

-- =====================
-- ActivityStatus
-- =====================
CREATE TABLE ActivityStatus (
    status_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- =====================
-- UserActivity
-- =====================
CREATE TABLE UserActivity (
    activity_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "User"(user_id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    status_id INT NOT NULL REFERENCES ActivityStatus(status_id) ON DELETE RESTRICT,
    rating NUMERIC,
    review TEXT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    source_platform UUID REFERENCES Platform(platform_id) ON DELETE SET NULL
);

-- =====================
-- Recommendation
-- =====================
CREATE TABLE Recommendation (
    recommendation_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "User"(user_id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    recommender_id UUID REFERENCES "User"(user_id) ON DELETE SET NULL,
    source TEXT,
    score FLOAT
);

-- =====================
-- Tag
-- =====================
CREATE TABLE Tag (
    tag_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    type TEXT NOT NULL
);

-- =====================
-- MediaTag (M:N between MediaItem and Tag)
-- =====================
CREATE TABLE MediaTag (
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES Tag(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (media_id, tag_id)
);

-- =====================
-- Rating (standalone)
-- =====================
CREATE TABLE Rating (
    user_id UUID NOT NULL REFERENCES "User"(user_id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    score NUMERIC NOT NULL,
    rated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, media_id)
);

-- =====================
-- Favorite
-- =====================
CREATE TABLE Favorite (
    user_id UUID NOT NULL REFERENCES "User"(user_id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, media_id)
);

-- =====================
-- Optional: ExternalIDs (if separate from MediaPlatform)
-- =====================
CREATE TABLE ExternalIDs (
    media_id UUID NOT NULL REFERENCES MediaItem(media_id) ON DELETE CASCADE,
    platform_id UUID NOT NULL REFERENCES Platform(platform_id) ON DELETE CASCADE,
    external_id TEXT NOT NULL,
    PRIMARY KEY (media_id, platform_id)
);
