var db = connect("mongodb://root:rootpassword@127.0.0.1:27017/admin");

db = db.getSiblingDB('stock');

db.createUser(
    {
        user: "stock_user",
        pwd: "stock_pass",
        roles: [{ role: "readWrite", db: "stock" }],
        passwordDigestor: "server",
    }
);

db.clients.createIndex({ "ident": 1 }, { unique: true });