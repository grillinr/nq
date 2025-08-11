use actix_web::{App, HttpResponse, HttpServer, Result, web};
use log;

mod db;
use db::{
    CreateMediaItemRequest, CreateRatingRequest, CreateUserActivityRequest, CreateUserRequest,
    Database, UpdateUserRequest,
};

// Database state
type AppState = web::Data<Database>;

// Create a new user
async fn create_user(
    user_data: web::Json<CreateUserRequest>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    match db.create_user(&user_data) {
        Ok(user) => Ok(HttpResponse::Created().json(user)),
        Err(e) => {
            log::error!("Failed to create user: {}", e);
            if e.to_string().contains("UNIQUE constraint failed") {
                Ok(HttpResponse::Conflict().json(serde_json::json!({
                    "error": "User with this email already exists"
                })))
            } else {
                Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": "Failed to create user"
                })))
            }
        }
    }
}

// Get a user by ID
async fn get_user(path: web::Path<String>, db: AppState) -> Result<HttpResponse, actix_web::Error> {
    let user_id = path.into_inner();

    match db.get_user(&user_id) {
        Ok(Some(user)) => Ok(HttpResponse::Ok().json(user)),
        Ok(None) => Ok(HttpResponse::NotFound().json(serde_json::json!({
            "error": "User not found"
        }))),
        Err(e) => {
            log::error!("Failed to get user: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to retrieve user"
            })))
        }
    }
}

// Update a user
async fn update_user(
    path: web::Path<String>,
    update_data: web::Json<UpdateUserRequest>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    let user_id = path.into_inner();

    match db.update_user(&user_id, &update_data) {
        Ok(Some(user)) => Ok(HttpResponse::Ok().json(user)),
        Ok(None) => Ok(HttpResponse::NotFound().json(serde_json::json!({
            "error": "User not found"
        }))),
        Err(e) => {
            log::error!("Failed to update user: {}", e);
            if e.to_string().contains("UNIQUE constraint failed") {
                Ok(HttpResponse::Conflict().json(serde_json::json!({
                    "error": "User with this email already exists"
                })))
            } else {
                Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": "Failed to update user"
                })))
            }
        }
    }
}

// Delete a user
async fn delete_user(
    path: web::Path<String>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    let user_id = path.into_inner();

    match db.delete_user(&user_id) {
        Ok(true) => Ok(HttpResponse::NoContent().finish()),
        Ok(false) => Ok(HttpResponse::NotFound().json(serde_json::json!({
            "error": "User not found"
        }))),
        Err(e) => {
            log::error!("Failed to delete user: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to delete user"
            })))
        }
    }
}

// List all users
async fn list_users(db: AppState) -> Result<HttpResponse, actix_web::Error> {
    match db.get_all_users() {
        Ok(users) => Ok(HttpResponse::Ok().json(users)),
        Err(e) => {
            log::error!("Failed to get users: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to retrieve users"
            })))
        }
    }
}

// Create a new media item
async fn create_media_item(
    media_data: web::Json<CreateMediaItemRequest>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    match db.create_media_item(&media_data) {
        Ok(media) => Ok(HttpResponse::Created().json(media)),
        Err(e) => {
            log::error!("Failed to create media item: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to create media item"
            })))
        }
    }
}

// Get a media item by ID
async fn get_media_item(
    path: web::Path<String>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    let media_id = path.into_inner();

    match db.get_media_item(&media_id) {
        Ok(Some(media)) => Ok(HttpResponse::Ok().json(media)),
        Ok(None) => Ok(HttpResponse::NotFound().json(serde_json::json!({
            "error": "Media item not found"
        }))),
        Err(e) => {
            log::error!("Failed to get media item: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to retrieve media item"
            })))
        }
    }
}

// List all media items
async fn list_media_items(db: AppState) -> Result<HttpResponse, actix_web::Error> {
    match db.get_all_media_items() {
        Ok(media_items) => Ok(HttpResponse::Ok().json(media_items)),
        Err(e) => {
            log::error!("Failed to get media items: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to retrieve media items"
            })))
        }
    }
}

// Create a rating
async fn create_rating(
    path: web::Path<String>,
    rating_data: web::Json<CreateRatingRequest>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    let user_id = path.into_inner();

    match db.create_rating(&user_id, &rating_data) {
        Ok(rating) => Ok(HttpResponse::Created().json(rating)),
        Err(e) => {
            log::error!("Failed to create rating: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to create rating"
            })))
        }
    }
}

// Get user activities
async fn get_user_activities(
    path: web::Path<String>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    let user_id = path.into_inner();

    match db.get_user_activities(&user_id) {
        Ok(activities) => Ok(HttpResponse::Ok().json(activities)),
        Err(e) => {
            log::error!("Failed to get user activities: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to retrieve user activities"
            })))
        }
    }
}

// Create user activity
async fn create_user_activity(
    path: web::Path<String>,
    activity_data: web::Json<CreateUserActivityRequest>,
    db: AppState,
) -> Result<HttpResponse, actix_web::Error> {
    let user_id = path.into_inner();

    match db.create_user_activity(&user_id, &activity_data) {
        Ok(activity) => Ok(HttpResponse::Created().json(activity)),
        Err(e) => {
            log::error!("Failed to create user activity: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to create user activity"
            })))
        }
    }
}

// Health check endpoint
async fn health_check() -> Result<HttpResponse> {
    Ok(HttpResponse::Ok().json(serde_json::json!({
        "status": "healthy",
        "service": "Cross-Media Tracking Platform API",
        "database": "SQLite"
    })))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    // Initialize SQLite database
    let db_path = "data/app.db";

    // Ensure data directory exists
    std::fs::create_dir_all("data").unwrap_or_else(|_| {
        log::warn!("Could not create data directory, using current directory");
    });

    let database = match Database::new(db_path) {
        Ok(db) => {
            log::info!("Connected to SQLite database at {}", db_path);
            db
        }
        Err(e) => {
            log::error!("Failed to connect to database: {}", e);
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                format!("Database connection failed: {}", e),
            ));
        }
    };

    // Seed the database with sample data
    if let Err(e) = database.seed_sample_data() {
        log::warn!("Failed to seed sample data: {}", e);
    }

    log::info!("Starting Cross-Media Tracking Platform API server at http://127.0.0.1:8080");
    log::info!("Database: SQLite at {}", db_path);

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(database.clone()))
            .route("/health", web::get().to(health_check))
            // User endpoints
            .route("/users", web::get().to(list_users))
            .route("/users", web::post().to(create_user))
            .route("/users/{id}", web::get().to(get_user))
            .route("/users/{id}", web::patch().to(update_user))
            .route("/users/{id}", web::delete().to(delete_user))
            // Media endpoints
            .route("/media", web::get().to(list_media_items))
            .route("/media", web::post().to(create_media_item))
            .route("/media/{id}", web::get().to(get_media_item))
            // Rating endpoints
            .route("/users/{id}/ratings", web::post().to(create_rating))
            // User activity endpoints
            .route("/users/{id}/activities", web::get().to(get_user_activities))
            .route(
                "/users/{id}/activities",
                web::post().to(create_user_activity),
            )
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
