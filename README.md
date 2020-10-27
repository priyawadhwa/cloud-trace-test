# cloud-trace-test

This is a POC that Trace ID can be propagated across multiple tools.
The Trace ID is saved in a file by the "parent tool", and is read from the file by the secondary tool.
This works as long as the "parent tool" is constantly running:


To test, replace any calls to `WithProjectID` with your own GCP project.

Then in one terminal, run:

```
go run service/service.go
```

to start the initial trace.

In a second terminal, run:

```
go run main.go
```

which should propagate the trace ID and run as a subspan.
