WebApp framework

Just another go framework to make handling api calls easy.

This framework is highly inspired by the Echo framework, but is different as its more configurable.

- You can bring your own router, comes with a slightly modified version of httprouter included.
- You can bring your own binder, comes with a binder that binds based on struct tags
- You can bring your own validator, comes with playground validator
- You can bring your own JSON encoder, comes with the std encoding
- Uses slog as the logging interface in the context
