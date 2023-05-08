## Overview

This project was built for learning purposes. 

The API lets you to manage collections of words with translations provided by [google.translate.com.](google.translate.com)   
For authorization/authentication opendID connect is used. Connection with [google.translate.com](google.translate.com) established through HTTP 2.0. 

It's only the backend of the whole application. The application itself can be found at: [https://github.com/Kin-dza-dzaa/flash_cards](https://github.com/Kin-dza-dzaa/flash_cards) 

TODO list for this project:

*   [x] Clean architecture
*   [x] Docker/docker-compose
*   [ ] K8S
*   [x] OpenTelemetry 
    *    [x]  Tracing (Jaeger)
    *    [ ]  Merics (Prometheus)
*   [x] Structured logging
*   [x] Swagger docs
*   [x] Google coding guidelines/ 12 factor app
*   [x] Migrations (go-migrate)

## Usage

**Run app:**

```plaintext
make run
```

**Run tests:**

**In order for tests to work you will need an available docker API on port 2375 with disabled tls.**

```plaintext
make test
```

**Run tests with cover:**

```plaintext
make cover
```