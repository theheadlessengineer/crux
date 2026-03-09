# Crux

Standards built in. Not bolted on.

---

## What is Crux

Crux is an internal CLI tool for generating production-ready microservice
skeletons. It is the starting point for every service built in this
organisation.

A single command generates a fully structured, runnable service with
company standards, security configuration, resilience patterns,
observability wiring, infrastructure as code, and CI/CD pipelines already
in place. Teams write business logic. Crux handles everything else.

Crux is built on a plugin architecture. Every integration beyond the
core — databases, caches, message brokers, cloud providers, AI tools,
observability backends — is a self-contained, versioned plugin. The core
enforces the standard. Plugins extend it.

```
crux new payment-service
```

---

## Vision

A world where every engineer in the organisation ships production-grade
services with confidence — where the distance between an idea and a
running, secure, observable, compliant service is measured in minutes,
not weeks.

## Mission

Crux gives every engineering team a single, extensible starting point —
embedding company standards, security, resilience, and observability
directly into the foundation of every service, so teams can focus entirely
on the problems only they can solve.