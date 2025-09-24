## Installing
### Linux
* You must have [Go v1.24.3](https://go.dev/doc/install) (or higher), [PostgreSQL](https://www.postgresql.org/download/), [goose](https://github.com/pressly/goose/releases) and [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html) installed on your system
* Clone this repo to your desired location
* Run ```go build``` on your local machine in cloned repository
### Windows
* You can repeat the linux steps, it will work just fine
* Or alternatively you can [download an executable from release page](https://github.com/lackingworth/Go-Chirpy/releases)

## Info
 Chirpy is a simple web-server emulating Twitter (X). 
 It makes use of goose and sqlc (postgres) to handle database migrations and SQL queries.

## Available API endpoints
* ```GET /``` - Displays simple placeholder html from ```http.FileServer```
* ```GET /api/chirps``` - Returns all the chirps (posts)
> [!NOTE]  
> 
> There are two possible queries:
> ```/api/chirps?author_id=authorID``` and ```/api/chirps?sort=desc```. If no authorID is specified then the method returns all chirps. If no sort value is specified then the method returns the chirps in default ascending order
> JSON Response Body
> ```json
>{
>   "id": "uuid",
>   "created_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "updated_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "body": "some text",
>   "user_id": "uuid of a user who created the chirp"
>}
> ```
* ```POST /api/chirps``` - Creates a chirp (post) while validating and censoring it
> [!NOTE]  
> 
> JSON Request Body
> ```json
> {
>   "body": "some text",
> }
> ```
> JSON Response Body
> ```json
> 
>{
>   "id": "uuid of a chirp",
>   "created_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "updated_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "body": "validated and censored text",
>   "user_id": "uuid of a user who created the chirp"
>}
> ```
* ```GET /api/chirps/{chirpID}``` - Returns the chirp (post) by specified ID
> [!NOTE]  
> 
> JSON Response Body
> ```json
>{
>   "id": "uuid of a chirp",
>   "created_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "updated_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "body": "some text",
>   "user_id": "uuid of a user who created the chirp"
>}
> ```
* ```DELETE /api/chirps/{chirpID}``` - Deletes the chirp (post) by specified ID
* ```GET /api/healthz``` - Displays the health of a web server
* ```POST /api/login``` - Authenticates user
> [!NOTE]  
> 
> JSON Request Body
> ```json
> {
>   "email": "test@example.com",
>   "password": "password_example"
> }
> ```
> JSON Response Body
> ```json
> {
>   "user":
>     {
>          "id": "uuid",
>          "created_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>          "updated_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>          "email": "test@example.com",
>          "is_chirpy_red": false
>     },
>   "token": "JWT token",
>   "refresh_token": "refresh token"
> }
> ```
*  ```POST /api/refresh``` - Refreshes token
> [!NOTE]  
> 
> JSON Request Body
> ```json
> {
>   "token": "JWT token"
> }
> ```
> JSON Response Body
> ```json
> {
>   "token": "refresh token",
> }
> ```
*  ```POST /api/revoke``` - Revokes refresh token
*  ```POST /api/users``` - Creates user
> [!NOTE]  
> 
> JSON Request Body
> ```json
> {
>   "email": "test@example.com",
>   "password": "password_example"
> }
> ```
> JSON Response Body
> ```json
> 
>{
>   "id": "uuid",
>   "created_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "updated_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "email": "test@example.com",
>   "is_chirpy_red": false
>}
> ```
* ```PUT /api/users``` - Updates specified user
> [!NOTE]  
> 
> JSON Request Body
> ```json
> {
>   "email": "test@example.com",
>   "password": "password_example"
> }
> ```
> JSON Response Body
> ```json
> 
>{
>   "id": "uuid",
>   "created_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "updated_at": "2025-07-14 15:50:13.793654 +0000 UTC",
>   "email": "test@example.com",
>   "is_chirpy_red": false
>}
> ```
* ```GET /admin/metrics``` - Displays the number of successful hits to the home page
* ```POST /admin/reset``` - Resets the number of successful hits to the home page and the database tables
> [!NOTE]  
> 
> Configure ```.env``` file correctly
> To successfuly hit ```/admin/``` endpoints - your ```PLATFORM``` environmental variable must equal to ```dev```
* ```POST /api/polka/webhooks``` - Webhook handler
> [!NOTE]  
> 
> JSON Request Body
> ```json
> {
>   "event": "user.upgraded",
>   "data":
>     {
>       "user_id" "uuid of a user"
>     }
> }
> ```
* New features to come

## Importan environmental variables

* ```DB_URL``` - URL to your database in ```postgres://username:password@domain:port/chirpy?sslmode=disable``` format
* ```PLATFORM``` - Monitors access to certain endpoints. Should equal to ```dev``` for full access
* ```JWT_SECRET``` - Your JWT secret. To generate one you can run ```openssl rand -base64 64``` command in tour terminal
* ```POLKA_API_KEY``` - API key for webhook handlers

## Version History

* v.0.0.1:

    * Initial Release
