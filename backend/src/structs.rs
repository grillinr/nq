use serde::{Deserialize, Serialize};
// =====================
// Data Structures
// =====================

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct User {
    pub user_id: String,
    pub name: String,
    pub email: String,
    pub auth_provider: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct MediaType {
    pub type_id: i32,
    pub type_name: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct MediaItem {
    pub media_id: String,
    pub title: String,
    pub type_id: i32,
    pub release_date: Option<String>,
    pub description: Option<String>,
    pub cover_url: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct CreatorRole {
    pub role_id: i32,
    pub role_name: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Creator {
    pub creator_id: String,
    pub name: String,
    pub role_id: i32,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Platform {
    pub platform_id: String,
    pub name: String,
    pub base_url: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ActivityStatus {
    pub status_id: i32,
    pub name: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct UserActivity {
    pub activity_id: String,
    pub user_id: String,
    pub media_id: String,
    pub status_id: i32,
    pub rating: Option<f64>,
    pub review: Option<String>,
    pub started_at: Option<String>,
    pub finished_at: Option<String>,
    pub source_platform: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Recommendation {
    pub recommendation_id: String,
    pub user_id: String,
    pub media_id: String,
    pub recommender_id: Option<String>,
    pub source: Option<String>,
    pub score: Option<f64>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Tag {
    pub tag_id: String,
    pub name: String,
    pub tag_type: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Rating {
    pub user_id: String,
    pub media_id: String,
    pub score: f64,
    pub rated_at: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Favorite {
    pub user_id: String,
    pub media_id: String,
    pub added_at: String,
}

// Request/Response structures
#[derive(Debug, Deserialize)]
pub struct CreateUserRequest {
    pub name: String,
    pub email: String,
    pub auth_provider: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct UpdateUserRequest {
    pub name: Option<String>,
    pub email: Option<String>,
    pub auth_provider: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct CreateMediaItemRequest {
    pub title: String,
    pub type_id: i32,
    pub release_date: Option<String>,
    pub description: Option<String>,
    pub cover_url: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct CreateRatingRequest {
    pub media_id: String,
    pub score: f64,
}

#[derive(Debug, Deserialize)]
pub struct CreateUserActivityRequest {
    pub media_id: String,
    pub status_id: i32,
    pub rating: Option<f64>,
    pub review: Option<String>,
    pub started_at: Option<String>,
    pub finished_at: Option<String>,
    pub source_platform: Option<String>,
}
