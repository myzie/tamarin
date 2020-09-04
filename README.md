# Tamarin

Cloud native scripting language and shell.

## Brainstorming

 * Developer ergonomics in the shell is priority
 * Like a cloud-aware, modernized bash
 * Document CRUD in S3 supported natively
 * Memory persists via DynamoDB or similar
 * Selectively load variables from environment
 * Built-in HTTP
 * Built-in secret retrieval from Secrets Manager
 * Named context for persisted memory, optionally shared by other shells
 * Message passing between shells via SNS
 * Batteries loaded defaults for working in AWS
 * Inspect variables, documents, DDB memory
 * Metadata definition for shell specifies how secrets, variables
   are to be loaded. And other settings like the context, bucket name, etc.
 * Option for aliases set up via metadata definition
 * Immutability, case statements
 * Dependency injection by name like pytest

## Load Common Definitions

For a given use case, the developer may choose to load common configurations
of imports, secrets, and built-ins. For example, a common "AWS developer"
configuration could be loaded which takes care of associated setup. These
configurations should document or publish what they provide.

## Developer Visibility

The shell should provide mechanisms to inspect state, execution history,
and available configuration and functions.

 * What common configuration is active?
 * What variables are defined and what are their values?
 * Last N network calls?
 * Which secrets are loaded?

## Developer Conveniences

The shell should feel lightweight, support concise commands, and be flexible
in terms of customization.

 * Command aliases
 * Great JSON support
 * Syntax highlighting
 * Easy to recall previous commands, even across shells, or across systems
 * Side-by-side views, or similar, allowing execution and monitoring of
   variables, state, and events
 * Browsable object hierarchy
 * Store state, recall elsewhere
 * Share shell state with another developer in realtime

Implicit `root` object contains top-level variables.

## Leverage Cloud Platforms

AWS Lambda and similar services have unique properties that have not been
fully leveraged in development environments. Consider how a shell could provide
the developer access to these services to solve common problems, while offering
significant parallel computing capabilties. Potentially interactions with your
local shell could be evaluated in a remote environment, for example.

## Language Features

Intentionally not many to start, but the language is here to support the
developer. Over time it will grow to be more expressive.

 * Should feel comparable to Python, Javascript, and Bash
 * For-in style iteration
 * Let and const for variables
 * First-class functions
 * Dynamically typed, strongly typed, like Python

## S3 Object Access

```
let doc = s3.get_json("path/to/foo.json")
print("retrieved document", keys=doc.keys())
```

Consider keeping imports external to the script source code. In this example,
the developer would specify `import cloud.aws.s3 as s3` externally in their
environment or the shell configuration.

## HTTP Requests

```
let content = http.get("https://api.ipify.org?format=json")
```

## Generate types

