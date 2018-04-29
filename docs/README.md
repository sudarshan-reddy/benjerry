## Architectural justification

For simplicity and demonstration it serves to currently have four APIs only,
one for each of CRUD.

Internally, we use the chi library because it is faster than the native
go router. Chi is only slightly behind httprouter. This is okay because chi 
lets you use context and that equalises the very minuscule speed advantage.

We use docker-compose to automate the whole thing because it runs the server 
and the postgres app together with migration.

Authentication is abstracted into implementable interfaces that would allow
an easy plug in of any type of Authentication. I am currently using a static
Bearer token to authenticate. If needed, we can implement a jwt microservice
and authenticate through that by implementing middleware.Authenticator.

The same goes for database storage. models.models.go is an abstracted 
storage interface that any database can implement. A postgres implementation
is provided.

The most important decision is `Update`. It is done as a single route 
with only the body changing. Update will be done only for the values that 
change.
