Create users in DB:
test> use admin
switched to db admin
admin> db.getUsers()
{ users: [], ok: 1 }
admin> db.createUser({"user":"user", "pwd":"user", "roles": []})
{ ok: 1 }
admin> db.createUser({"user":"admin", "pwd":"admin", "roles": []})
{ ok: 1 }
admin> db.getUsers()
{
users: [
{
_id: 'admin.admin',
userId: UUID('d74ec7d0-6e9e-4869-bafb-3a865fa88b5d'),
user: 'admin',
db: 'admin',
roles: [],
mechanisms: [ 'SCRAM-SHA-1', 'SCRAM-SHA-256' ]
},
{
_id: 'admin.user',
userId: UUID('2510fbf3-71a4-4397-8f6d-3e80098c947f'),
user: 'user',
db: 'admin',
roles: [],
mechanisms: [ 'SCRAM-SHA-1', 'SCRAM-SHA-256' ]
}
],
ok: 1
}
admin> use registry
switched to db registry
registry> db.dropAllUsers()
{ n: 2, ok: 1 }
registry> db.getUsers()
{ users: [], ok: 1 }
registry> use admin
switched to db admin
admin> db.getUsers()
{
users: [
{
_id: 'admin.admin',
userId: UUID('d74ec7d0-6e9e-4869-bafb-3a865fa88b5d'),
user: 'admin',
db: 'admin',
roles: [],
mechanisms: [ 'SCRAM-SHA-1', 'SCRAM-SHA-256' ]
},
{
_id: 'admin.user',
userId: UUID('2510fbf3-71a4-4397-8f6d-3e80098c947f'),
user: 'user',
db: 'admin',
roles: [],
mechanisms: [ 'SCRAM-SHA-1', 'SCRAM-SHA-256' ]
}
],
ok: 1
}
admin> exit

