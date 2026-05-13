# SnipQURL Technical Documentation

## Overview
SnipQURL is a high-performance URL shortening service designed for efficiency and reliability. The platform provides automated URL normalization, customizable aliases, and managed link lifecycle through enforced expiration and background cleanup processes.

## Key Features

### URL Protocol Recognition and Normalization
The system automatically normalizes input URLs to ensure consistency. If a user provides a destination link without a protocol scheme (e.g., "google.com"), the platform automatically prepends "http://" before processing. This ensures that all shortened links are valid and functional without requiring manual protocol entry from the end-user.

### Custom Shorten Link Aliases
Users have the option to specify custom aliases for their shortened links.
- **Maximum Length:** 30 characters.
- **Uniqueness:** The system enforces strict uniqueness for aliases. If a requested alias is already in use, the system will return a descriptive error message.
- **Random Fallback:** If no alias is provided, the system generates a secure, random 8-character code.

### Managed Link Expiration
To maintain database performance and ensure the availability of short codes, SnipQURL enforces a mandatory expiration policy.
- **Configurable Expiration:** Users can choose from several intervals: 1 Hour, 1 Day, 1 Week, or 1 Month.
- **Default Expiration:** Requests that do not specify an expiration interval are automatically assigned a default lifetime of 30 days.
- **Expiration Behavior:** Once a link reaches its expiration timestamp, it becomes inaccessible. Any attempt to access an expired link will result in an HTTP 410 (Gone) status code.

### Automated Background Cleanup
The platform includes an automated background worker that executes every 60 minutes. This worker identifies and permanently prunes expired records from the database, ensuring that storage is optimized and expired aliases are freed for future use.

## API Specification

### Shorten URL
`POST /api/shorten`

**Request Body (JSON):**
| Field | Type | Description |
| :--- | :--- | :--- |
| `url` | string | The destination URL (Required). |
| `alias` | string | Preferred custom alias (Optional). |
| `expires_in` | string | Expiration duration (e.g., "1h", "24h", "720h"). Defaults to "720h" (Optional). |

**Response (JSON):**
| Field | Type | Description |
| :--- | :--- | :--- |
| `code` | string | The fully qualified shortened URL. |

### Generate QR Code
`POST /api/qr`

**Request Body (JSON):**
| Field | Type | Description |
| :--- | :--- | :--- |
| `url` | string | The URL to encode into the QR code (Required). |

## Deployment and Local Testing

### Prerequisites
- Docker
- Docker Compose

### Local Setup
1. Clone the repository.
2. Initialize the environment:
   ```bash
   docker-compose up -d postgres
   ```
3. Apply migrations:
   ```bash
   docker-compose exec -T postgres psql -U snipqurl -d snipqurl < db/migrations/000001_create_urls_table.up.sql
   docker-compose exec -T postgres psql -U snipqurl -d snipqurl < db/migrations/000002_add_alias_and_expiration.up.sql
   ```
4. Start the application:
   ```bash
   docker-compose up --build app
   ```
The application will be accessible at `http://localhost:8080`.
