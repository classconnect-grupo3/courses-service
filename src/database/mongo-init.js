print("MongoDB initialization script starting...");

db = db.getSiblingDB('admin');
db.createUser({
    user: 'admin',
    pwd: 'password',
    roles: [
        { role: 'userAdminAnyDatabase', db: 'admin' },
        { role: 'readWriteAnyDatabase', db: 'admin' },
        { role: 'dbAdminAnyDatabase', db: 'admin' }
    ]
});

db = db.getSiblingDB('courses_db');
db.createCollection('courses');

print("MongoDB initialization completed successfully!"); 