# Subject: Image Processing Service in Go

Goal
- Build a scalable backend system for image processing similar to Cloudinary, featuring user authentication, image uploads, and various image transformations.

You'll Practice
- User authentication with JWT tokens
- Multipart form-data file uploads
- Image processing and transformation operations (resize, crop, rotate, filters, etc.)
- Cloud storage integration (AWS S3, Cloudflare R2, or Google Cloud Storage)
- API design with RESTful endpoints
- Input validation and error handling
- Rate limiting and caching strategies
- Optional: asynchronous processing with message queues (RabbitMQ, Kafka)

Run
1) cd "Image Process Service"
2) go run .
3) API at http://localhost:8080

Authentication Endpoints
- POST   /register                → register a new user
- POST   /login                   → login and get JWT token

Image Management Endpoints
- POST   /images                  → upload an image
- GET    /images                  → list all user images (paginated)
- GET    /images/{id}             → retrieve a specific image
- POST   /images/{id}/transform   → apply transformations to an image
- DELETE /images/{id}             → delete an image

Request/Response Examples

Register:
```
POST /register
{
  "username": "user1",
  "password": "password123"
}
Response: { "id": "uuid", "username": "user1", "jwt": "token" }
```

Login:
```
POST /login
{
  "username": "user1",
  "password": "password123"
}
Response: { "id": "uuid", "username": "user1", "jwt": "token" }
```

Upload Image:
```
POST /images
Content-Type: multipart/form-data
File: [binary image data]
Response: { "id": "img_uuid", "filename": "photo.jpg", "url": "...", "size": 2048, "width": 1920, "height": 1080 }
```

Apply Transformations:
```
POST /images/{id}/transform
{
  "transformations": {
    "resize": { "width": 800, "height": 600 },
    "crop": { "width": 400, "height": 300, "x": 0, "y": 0 },
    "rotate": 90,
    "format": "webp",
    "filters": {
      "grayscale": false,
      "sepia": true
    }
  }
}
Response: { "id": "transform_uuid", "original_id": "img_uuid", "url": "...", "transformations": {...} }
```

List Images (Paginated):
```
GET /images?page=1&limit=10
Response: {
  "images": [...],
  "page": 1,
  "limit": 10,
  "total": 45
}
```

Available Image Transformations
- Resize: Scale image to specific dimensions
- Crop: Extract a portion of the image
- Rotate: Rotate image by angle (0-360 degrees)
- Flip: Horizontal or vertical flip
- Mirror: Create mirrored effect
- Compress: Reduce file size while maintaining quality
- Format Conversion: Convert between JPEG, PNG, WebP, GIF, etc.
- Filters: Grayscale, sepia, blur, brightness, contrast, etc.

Data Model

User:
- id: UUID
- username: string (unique)
- password_hash: string (bcrypt)
- created_at: timestamp

Image:
- id: UUID
- user_id: UUID (foreign key)
- filename: string
- original_url: string (cloud storage URL)
- width: int
- height: int
- size: int (bytes)
- format: string (JPEG, PNG, etc.)
- created_at: timestamp

Transformation:
- id: UUID
- user_id: UUID (foreign key)
- original_image_id: UUID (foreign key)
- transformations: JSON (applied operations)
- result_url: string (cloud storage URL)
- created_at: timestamp

Behavior
- All endpoints except /register and /login require JWT in Authorization header
- Images are stored in cloud storage (S3, R2, or GCS)
- Transformations can be chained (multiple applied in sequence)
- Rate limiting: max 100 transformations per hour per user
- Transformed images are cached for 24 hours
- Invalid file types are rejected (accept: image/*, video)
- File size limit: 50MB per image
- HTTP responses use application/json content-type

Try
Authentication:
```
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}'

curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}'
```

Upload Image:
```
curl -X POST http://localhost:8080/images \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -F "file=@/path/to/image.jpg"
```

List Images:
```
curl http://localhost:8080/images?page=1&limit=10 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Apply Transformation:
```
curl -X POST http://localhost:8080/images/<IMAGE_ID>/transform \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "transformations": {
      "resize": {"width": 800, "height": 600},
      "format": "webp"
    }
  }'
```

Project Structure
- main.go                → application entry point
- handlers/              → HTTP request handlers
- models/               → data structures
- middleware/           → JWT authentication, error handling
- storage/              → cloud storage integration
- processor/            → image transformation logic
- config/               → configuration management
- go.mod               → module dependencies

Key Libraries
- github.com/gorilla/mux → HTTP routing
- github.com/golang-jwt/jwt → JWT authentication
- golang.org/x/crypto/bcrypt → password hashing
- image/*, image/draw → image processing
- aws/aws-sdk-go-v2 → AWS S3 integration (or similar for other cloud providers)
