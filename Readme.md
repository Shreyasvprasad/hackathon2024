Hack_to_Intern Submission
CloudKeep: a personal cloud storage application written in go,also using Scylladb and MinIo

Setup Instructions: 
download go or goland
download docker
run a scylladb image in docker
run a minio image in the docker
open GCP and get the credentials for oAuth (paswordless login)
save the client secret credentials and minio credentials in environment variables
connect the go file to the scylladb, minio instance.

Assumptions:
The frontend of this app is going to be developed by someone else, we only have to focus on backend development.
The google sign in comes before all the  functionality of the website.
We cannot link the tables(use foreign keys) in the database together as it is a nosql type db, so we use effecient querying techniques like indexes.

Database schema:
![image](https://github.com/user-attachments/assets/07b9bd05-e7e4-4c9e-a866-ca6b6e6dc8c0)
we have created three tables in the scylladb instance
Users - Stores user information from Google OAuth.
Files - Stores user notes with support for real-time synchronization.
Notes -Stores metadata for uploaded files.
Files- file_id-A unique identifier for each file. Ensures that every file has a distinct and immutable ID.
       file_size
       file_type	The type of the file (e.g., image/jpeg, application/pdf). Helps identify file formats.
       file_url		The URL where the file is stored by MinIo.
       filename		The name of the file uploaded by the user.
       last_accessed		The last time the file was accessed or downloaded.It monitors activity.
Users- user_id	A unique identifier for each user. Acts as the primary key to identify users in the system.
      created_at		The date and time when the user account was created.
      email		The user’s email address. 
      full_name		The full name of the user. 
      google_id	The unique ID provided by Google OAuth. Used for authentication and linking the user’s Google account.
      profile_picture		The URL of the user’s profile picture fetched from Google OAuth. 
      updated_at	TIMESTAMP	The last time the user’s profile information was updated.
Notes- note_id		A unique identifier for each note. 
       content		The main content of the note. Stores the text or information entered by the user.
      created_at		The date and time when the note was created. 
      title		The title of the note. 
      updated_at	The last time the note was modified. Helps track changes and synchronize updates across devices.
       user_id	The ID of the user who created the note. 

Working Demo link:
