# API Endpoints

## User Endpoints

1. GET /api/user
   Description: Get current user's metadata.
   Auth: Bearer token (any authenticated user)
   Response 200 JSON: User object
   User fields:

- id (int64)
- username (string)
- phone (int64)
- email (string)
- password (string, hashed)

## Log Head Endpoints

2. GET /api/log-heads
   Description: List all log heads.
   Auth: Bearer token
   Response 200 JSON: Array of LogHead objects
   LogHead fields:
   - id (uint)
   - subject (string)
   - start_date (string, RFC3339 timestamp)
   - end_date (string, RFC3339 timestamp)
   - writer_id_list (array of int64)
   - owner_id (uint)
   - owner_name (string)
   - log_contents (array of LogContent, may be null/empty if not preloaded)
     LogContent fields (within log_contents):
   - id (uint)
   - log_head_id (uint)
   - writer_id (uint)
   - content (string)
   - date (string, RFC3339 timestamp)
   - writer_name (string)

## Admin Endpoints

3. PUT/PATCH /api/admin/log-heads/:id
   Description: Update log head metadata. All fields are optional (partial update).
   Auth: Bearer token (admin role required)
   Request Body JSON (all fields optional):

- subject (string)
- start_date (time, ISO format)
- end_date (time, ISO format)
- writer_id_list (array of int64)
- owner_id (uint)
  Response 200 JSON: Updated LogHead object
  LogHead fields:
- id (uint)
- subject (string)
- start_date (time)
- end_date (time)
- writer_id_list (array of int64)
- owner_id (uint)
- log_contents (array of LogContent, auto-populated)

4. GET /api/members/search?q=...
   Description: Search for members by username, email, or phone number.
   Auth: Bearer token (admin role required)
   Query Parameters:

- q (string, required): Search query
  Response 200 JSON: Array of Account objects
  Account fields:
- id (int64)
- username (string)
- phone (int64)
- email (string)
- password (string, hashed)

## Log Content Endpoints

5. PUT/PATCH /api/log-contents/:id
   Description: Update log content. All fields are optional (partial update).
   Auth: Bearer token (any authenticated user)
   Permission: Admin or the original writer of the log content
   Request Body JSON (all fields optional):

- content (string)
- date (time, ISO format)
  Response 200 JSON: Updated LogContent object
  LogContent fields:
- id (uint)
- log_head_id (uint)
- writer_id (uint)
- content (string)
- date (time)

6. DELETE /api/log-contents/:id
   Description: Delete log content.
   Auth: Bearer token (any authenticated user)
   Permission: Admin or the original writer of the log content
   Response 204: No content
