# Features

The generator dredger generates from an _openapi.spec_ file the basic code for a microservice in Go. The code includes not only the HTTP server and its handler. It also includes typically required, sophisticated functionally for an API, like security and observability functions.

The generated code can be enhanced by the business code as required. Files, which can be adapted by the developer will not be overwritten, if the generator runs again.

## Basics

The code is basically structured using the clean architecture principles. The core includes code for configuration, logging, tracing, the command line and the _echo_ HTTP server. Further, for the given schematas in the OpenAPI specification the code for the entities and the validation of the data of the entities will be generated. The entities includes the field tags for marshalling and unmarshalling data to YAML, JSON and XML format.

### Handler

For each API endpoint in the OpenAPI specification a basic stub handler code will be generated. The name of the handler is based on the _operationId_ of the API endpoint. The handler code includes general code for the observation and logging. It can be adapted by adding the required business logic.

If an API endpoint contains a tag _builtin_ in your _tags_, then no generic handler code will be generated. The tag _builtin_ should especially be used for default API endpoints for the lifecycle functions, the default stylesheets, the default javascript files and the _/index.html_ or _/content.html_ files.

If an API endpoint contains a tag _page_ in your _tags_, then a [templ](https://templ.guide/) template will also be created. A templ template allows to write HTML pages mixed with go code and generate a go function, which can be used easily in your handlers. Further, a localizer and a language selector (_languages.templ_) is setup to translate strings using a [i18n library](github.com/nicksnyder/go-i18n/v2/i18n) for internationalization.

### Server Side Events

If in the OpenAPI specification for the API endpoints the path "/events" with the builtin operations _get_ and _post_ are given, additional code for Server Side Events (_SSE_) (_rest/handleEvents.go_) will be generated. Especially, for tasks, which will need longer the functions _ProgressPico_ and _ProgressBootstrap_ (_rest/progress.go_) can be used to send using server side events a _progress bar_ code to HTMX for Pico and Bootstrap CSS, e.g.

	f := func() {
		_, err := http.Get("http://localhost:9090/slowz")
		if err != nil {
			log.Warn().Err(err).Msg("Slow call failed")
		}
	}
	ProgressPico(f)

The _progress bar_ itself need to be declared in the frontend, e.g.

    <script src="js/sse.js"></script>

    <div hx-ext="sse" sse-connect="/events?stream=progress" sse-swap="Progress"></div>

and will be visible, when a call start, progress over time and reappear at the end.

### Configuration

The generated service can be configured using the default values, a _.env_ file, environment variables and the command line options (highest priority).

### Lifecycle

Standard handler _livez_ and _readyz_ for the lifecycle of the service will be generated. Further, an handler _infoz_ iss added to get the meta information about the service at runtime.

To restrict the web crawlers a handler for _robots.txt_ is generated and adaptable.

### Doc

The comprehensive documentation of a service should be contained in the OpenAPI specification. Therefor, the OpenAPI documentation will also be embedded with the service and online available at runtime (_/doc_). Two viewers (_rapidoc_ and _Stoplight elements_) are available and configurable.

### Testing

When offering an API it should be proper tested. Therefore, the generator optionally provides code for BDD based testing using [godog](https://github.com/cucumber/godog) and _specification by example_ to describe tests.

## Security

To ensure security and data privacy also the basic security functions will be added to the code. This allows rather easily to implement a _policy enforcement point_.

### Validating data

One aspect of a _policy enforcement point_ is to validate the input of a request by validating the data, e.g. by checking minimum and maximum values or matching pattern. The validation code will be generated from the given restrictions in the OpenAPI specification.

### Security policies

Further, a valid access must be proved by checking the authentication of users and their roles or an API key as well as checking the authorization using rules, which could be provided using the [Open Policy Agent](https://www.openpolicyagent.org/) technology with its access rules.

## Observability

When a distributed system doesn't work corrrectly, you have to find the origin of the failure. This can be a tedious task in a distributed system. To find a problem you need observe what is happen at which time and in which sequence and what are the relevant information and details to identify the problem. Therefore it is essential to instrument the code, so it can be observered easily. This requires logging of messages in a useful way as well as providing tracing data. Both have to be collected in a central place, and should be monitored in dashboards for the application support and others, e.g. using the [Grafana](https://grafana.com/) toolbox.

### Logging

The generated code supports logging and tracing. The logging is based on the zerolog interface. The log messages supports the [12 factor app](https://12factor.net/) principles to log messages for the _system administration_ and also direct logging to the _application support_ using log files and log forwarders. Also, logs could be directly send to Grafana Loki using the Loki API. Optionally, debug logging can be enabled also on the command line.

### Tracing and metrics

The stub handlers are instrumented using the [Open Telemetry](https://pkg.go.dev/go.opentelemetry.io/otel) libraries and tools for tracing and providing metrics for key indicators.

## Deployment

As micro services are mostly deployt as OCI conainers a standard Dockerfile and image manifest are provided to create efficient, small images for AMD64 and ARM64 based systems.

## Development

To support the development an example file for _.gitlab-ci.yml_ will be provided. Further, a [Justfile](https://just.systems/) is available, which includes typical task required by a developer.
