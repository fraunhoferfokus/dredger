SIMPLE_API_DEBUG   enable debug level for logging
SIMPLE_API_TRACING   enable tracing level for logging
SIMPLE_API_NAME   set the name of the instance of the service
SIMPLE_API_TITLE   set the title in the web page
SIMPLE_API_PORT_NB   the local port of the web service (default=8080)
SIMPLE_API_API_KEYS   space separated list of valid API keys
SIMPLE_API_SESSION_KEY
SIMPLE_API_POLICY   OPA policy for access control
SIMPLE_API_OPASVC   OPA service port to get the OPA policy for access control
SIMPLE_API_REALM   Basic authentication realm
SIMPLE_API_STAFF_USER   username of the administrator
SIMPLE_API_STAFF_PASSWORD   password of the administrator
SIMPLE_API_PARTICIPANT_USER   username of the user
SIMPLE_API_PARTICIPANT_PASSWORD   password of the user
SIMPLE_API_CERT_PEM   certificate for TLS (HTTPS) communication
SIMPLE_API_KEY_PEM   key for TLS (HTTPS) communication
SIMPLE_API_LOG_FILE   filename of the logging file or if it is "-" log all messages to the console
SIMPLE_API_LOKI_SERVER   URL of the Loki Server, e.g. the Grafana Cloud
SIMPLE_API_LOKI_USER   user name as defined for the data source as basic authentication
SIMPLE_API_LOKI_PASSWORD   password as defined for the data source as basic authentication
SIMPLE_API_LOKI_KEY   key/token as defined for the data source
SIMPLE_API_LABELS   set of labels for Grafana Loki, like "app:myapp,tenant:mycustomer"
SIMPLE_API_BUFFERSIZE   number of log lines cached before sending to Grafana Loki
SIMPLE_API_MAX_DELAY   max time in seconds after which logs are flushed from buffer
SIMPLE_API_LANGUAGE
SIMPLE_API_LANGUAGES
SIMPLE_API_USE_SSE   enable support for _server side event_ communication (default=false)
SIMPLE_API_PROGRESS_DURATION   default duration of the progress bar (default=100ms)
SIMPLE_API_RAPIDOC_DOC   enable Rapidoc for the OpenAPI viewer (default=false)
SIMPLE_API_ELEMENTS_DOC   enable Elements for the OpenAPI viewer (default=false)
