---
title: Engine Overview
nav_order: 2
has_children: true
---

In this system, the **Minnow** is the core concept: it represents a single profiled resource with defined behaviors. Each Minnow is identified by a unique **Resource ID** and corresponds to a specific resource, defined by its **path**, **HTTP method**, and optional **special headers**. Each Minnow has a **profile**, which can include multiple **behaviors**. Each behavior can be triggered when certain conditions are met in the request data.   
Additionally, a single Minnow can represent a group of similar RESTful resources differentiating them by path parameters evaluation, allowing the same profile to handle different resources with different behaviors.  
The REMORA Engine operates through multiple layers: **Harbor** as the entry point, **Tidal** as the processing core, and **Benthos** as the persistence abstraction. In essence, when a request arrives Harbor forwards it to Tidal, which identifies the appropriate Minnow, validates the request against the defined conditions, executes the corresponding behavior and returns the resulting effect.


### 2.1. Harbor: the proxy <a name="overview-harbor"></a>
**Harbor** is the public-facing interface of the REMORA Engine. It is composed of two primary components that manage the lifecycle of an incoming HTTP request:

-  **Proxy**:  
    This component is the actual HTTP server. It starts an HTTP listener on a configured port and base path and its responsibility is to receive raw HTTP requests, extract their data (URL, method, headers, body), and package them into a internal transfer message. Then, this transfer message is immediately sent to the Harbor Dispatcher's channel for asynchronous processing. It also manages other settings and processes such as request timeouts.
-  **Dispatcher**:  
    This component is the coordination hub: it receives the transfer message from the proxy via an asynchronous channel, then orchestrates the communication with Tidal engine. It first calls the **Search Engine** API to find a matching resource and if a resource is found, it calls the **Validation Engine** API to evaluate the request against the resource's defined behaviors. Then it return a transfer message on the response channel to the Proxy component based on the validation result. To achieve this, it controls concurrency with a semaphore that limits the number of requests processed at the same time.


### 2.2. Tidal: the engine <a name="overview-tidal"></a>
**Tidal** is the core "brain" of REMORA Engine. It is not a single component but a collection of modules orchestrated by the Harbor Dispatcher. Its process is split into two main phases: **search** and **validation**.

-  **Search phase**:  
    The **Search Engine** is responsible for matching an incoming request to an indexed Minnow. It uses the reference indexes defined in the Reference Registry to perform an efficient lookup based on the request's path, method, and headers. If a match is found, the Search Engine returns the Resource ID and any extracted path parameters to the dispatcher.

-  **Validation phase**:  
    Once the Search Engine identifies a resource, the Dispatcher passes the request and the Resource ID to the **Validation Engine**. This micro-engine retrieves the full Minnow profile for the indexed resource from the Profile Registry, then executes the behavior validation (compiled before in a **Coral Syntax Tree**) against the request's full context (headers, query params, body) to find a matching condition. The result of this validation is a well-defined **behavior**, generate as a transfer message (either containing status code, body and other or an error) which the Dispatcher uses to build the final HTTP response.


### 2.3. Benthos: the persistence registry <a name="overview-registry"></a>
**Benthos** is the persistence and data access layer of the entire REMORA Engine, providing an abstraction layer for retrieving Minnow definitions.
At startup, a factory initializes a bundle of concrete registry implementations based on the configuration properties. Depending on the user’s configuration, the persistence layer can rely on different backends: **MongoDB**, **etcd** or a local/remote **filesystem**. This component exposes two main interfaces, both consumed by the Tidal Engine:

 - **Reference Registry**:  
    Provides access to Minnow references through an efficient in-memory **search tree**, allowing rapid matching of incoming requests to a Minnow Resource ID without loading the full behavior logic at runtime.

 - **Profile Registry**:  
    Manages the dictionary of indexed Minnow profiles, whose behaviors are compiled into a **Coral Syntax Tree** and used during the validation process once a matching resource is found.


### 2.4. Coral: the Behavior Syntax Tree <a name="overview-coral"></a>
**Coral** is the internal language and parser that give form to the Minnow behavior through a list of **conditions**, turning a simple rule string into an executable Abstract Syntax Tree (AST). So, Coral's purpose is to parse a text-based condition (like `field1 eq_s "value" and (field2 gt_n 10 or field3 is_empty)`) into a structured, tree-like object that the Validation Engine can easily evaluate.  
Coral language defines the structure required for generate the Syntax Tree:  
 - Each **Minnow Profile** contains a list of **Behaviors**
 - Each Behavior contains a **Condition** and an **Effect** related to the successful validation of the condition
 - Each Condition contains one  **Clause** and/or more nested Conditions, joined each other with logical connectors
 - Each Clause contains a **Field**, an **Operator** (unary or binary) and a **Value** to check against during validation

Coral follows a multi-stage compilation process to transform a raw rule string into an executable structure:
 - **Tokenization**: The **Tokenizer** first scans the condition string and breaks it down into a list of basic **Tokens**
 - **Parsing**: The **Parser** receives these tokens, it understands the language grammar by handling logical AND/OR precedence and parenthetical groups in order to build the hierarchical AST.
 - **Compilation**: The generated syntax tree is then processed by a **Compiler**. This step performs optimizations, such as pre-compiling any regular expressions defined in regex clauses, to make the final evaluation more efficient.
 - **Execution**: Finally, the **Runner** (invoked by the Validation Engine) takes the compiled Profile and executes it against the input data from the request, traversing the Condition tree to find a valid Effect.

Coral is the _engine-within-the-engine_ that gives REMORA its dynamic, behavior-based power, allowing it to validate complex conditions beyond simple path matching.
