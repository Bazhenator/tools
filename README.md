# tools
tools is a supplimentary service with grpc and other toolkit for cleaning-service. 

**Accessable packages:**
- `middleware` - contains logs producer with tracing id for grpc server and auth options (in future).
- `server` - contains ***http*** and ***grpc*** server's options for listening and shutdown servers.
- `logger` - contains implementation of project ***logger***, using `zap`.
