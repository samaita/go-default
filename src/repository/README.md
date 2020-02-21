Repository will store any Database handler. Querying, or Creating/ Inserting into any database will stored here. This layer will act for CRUD to database only. No business process happen here. Only plain function to Database.

This layer also have responsibility to choose what DB will used in Application. Could be Mysql, MongoDB, MariaDB, Postgresql whatever, will decided here.

If using ORM, this layer will control the input, and give it directly to ORM services.

If calling microservices, will handled here. Create HTTP Request to other services, and sanitize the data. This layer, must fully act as a repository. Handle all data input - output no specific logic happen.

This Repository layer will depends to Connected DB , or other microservices if exists.