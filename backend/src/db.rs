use anyhow::Result;
use backend::structs::*;
use chrono::Utc;
use rusqlite::{Connection, params};
use std::sync::{Arc, Mutex};
use uuid::Uuid;

#[derive(Clone)]
pub struct Database {
    conn: Arc<Mutex<Connection>>,
}

impl Database {
    pub fn new(db_path: &str) -> Result<Self> {
        let conn = Connection::open(db_path)?;
        let db = Database {
            conn: Arc::new(Mutex::new(conn)),
        };
        db.init()?;
        Ok(db)
    }

    fn init(&self) -> Result<()> {
        let conn = self.conn.lock().unwrap();

        // Create all tables
        conn.execute_batch(
            "PRAGMA foreign_keys = ON;
            
            -- User table
            CREATE TABLE IF NOT EXISTS users (
                user_id TEXT PRIMARY KEY,
                name TEXT NOT NULL,
                email TEXT UNIQUE NOT NULL,
                auth_provider TEXT
            );
            
            -- MediaType table
            CREATE TABLE IF NOT EXISTS media_types (
                type_id INTEGER PRIMARY KEY AUTOINCREMENT,
                type_name TEXT NOT NULL UNIQUE
            );
            
            -- MediaItem table
            CREATE TABLE IF NOT EXISTS media_items (
                media_id TEXT PRIMARY KEY,
                title TEXT NOT NULL,
                type_id INTEGER NOT NULL,
                release_date TEXT,
                description TEXT,
                cover_url TEXT,
                FOREIGN KEY (type_id) REFERENCES media_types(type_id)
            );
            
            -- CreatorRole table
            CREATE TABLE IF NOT EXISTS creator_roles (
                role_id INTEGER PRIMARY KEY AUTOINCREMENT,
                role_name TEXT NOT NULL UNIQUE
            );
            
            -- Creator table
            CREATE TABLE IF NOT EXISTS creators (
                creator_id TEXT PRIMARY KEY,
                name TEXT NOT NULL,
                role_id INTEGER NOT NULL,
                FOREIGN KEY (role_id) REFERENCES creator_roles(role_id)
            );
            
            -- MediaCreator junction table
            CREATE TABLE IF NOT EXISTS media_creators (
                media_id TEXT NOT NULL,
                creator_id TEXT NOT NULL,
                PRIMARY KEY (media_id, creator_id),
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE,
                FOREIGN KEY (creator_id) REFERENCES creators(creator_id) ON DELETE CASCADE
            );
            
            -- Platform table
            CREATE TABLE IF NOT EXISTS platforms (
                platform_id TEXT PRIMARY KEY,
                name TEXT NOT NULL UNIQUE,
                base_url TEXT
            );
            
            -- MediaPlatform junction table
            CREATE TABLE IF NOT EXISTS media_platforms (
                media_id TEXT NOT NULL,
                platform_id TEXT NOT NULL,
                external_id TEXT,
                PRIMARY KEY (media_id, platform_id),
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE,
                FOREIGN KEY (platform_id) REFERENCES platforms(platform_id) ON DELETE CASCADE
            );
            
            -- ActivityStatus table
            CREATE TABLE IF NOT EXISTS activity_statuses (
                status_id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL UNIQUE
            );
            
            -- UserActivity table
            CREATE TABLE IF NOT EXISTS user_activities (
                activity_id TEXT PRIMARY KEY,
                user_id TEXT NOT NULL,
                media_id TEXT NOT NULL,
                status_id INTEGER NOT NULL,
                rating REAL,
                review TEXT,
                started_at TEXT,
                finished_at TEXT,
                source_platform TEXT,
                FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE,
                FOREIGN KEY (status_id) REFERENCES activity_statuses(status_id),
                FOREIGN KEY (source_platform) REFERENCES platforms(platform_id) ON DELETE SET NULL
            );
            
            -- Recommendation table
            CREATE TABLE IF NOT EXISTS recommendations (
                recommendation_id TEXT PRIMARY KEY,
                user_id TEXT NOT NULL,
                media_id TEXT NOT NULL,
                recommender_id TEXT,
                source TEXT,
                score REAL,
                FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE,
                FOREIGN KEY (recommender_id) REFERENCES users(user_id) ON DELETE SET NULL
            );
            
            -- Tag table
            CREATE TABLE IF NOT EXISTS tags (
                tag_id TEXT PRIMARY KEY,
                name TEXT NOT NULL,
                tag_type TEXT NOT NULL
            );
            
            -- MediaTag junction table
            CREATE TABLE IF NOT EXISTS media_tags (
                media_id TEXT NOT NULL,
                tag_id TEXT NOT NULL,
                PRIMARY KEY (media_id, tag_id),
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE,
                FOREIGN KEY (tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE
            );
            
            -- Rating table
            CREATE TABLE IF NOT EXISTS ratings (
                user_id TEXT NOT NULL,
                media_id TEXT NOT NULL,
                score REAL NOT NULL,
                rated_at TEXT NOT NULL,
                PRIMARY KEY (user_id, media_id),
                FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE
            );
            
            -- Favorite table
            CREATE TABLE IF NOT EXISTS favorites (
                user_id TEXT NOT NULL,
                media_id TEXT NOT NULL,
                added_at TEXT NOT NULL,
                PRIMARY KEY (user_id, media_id),
                FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE
            );
            
            -- ExternalIDs table
            CREATE TABLE IF NOT EXISTS external_ids (
                media_id TEXT NOT NULL,
                platform_id TEXT NOT NULL,
                external_id TEXT NOT NULL,
                PRIMARY KEY (media_id, platform_id),
                FOREIGN KEY (media_id) REFERENCES media_items(media_id) ON DELETE CASCADE,
                FOREIGN KEY (platform_id) REFERENCES platforms(platform_id) ON DELETE CASCADE
            );",
        )?;

        Ok(())
    }

    // =====================
    // User Operations
    // =====================

    pub fn create_user(&self, user: &CreateUserRequest) -> Result<User> {
        let conn = self.conn.lock().unwrap();
        let user_id = Uuid::new_v4().to_string();

        conn.execute(
            "INSERT INTO users (user_id, name, email, auth_provider) VALUES (?1, ?2, ?3, ?4)",
            params![user_id, user.name, user.email, user.auth_provider],
        )?;

        Ok(User {
            user_id,
            name: user.name.clone(),
            email: user.email.clone(),
            auth_provider: user.auth_provider.clone(),
        })
    }

    pub fn get_user(&self, id: &str) -> Result<Option<User>> {
        let conn = self.conn.lock().unwrap();
        let mut stmt = conn
            .prepare("SELECT user_id, name, email, auth_provider FROM users WHERE user_id = ?1")?;
        let mut user_iter = stmt.query_map(params![id], |row| {
            Ok(User {
                user_id: row.get(0)?,
                name: row.get(1)?,
                email: row.get(2)?,
                auth_provider: row.get(3)?,
            })
        })?;

        let user = user_iter.next().transpose()?;
        Ok(user)
    }

    pub fn get_all_users(&self) -> Result<Vec<User>> {
        let conn = self.conn.lock().unwrap();
        let mut stmt =
            conn.prepare("SELECT user_id, name, email, auth_provider FROM users ORDER BY name")?;
        let user_iter = stmt.query_map([], |row| {
            Ok(User {
                user_id: row.get(0)?,
                name: row.get(1)?,
                email: row.get(2)?,
                auth_provider: row.get(3)?,
            })
        })?;

        let mut users = Vec::new();
        for user_result in user_iter {
            users.push(user_result?);
        }
        Ok(users)
    }

    pub fn update_user(&self, id: &str, updates: &UpdateUserRequest) -> Result<Option<User>> {
        // First check if user exists
        let user = self.get_user(id)?;
        if user.is_none() {
            return Ok(None);
        }

        let conn = self.conn.lock().unwrap();

        if let Some(name) = &updates.name {
            conn.execute(
                "UPDATE users SET name = ?1 WHERE user_id = ?2",
                params![name, id],
            )?;
        }

        if let Some(email) = &updates.email {
            conn.execute(
                "UPDATE users SET email = ?1 WHERE user_id = ?2",
                params![email, id],
            )?;
        }

        if let Some(auth_provider) = &updates.auth_provider {
            conn.execute(
                "UPDATE users SET auth_provider = ?1 WHERE user_id = ?2",
                params![auth_provider, id],
            )?;
        }

        // Return updated user
        self.get_user(id)
    }

    pub fn delete_user(&self, id: &str) -> Result<bool> {
        let conn = self.conn.lock().unwrap();
        let rows_affected = conn.execute("DELETE FROM users WHERE user_id = ?1", params![id])?;

        Ok(rows_affected > 0)
    }

    // =====================
    // Media Operations
    // =====================

    pub fn create_media_item(&self, media: &CreateMediaItemRequest) -> Result<MediaItem> {
        let conn = self.conn.lock().unwrap();
        let media_id = Uuid::new_v4().to_string();

        conn.execute(
            "INSERT INTO media_items (media_id, title, type_id, release_date, description, cover_url) VALUES (?1, ?2, ?3, ?4, ?5, ?6)",
            params![media_id, media.title, media.type_id, media.release_date, media.description, media.cover_url],
        )?;

        Ok(MediaItem {
            media_id,
            title: media.title.clone(),
            type_id: media.type_id,
            release_date: media.release_date.clone(),
            description: media.description.clone(),
            cover_url: media.cover_url.clone(),
        })
    }

    pub fn get_media_item(&self, id: &str) -> Result<Option<MediaItem>> {
        let conn = self.conn.lock().unwrap();
        let mut stmt = conn.prepare("SELECT media_id, title, type_id, release_date, description, cover_url FROM media_items WHERE media_id = ?1")?;
        let mut media_iter = stmt.query_map(params![id], |row| {
            Ok(MediaItem {
                media_id: row.get(0)?,
                title: row.get(1)?,
                type_id: row.get(2)?,
                release_date: row.get(3)?,
                description: row.get(4)?,
                cover_url: row.get(5)?,
            })
        })?;

        let media = media_iter.next().transpose()?;
        Ok(media)
    }

    pub fn get_all_media_items(&self) -> Result<Vec<MediaItem>> {
        let conn = self.conn.lock().unwrap();
        let mut stmt = conn.prepare("SELECT media_id, title, type_id, release_date, description, cover_url FROM media_items ORDER BY title")?;
        let media_iter = stmt.query_map([], |row| {
            Ok(MediaItem {
                media_id: row.get(0)?,
                title: row.get(1)?,
                type_id: row.get(2)?,
                release_date: row.get(3)?,
                description: row.get(4)?,
                cover_url: row.get(5)?,
            })
        })?;

        let mut media_items = Vec::new();
        for media_result in media_iter {
            media_items.push(media_result?);
        }
        Ok(media_items)
    }

    // =====================
    // Rating Operations
    // =====================

    pub fn create_rating(&self, user_id: &str, rating: &CreateRatingRequest) -> Result<Rating> {
        let conn = self.conn.lock().unwrap();
        let now = Utc::now()
            .naive_utc()
            .format("%Y-%m-%d %H:%M:%S")
            .to_string();

        conn.execute(
            "INSERT OR REPLACE INTO ratings (user_id, media_id, score, rated_at) VALUES (?1, ?2, ?3, ?4)",
            params![user_id, rating.media_id, rating.score, now],
        )?;

        Ok(Rating {
            user_id: user_id.to_string(),
            media_id: rating.media_id.clone(),
            score: rating.score,
            rated_at: now,
        })
    }

    pub fn get_user_rating(&self, user_id: &str, media_id: &str) -> Result<Option<Rating>> {
        let conn = self.conn.lock().unwrap();
        let mut stmt = conn.prepare("SELECT user_id, media_id, score, rated_at FROM ratings WHERE user_id = ?1 AND media_id = ?2")?;
        let mut rating_iter = stmt.query_map(params![user_id, media_id], |row| {
            Ok(Rating {
                user_id: row.get(0)?,
                media_id: row.get(1)?,
                score: row.get(2)?,
                rated_at: row.get(3)?,
            })
        })?;

        let rating = rating_iter.next().transpose()?;
        Ok(rating)
    }

    // =====================
    // User Activity Operations
    // =====================

    pub fn create_user_activity(
        &self,
        user_id: &str,
        activity: &CreateUserActivityRequest,
    ) -> Result<UserActivity> {
        let conn = self.conn.lock().unwrap();
        let activity_id = Uuid::new_v4().to_string();

        conn.execute(
            "INSERT INTO user_activities (activity_id, user_id, media_id, status_id, rating, review, started_at, finished_at, source_platform) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)",
            params![activity_id, user_id, activity.media_id, activity.status_id, activity.rating, activity.review, activity.started_at, activity.finished_at, activity.source_platform],
        )?;

        Ok(UserActivity {
            activity_id,
            user_id: user_id.to_string(),
            media_id: activity.media_id.clone(),
            status_id: activity.status_id,
            rating: activity.rating,
            review: activity.review.clone(),
            started_at: activity.started_at.clone(),
            finished_at: activity.finished_at.clone(),
            source_platform: activity.source_platform.clone(),
        })
    }

    pub fn get_user_activities(&self, user_id: &str) -> Result<Vec<UserActivity>> {
        let conn = self.conn.lock().unwrap();
        let mut stmt = conn.prepare("SELECT activity_id, user_id, media_id, status_id, rating, review, started_at, finished_at, source_platform FROM user_activities WHERE user_id = ?1 ORDER BY started_at DESC")?;
        let activity_iter = stmt.query_map(params![user_id], |row| {
            Ok(UserActivity {
                activity_id: row.get(0)?,
                user_id: row.get(1)?,
                media_id: row.get(2)?,
                status_id: row.get(3)?,
                rating: row.get(4)?,
                review: row.get(5)?,
                started_at: row.get(6)?,
                finished_at: row.get(7)?,
                source_platform: row.get(8)?,
            })
        })?;

        let mut activities = Vec::new();
        for activity_result in activity_iter {
            activities.push(activity_result?);
        }
        Ok(activities)
    }

    // =====================
    // Seed Data
    // =====================

    pub fn seed_sample_data(&self) -> Result<()> {
        let conn = self.conn.lock().unwrap();

        // Check if we already have users
        let count: i64 = conn.query_row("SELECT COUNT(*) FROM users", [], |row| row.get(0))?;

        if count == 0 {
            // Seed media types
            conn.execute(
                "INSERT OR IGNORE INTO media_types (type_name) VALUES ('Movie')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO media_types (type_name) VALUES ('TV Show')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO media_types (type_name) VALUES ('Book')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO media_types (type_name) VALUES ('Game')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO media_types (type_name) VALUES ('Music')",
                [],
            )?;

            // Seed creator roles
            conn.execute(
                "INSERT OR IGNORE INTO creator_roles (role_name) VALUES ('Director')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO creator_roles (role_name) VALUES ('Actor')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO creator_roles (role_name) VALUES ('Author')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO creator_roles (role_name) VALUES ('Developer')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO creator_roles (role_name) VALUES ('Artist')",
                [],
            )?;

            // Seed activity statuses
            conn.execute(
                "INSERT OR IGNORE INTO activity_statuses (name) VALUES ('Want to Watch/Read/Play')",
                [],
            )?;
            conn.execute("INSERT OR IGNORE INTO activity_statuses (name) VALUES ('Currently Watching/Reading/Playing')", [])?;
            conn.execute(
                "INSERT OR IGNORE INTO activity_statuses (name) VALUES ('Completed')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO activity_statuses (name) VALUES ('Dropped')",
                [],
            )?;
            conn.execute(
                "INSERT OR IGNORE INTO activity_statuses (name) VALUES ('On Hold')",
                [],
            )?;

            // Seed platforms
            let netflix_id = Uuid::new_v4().to_string();
            let amazon_id = Uuid::new_v4().to_string();
            let steam_id = Uuid::new_v4().to_string();

            conn.execute("INSERT OR IGNORE INTO platforms (platform_id, name, base_url) VALUES (?1, 'Netflix', 'https://netflix.com')", params![netflix_id])?;
            conn.execute("INSERT OR IGNORE INTO platforms (platform_id, name, base_url) VALUES (?1, 'Amazon Prime', 'https://amazon.com/prime')", params![amazon_id])?;
            conn.execute("INSERT OR IGNORE INTO platforms (platform_id, name, base_url) VALUES (?1, 'Steam', 'https://steam.com')", params![steam_id])?;

            // Seed sample users
            let sample_users: Vec<(&str, &str, Option<&str>)> = vec![
                ("John Doe", "john.doe@example.com", None),
                ("Jane Smith", "jane.smith@example.com", None),
                ("Bob Johnson", "bob.johnson@example.com", None),
            ];

            for (name, email, auth_provider) in &sample_users {
                let user_id = Uuid::new_v4().to_string();
                conn.execute(
                    "INSERT INTO users (user_id, name, email, auth_provider) VALUES (?1, ?2, ?3, ?4)",
                    params![user_id, name, email, auth_provider],
                )?;
            }

            log::info!(
                "Seeded database with {} sample users and reference data",
                sample_users.len()
            );
        }

        Ok(())
    }
}
