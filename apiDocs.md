# API Documentation - Certify NFT Backend

## Base URL
```
https://api.gpadaka.com/api3
```

## Authentication
Currently, the API uses wallet address-based authentication. Users and vendors are identified by their wallet addresses.

---

## Endpoint Summary

### Authentication
- **`POST /api/auth/login`**: Authenticates a user or vendor.

### Users
- **`POST /api/users/register`**: Registers a new user account.
- **`GET /api/users/:walletAddress/events`**: Retrieves all events for a specific user.
- **`GET /api/users/:walletAddress/certificate`**: Retrieves all certificates for a specific user.
- **`POST /api/users/whitelist`**: Creates a whitelist entry for an event.
- **`POST /api/users/whitelist/cancel`**: Cancels a user's whitelist entry for an event.

### Vendors
- **`POST /api/vendors/register`**: Registers a new vendor account.
- **`GET /api/vendors/:walletAddress/events`**: Retrieves all events for a specific vendor.

### Events
- **`GET /api/events/all`**: Retrieves all events.
- **`POST /api/events/create`**: Creates a new event.
- **`GET /api/events/:id`**: Retrieves detailed information about a specific event.

---

## Endpoints

### Auth

#### 1. Login
**POST** `/api/auth/login`

Authenticates a user or vendor using their wallet address.

#### Request
**Content-Type:** `application/json`

```json
{
  "wallet_address": "0x1234567890abcdef..."
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| wallet_address | string | Yes | User's or vendor's wallet address |

#### Response
**Success (200 OK) - Existing User**
```json
{
  "isNewUser": false,
  "wallet_address": "0x1234567890abcdef...",
  "role": "users"
}
```

**Success (200 OK) - Existing Vendor**
```json
{
  "isNewUser": false,
  "wallet_address": "0x1234567890abcdef...",
  "role": "vendors"
}
```

**Success (200 OK) - New User**
```json
{
  "isNewUser": true
}
```

**Error Responses**

**400 Bad Request - Invalid Request**
```json
{
  "error": "Invalid request"
}
```

**400 Bad Request - Missing Wallet Address**
```json
{
  "error": "wallet_address is required"
}
```

**500 Internal Server Error**
```json
{
  "error": "database error message"
}
```

---

### Users

#### 2. Register User
**POST** `/api/users/register`

Registers a new user account.

#### Request
**Content-Type:** `application/json`

```json
{
  "email": "user@example.com",
  "wallet_address": "0x1234567890abcdef...",
  "name": "John Doe"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | No | User's email address |
| wallet_address | string | Yes | User's wallet address (must be unique) |
| name | string | Yes | User's full name |

#### Response
**Success (201 Created)**
```json
{
  "id": 1,
  "email": "user@example.com",
  "wallet_address": "0x1234567890abcdef...",
  "name": "John Doe",
  "created_at": "2024-06-10T10:00:00Z",
  "updated_at": "2024-06-10T10:00:00Z"
}
```

**Error Responses**

**400 Bad Request - Invalid Request**
```json
{
  "error": "invalid request"
}
```

**400 Bad Request - Missing Required Fields**
```json
{
  "error": "name and wallet_address are required"
}
```

**409 Conflict - Wallet Already Registered**
```json
{
  "error": "wallet address already registered with another account"
}
```

**500 Internal Server Error**
```json
{
  "error": "server error"
}
```

---

#### 3. Get Events by User
**GET** `/api/users/:walletAddress/events`

Retrieves all events a specific user is registered for, based on their wallet address.

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| walletAddress | string | Yes | The user's wallet address. |

#### Response
**Success (200 OK)**
Returns an array of `EventWithUserStatus` objects.
```json
[
  {
    "id": 1,
    "title": "NFT Conference 2024",
    "description": "Annual NFT conference",
    "vendor_id": 1,
    "start_date": "2024-06-15T09:00:00Z",
    "end_date": "2024-06-15T17:00:00Z",
    "status": "upcoming",
    "created_at": "2024-06-10T10:00:00Z",
    "updated_at": "2024-06-10T10:00:00Z",
    "picture": "uploads/1718000000_event.jpg",
    "maxattendees": 100,
    "location": "Jakarta Convention Center",
    "attendees": 25,
    "user_status": "present",
    "requirements": {"items": ["Laptop", "Notebook"]},
    "agenda": {"sessions": [{"time": "09:00", "topic": "Introduction"}]}
  }
]
```

**Error Responses**

**400 Bad Request - Missing Wallet Address**
```json
{
  "error": "walletAddress is required"
}
```

**404 Not Found**
```json
{
  "error": "user not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "database error message"
}
```

---

#### 4. Get Certificates by User
**GET** `/api/users/:walletAddress/certificate`

Retrieves all certificates for a specific user, combined with event details.

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| walletAddress | string | Yes | The user's wallet address. |

#### Response
**Success (200 OK)**
Returns an array of `CertificateWithEvent` objects.
```json
[
  {
    "id": 1,
    "event_id": 12,
    "user_id": 1,
    "certificate_data": "ipfs://bafybeig.../metadata.json",
    "mint_status": "minted",
    "mint_transaction_hash": "0xabc...",
    "created_at": "2024-06-16T10:00:00Z",
    "updated_at": "2024-06-16T10:00:00Z",
    "event_title": "NFT Conference 2024",
    "event_description": "Annual NFT conference",
    "event_start_date": "2024-06-15T09:00:00Z",
    "event_location": "Jakarta Convention Center",
    "event_picture": "https://api.gpadaka.com/uploads/1718000000_event.jpg"
  }
]
```

**Error Responses**

**400 Bad Request - Missing Wallet Address**
```json
{
  "error": "walletAddress is required"
}
```

**404 Not Found**
```json
{
  "error": "user not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "database error message"
}
```

---

#### 5. Create Whitelist Entry
**POST** `/api/users/whitelist`

Creates a whitelist entry for a user to attend an event. The system automatically determines if the user gets approved or pending status based on available quota.

#### Request
**Content-Type:** `application/json`

```json
{
  "event_id": 1,
  "wallet_address": "0x1234567890abcdef..."
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| event_id | integer | Yes | ID of the event to register for |
| wallet_address | string | Yes | User's wallet address |

#### Response
**Success (201 Created) - Approved**
```json
{
  "message": "Whitelist successful! You are registered as an event participant.",
  "data": {
    "id": 1,
    "event_id": 1,
    "user_id": 5,
    "wallet_address": "0x1234567890abcdef...",
    "status": "approved",
    "created_at": "2024-06-10T10:00:00Z",
    "updated_at": "2024-06-10T10:00:00Z"
  }
}
```

**Success (201 Created) - Pending (Quota Full)**
```json
{
  "message": "Whitelist successful, but quota is full. You are on the waiting list.",
  "data": {
    "id": 2,
    "event_id": 1,
    "user_id": 6,
    "wallet_address": "0xabcdef1234567890...",
    "status": "pending",
    "created_at": "2024-06-10T10:00:00Z",
    "updated_at": "2024-06-10T10:00:00Z"
  }
}
```

**Error Responses**

**400 Bad Request - Invalid Request**
```json
{
  "error": "Invalid request body"
}
```

**400 Bad Request - Missing Required Fields**
```json
{
  "error": "event_id and wallet_address are required"
}
```

**404 Not Found - User Not Found**
```json
{
  "error": "Wallet address is not registered as a user"
}
```

**404 Not Found - Event Not Found**
```json
{
  "error": "Event not found"
}
```

**409 Conflict - Already Registered**
```json
{
  "error": "You have already registered for this event"
}
```

**500 Internal Server Error**
```json
{
  "error": "Failed to find user"
}
```

---

#### 6. Cancel Whitelist Entry
**POST** `/api/users/whitelist/cancel`

Cancels a user's whitelist registration for an event.

#### Request
**Content-Type:** `application/json`

```json
{
  "event_id": 1,
  "wallet_address": "0x1234567890abcdef..."
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| event_id | integer | Yes | ID of the event to cancel registration for |
| wallet_address | string | Yes | User's wallet address |

#### Response
**Success (200 OK)**
```json
{
  "message": "Whitelist cancelled successfully"
}
```

**Error Responses**

**400 Bad Request - Invalid Payload**
```json
{
  "error": "invalid request payload"
}
```

**400 Bad Request - Missing Required Fields**
```json
{
  "error": "event_id and wallet_address are required"
}
```

**404 Not Found - User Not Found**
```json
{
  "error": "user not found"
}
```

**404 Not Found - Whitelist Entry Not Found**
```json
{
  "error": "whitelist entry not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to cancel whitelist"
}
```

---

### Vendors

#### 7. Register Vendor
**POST** `/api/vendors/register`

Registers a new vendor account.

#### Request
**Content-Type:** `application/json`

```json
{
  "vendor_name": "NFT Events Co.",
  "email": "vendor@example.com",
  "contact_info": "+62-812-3456-7890",
  "wallet_address": "0xabcdef1234567890..."
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| vendor_name | string | Yes | Vendor's business name |
| email | string | Yes | Vendor's email address |
| contact_info | string | No | Vendor's contact information |
| wallet_address | string | Yes | Vendor's wallet address (must be unique) |

#### Response
**Success (201 Created)**
```json
{
  "id": 1,
  "vendor_name": "NFT Events Co.",
  "email": "vendor@example.com",
  "contact_info": "+62-812-3456-7890",
  "wallet_address": "0xabcdef1234567890...",
  "created_at": "2024-06-10T10:00:00Z",
  "updated_at": "2024-06-10T10:00:00Z"
}
```

**Error Responses**

**400 Bad Request - Invalid Request**
```json
{
  "error": "invalid request"
}
```

**400 Bad Request - Missing Required Fields**
```json
{
  "error": "vendor_name, email, and wallet_address are required"
}
```

**409 Conflict - Wallet Already Registered**
```json
{
  "error": "wallet address already registered with another account"
}
```

**500 Internal Server Error**
```json
{
  "error": "server error"
}
```

---

#### 8. Get Events by Vendor
**GET** `/api/vendors/:walletAddress/events`

Retrieves all events created by a specific vendor, based on their wallet address.

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| walletAddress | string | Yes | The vendor's wallet address. |

#### Response
**Success (200 OK)**
Returns an array of `Event` objects.
```json
[
  {
    "id": 1,
    "title": "NFT Conference 2024",
    "description": "Annual NFT conference",
    "vendor_id": 1,
    "start_date": "2024-06-15T09:00:00Z",
    "end_date": "2024-06-15T17:00:00Z",
    "status": "upcoming",
    "created_at": "2024-06-10T10:00:00Z",
    "updated_at": "2024-06-10T10:00:00Z",
    "picture": "uploads/1718000000_event.jpg",
    "maxattendees": 100,
    "location": "Jakarta Convention Center",
    "attendees": 25,
    "requirements": {"items": ["Laptop", "Notebook"]},
    "agenda": {"sessions": [{"time": "09:00", "topic": "Introduction"}]}
  }
]
```

**Error Responses**

**400 Bad Request - Missing Wallet Address**
```json
{
  "error": "walletAddress is required"
}
```

**404 Not Found**
```json
{
  "error": "vendor not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "database error message"
}
```

---

### Events

#### 9. Get All Events
**GET** `/api/events/all`

Retrieves all events with organizer and whitelisted count.

#### Response
**Success (200 OK)**
```json
[
  {
    "id": 1,
    "title": "NFT Conference 2024",
    "description": "Annual NFT conference",
    "vendor_id": 1,
    "start_date": "2024-06-15T09:00:00Z",
    "end_date": "2024-06-15T17:00:00Z",
    "status": "upcoming",
    "created_at": "2024-06-10T10:00:00Z",
    "updated_at": "2024-06-10T10:00:00Z",
    "picture": "uploads/1718000000_event.jpg",
    "maxattendees": 100,
    "location": "Jakarta Convention Center",
    "organizer": "NFT Events Co.",
    "whitelisted": 25
  }
]
```

**Error (500 Internal Server Error)**
```json
{
  "error": "database connection error"
}
```

---

#### 10. Create Event
**POST** `/api/events/create`

Creates a new event. Requires multipart form data.

#### Request
**Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Event title |
| description | string | No | Event description |
| vendor_id | integer | Yes | ID of the vendor creating the event |
| start_date | string | Yes | Start date in RFC3339 format (e.g., "2024-06-15T09:00:00Z") |
| end_date | string | Yes | End date in RFC3339 format (e.g., "2024-06-15T17:00:00Z") |
| status | string | Yes | Event status (e.g., "upcoming", "ongoing", "completed") |
| maxattendees | integer | Yes | Maximum number of attendees |
| location | string | No | Event location |
| requirements | string | No | JSON string containing event requirements |
| agenda | string | No | JSON string containing event agenda |
| picture | file | Yes | Event image file |

#### Response
**Success (201 Created)**
```json
{
  "id": 12,
  "title": "NFT Conference 2024",
  "description": "Annual NFT conference",
  "vendor_id": 1,
  "start_date": "2024-06-15T09:00:00Z",
  "end_date": "2024-06-15T17:00:00Z",
  "status": "upcoming",
  "created_at": "2024-06-10T10:00:00Z",
  "updated_at": "2024-06-10T10:00:00Z",
  "picture": "uploads/1718000000_event.jpg",
  "maxattendees": 100,
  "location": "Jakarta Convention Center",
  "attendees": 0,
  "requirements": {"items": ["Laptop", "Notebook"]},
  "agenda": {"sessions": [{"time": "09:00", "topic": "Introduction"}]}
}
```

**Error Responses**

**400 Bad Request - Missing Required Fields**
```json
{
  "error": "title is required"
}
```

**400 Bad Request - Invalid Data Types**
```json
{
  "error": "vendor_id must be a number"
}
```

**400 Bad Request - Invalid Date Format**
```json
{
  "error": "invalid start_date format"
}
```

**400 Bad Request - Invalid Date Logic**
```json
{
  "error": "start_date must be before end_date"
}
```

**400 Bad Request - Invalid Max Attendees**
```json
{
  "error": "maxattendees must be a positive number"
}
```

**400 Bad Request - Missing Picture**
```json
{
  "error": "picture file is required"
}
```

**400 Bad Request - Invalid JSON**
```json
{
  "error": "invalid requirements JSON"
}
```

**500 Internal Server Error - File Upload Issues**
```json
{
  "error": "failed to open uploaded file"
}
```

**500 Internal Server Error - Database Issues**
```json
{
  "error": "DB scan failed"
}
```

---

#### 11. Get Event Detail
**GET** `/api/events/:id`

Retrieves detailed information about a specific event, including organizer and attendee counts.

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | The event ID |

#### Response
**Success (200 OK)**
```json
{
  "id": 1,
  "organizer": "NFT Events Co.",
  "title": "NFT Conference 2024",
  "description": "Annual NFT conference",
  "vendor_id": 1,
  "start_date": "2024-06-15T09:00:00Z",
  "end_date": "2024-06-15T17:00:00Z",
  "status": "upcoming",
  "created_at": "2024-06-10T10:00:00Z",
  "updated_at": "2024-06-10T10:00:00Z",
  "picture": "uploads/1718000000_event.jpg",
  "maxattendees": 100,
  "location": "Jakarta Convention Center",
  "attendees": 25,
  "whitelisted": 50,
  "requirements": {"items": ["Laptop", "Notebook"]},
  "agenda": {"sessions": [{"time": "09:00", "topic": "Introduction"}]}
}
```

**Error Responses**

**400 Bad Request - Missing Event ID**
```json
{
  "error": "event ID is required"
}
```

**404 Not Found**
```json
{
  "error": "event not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "database error message"
}
```

---

## Data Models

### Event
```json
{
  "id": "integer",
  "title": "string",
  "description": "string",
  "vendor_id": "integer",
  "start_date": "datetime (RFC3339)",
  "end_date": "datetime (RFC3339)",
  "status": "string",
  "created_at": "datetime (RFC3339)",
  "updated_at": "datetime (RFC3339)",
  "picture": "string (file path)",
  "maxattendees": "integer",
  "location": "string",
  "attendees": "integer",
  "requirements": "json (optional)",
  "agenda": "json (optional)"
}
```

### User
```json
{
  "id": "integer",
  "email": "string",
  "wallet_address": "string",
  "name": "string",
  "created_at": "datetime (RFC3339)",
  "updated_at": "datetime (RFC3339)"
}
```

### Vendor
```json
{
  "id": "integer",
  "vendor_name": "string",
  "email": "string",
  "contact_info": "string",
  "wallet_address": "string",
  "created_at": "datetime (RFC3339)",
  "updated_at": "datetime (RFC3339)"
}
```

### CertificateWithEvent
```json
{
  "id": "integer",
  "event_id": "integer",
  "user_id": "integer",
  "certificate_data": "string",
  "mint_status": "string",
  "mint_transaction_hash": "string",
  "created_at": "datetime (RFC3339)",
  "updated_at": "datetime (RFC3339)",
  "event_title": "string",
  "event_description": "string",
  "event_start_date": "datetime (RFC3339)",
  "event_location": "string",
  "event_picture": "string (full URL)"
}
```

### Whitelist Entry
```json
{
  "id": "integer",
  "event_id": "integer",
  "user_id": "integer",
  "wallet_address": "string",
  "status": "string (approved|pending)",
  "created_at": "datetime (RFC3339)",
  "updated_at": "datetime (RFC3339)"
}
```

### EventWithUserStatus
This model is returned by `GET /api/users/:walletAddress/events` and includes the user's status for that event.
```json
{
  "id": "integer",
  "title": "string",
  "description": "string",
  "vendor_id": "integer",
  "start_date": "datetime (RFC3339)",
  "end_date": "datetime (RFC3339)",
  "status": "string",
  "created_at": "datetime (RFC3339)",
  "updated_at": "datetime (RFC3339)",
  "picture": "string (file path)",
  "maxattendees": "integer",
  "location": "string",
  "attendees": "integer",
  "user_status": "string (present|whitelisted)",
  "requirements": "json (optional)",
  "agenda": "json (optional)"
}
```

---

## Error Handling

All endpoints return appropriate HTTP status codes:

- **200 OK**: Successful GET requests
- **201 Created**: Successful POST requests that create resources
- **400 Bad Request**: Invalid input data
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists (e.g., duplicate wallet address, already registered for event)
- **500 Internal Server Error**: Server-side errors

Error responses follow this format:
```json
{
  "error": "error message description"
}
```

---

## CORS Configuration

The API supports CORS with the following configuration:
- **Allow Origins**: All origins (`*`)
- **Allow Methods**: GET, POST, OPTIONS
- **Allow Headers**: All headers

---

## File Upload

Event images are stored in the `uploads/` directory with the following naming convention:
```
uploads/{timestamp}_{original_filename}
```

Example: `uploads/1718000000_event_banner.jpg`

---

## Whitelist System

The whitelist system automatically manages event registration:

1. **Quota Check**: When a user registers for an event, the system checks if there are available spots
2. **Status Assignment**: 
   - `approved`: User gets a confirmed spot if quota is available
   - `pending`: User is placed on waiting list if quota is full
3. **Duplicate Prevention**: Users cannot register for the same event multiple times
4. **User Validation**: Only registered users can create or cancel whitelist entries

---

## Notes

1. **Wallet Address Uniqueness**: Wallet addresses must be unique across both users and vendors tables.
2. **Date Format**: All dates should be in RFC3339 format (ISO 8601).
3. **File Upload**: Only image files are supported for event pictures.
4. **Database**: The API uses PostgreSQL as the database backend.
5. **Environment**: Configuration is loaded from `.env` file.
6. **JSON Fields**: Requirements and agenda fields accept JSON strings that are validated before storage.
7. **Whitelist Logic**: The system automatically determines whitelist status based on event capacity. 