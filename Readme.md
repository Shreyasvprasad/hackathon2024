CloudKeep: a personal cloud storage application written in go, also using Scylladb and MinIo

Setup Instructions:
download go or goland  
download docker  
run a scylladb image in docker - docker pull scylladb/scylla, docker run --name scylla -d scylladb/scylla
run a minio image in the docker - docker pull docker://minio/minio , setup a local environment variable and then run it 
'docker run -d -p 9000:9000 -e "MINIO_ACCESS_KEY=minioadmin" -e "MINIO_SECRET_KEY=minioadmin" --name minio minio/minio server /data'

open GCP and get the credentials for google oAuth (passwordless login)  
save the client secret credentials and minio credentials in environment variables 
connect the go file to the scylladb, minio instance and monitor and troubleshoot the connection using the logd.  

Assumptions:  
The frontend of this app is going to be developed by someone else, we only have to focus on backend development.  
The google sign-in comes before all the functionality of the website.  
We cannot link the tables (use foreign keys) in the database together as it is a NoSQL-type DB, so we use efficient querying techniques like indexes.  

Database schema:  
![image](https://github.com/user-attachments/assets/07b9bd05-e7e4-4c9e-a866-ca6b6e6dc8c0)  

We have created three tables in the ScyllaDB instance:  
- **Users** - Stores user information from Google OAuth.  
- **Files** - Stores user notes with support for real-time synchronization.  
- **Notes** - Stores metadata for uploaded files.  

### Schema Details  

**Files Table:**  
- `file_id` - A unique identifier for each file. Ensures that every file has a distinct and immutable ID.  
- `file_size`  
- `file_type` - The type of the file (e.g., image/jpeg, application/pdf). Helps identify file formats.  
- `file_url` - The URL where the file is stored by MinIO.  
- `filename` - The name of the file uploaded by the user.  
- `last_accessed` - The last time the file was accessed or downloaded. It monitors activity.  

**Users Table:**  
- `user_id` - A unique identifier for each user. Acts as the primary key to identify users in the system.  
- `created_at` - The date and time when the user account was created.  
- `email` - The user’s email address.  
- `full_name` - The full name of the user.  
- `google_id` - The unique ID provided by Google OAuth. Used for authentication and linking the user’s Google account.  
- `profile_picture` - The URL of the user’s profile picture fetched from Google OAuth.  
- `updated_at` - TIMESTAMP: The last time the user’s profile information was updated.  

**Notes Table:**  
- `note_id` - A unique identifier for each note.  
- `content` - The main content of the note. Stores the text or information entered by the user.  
- `created_at` - The date and time when the note was created.  
- `title` - The title of the note.  
- `updated_at` - The last time the note was modified. Helps track changes and synchronize updates across devices.  
- `user_id` - The ID of the user who created the note.  

Additional Features:  
Indexing: In ScyllaDB, I have implemented two indexes that will allow for efficient querying of the database.  

Working Demo link:  https://youtu.be/9gSqgzjlc5o
