Use the compass app from mongoDB to create a new database called registry, with a collection called documents.

Create a user with readWrite access to the registry database.
```mongosh
use registry
db.createUser(
   {
     user: "user",
     pwd: "user",
     roles: [ "readWrite" ]
   }
)
```

Create a superuser with readWrite and dbAdmin access to the registry database.
```mongosh
use registry
db.createUser(
   {
     user: "superuser",
     pwd: passwordPrompt(),  // Or  "<cleartext password>"
     roles: [ "readWrite", "dbAdmin" ]
   }
)
```
