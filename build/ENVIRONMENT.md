BUILD_DEBUG   enable debug level for logging
BUILD_TRACING   enable tracing level for logging
BUILD_NAME   set the name of the instance of the service
BUILD_TITLE   set the title in the web page
BUILD_PORT_NB   the local port of the web service (default=8080)
BUILD_API_KEYS   space separated list of valid API keys
BUILD_SESSION_KEY
BUILD_POLICY   OPA policy for access control
BUILD_OPASVC   OPA service port to get the OPA policy for access control
BUILD_REALM   Basic authentication realm
BUILD_STAFF_USER   username of the administrator
BUILD_STAFF_PASSWORD   password of the administrator
BUILD_PARTICIPANT_USER   username of the user
BUILD_PARTICIPANT_PASSWORD   password of the user
BUILD_CERT_PEM   certificate for TLS (HTTPS) communication
BUILD_KEY_PEM   key for TLS (HTTPS) communication
BUILD_LOG_FILE   filename of the logging file or if it is "-" log all messages to the console
BUILD_LOKI_SERVER   URL of the Loki Server, e.g. the Grafana Cloud
BUILD_LOKI_USER   user name as defined for the data source as basic authentication
BUILD_LOKI_PASSWORD   password as defined for the data source as basic authentication
BUILD_LOKI_KEY   key/token as defined for the data source
BUILD_LABELS   set of labels for Grafana Loki, like "app:myapp,tenant:mycustomer"
BUILD_BUFFERSIZE   number of log lines cached before sending to Grafana Loki
BUILD_MAX_DELAY   max time in seconds after which logs are flushed from buffer
BUILD_LANGUAGE
BUILD_LANGUAGES
BUILD_USE_SSE   enable support for _server side event_ communication (default=false)
BUILD_PROGRESS_DURATION   default duration of the progress bar (default=100ms)
BUILD_RAPIDOC_DOC   enable Rapidoc for the OpenAPI viewer (default=false)
BUILD_ELEMENTS_DOC   enable Elements for the OpenAPI viewer (default=false)
