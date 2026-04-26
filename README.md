# aws-iac-pulumi

A hands-on learning repository for building common AWS architectures using
Pulumi (Go) as the Infrastructure as Code tool, with Datadog integrated for
observability at every stage.

The architectures progress deliberately from simple to complex — each one
introduces new AWS services and patterns while reusing and building on
everything that came before.

---

## Purpose

This repository exists to develop practical, hands-on experience across three
skill areas in parallel:

- **AWS** — understanding core services and how they compose into real architectures
- **Pulumi with Go** — writing infrastructure as code using reusable, well-structured components
- **Datadog** — integrating observability from the ground up, not as an afterthought

Each architecture is intentionally kept self-contained with its own documentation,
so the learning journey for any given topic is clear and focused.

---

## Repository Structure

```
aws-iac-pulumi/
│
├── architectures/               # Individual architecture projects
│   ├── 01-static-website/       # Static website hosting using S3 and CloudFront with OAC
│   └── ...                      # New architectures will be added progressiveley
│
├── helper-modules/              # Reusable utilities with no direct AWS resources
│   └── tags/                    # Standard Pulumi tagging helper
│   └── ...                      # Grows based on requirement
│
├── infra-modules/               # Reusable AWS ComponentResource building blocks
│   ├── cloudfront/              # CloudFront distribution component
│   ├── s3/                      # S3 bucket component
│   └── ...                      # Grows as new architectures are added
│
├── go.work                      # Go workspace (see note below for the reason on why this is added to the source control)
├── go.work.sum
├── LICENSE
└── README.md
```

### Module Design

**`helper-modules`** contains pure Go utilities that support infrastructure
code — things like tagging helpers, naming conventions, and config utilities.
These have no direct dependency on AWS resources.

**`infra-modules`** contains reusable Pulumi ComponentResources — each one
wraps one or more AWS resources into a well-defined, configurable building
block. Architecture projects import from here rather than defining raw AWS
resources inline.

This separation keeps architecture code focused on wiring components together
rather than low-level resource configuration.

---

## Architecture Index

Each architecture folder contains its own README with:

- What the architecture does and a diagram
- Which AWS services are involved and why
- What you will learn from building it
- Prerequisites and deployment instructions
- What to look at in Datadog after deploying

---

## Prerequisites

Before deploying any architecture you will need:

- [Pulumi CLI](https://www.pulumi.com/docs/install/) installed
- [Go 1.26.2+](https://go.dev/dl/) installed
- AWS account with appropriate permissions
- Datadog account with API and App keys

---

## Pulumi Backend

This repository was developed using the local filesystem as the Pulumi
backend. It is fully compatible with Pulumi Cloud — no code changes are
required.

Architectures that depend on the shared VPC stack use a configurable stack
reference string. The correct value to set for your backend is documented
in the README of the particular architecture.

---

## Getting Started

Each architecture has its own deployment instructions, prerequisites, and
configuration steps documented in its README. Start there rather than
following a generic flow, as dependencies and setup vary significantly
between architectures.

If this is your first time with the repository, start with:
  architectures/01-static-website/README.md

---

## A Note on `go.work`

Go's general guidelines recommend against committing `go.work` for libraries
and tools that others will consume as dependencies, as it can interfere with
their local workspace setup.

This repository is a self-contained monorepo with no external module consumers.
The `go.work` file is the intended mechanism for linking `helper-modules`,
`infra-modules`, and `architectures` as local modules without publishing them
to a package registry. Committing it ensures the repository works out of the
box on any machine without additional setup steps.

---

## Conventions

- All infrastructure is managed exclusively through Pulumi — no manual changes in the AWS console
- Every AWS resource is tagged with `Creator`, `Team`, `Service`, `Env` and `ManagedBy` via the tags helper-module
- Secrets are managed via `pulumi config set --secret` — no credentials are stored in plaintext
- Commit messages follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
