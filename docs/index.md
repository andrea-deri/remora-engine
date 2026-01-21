---
title: What is REMORA Engine?
nav_order: 1
---

**REMORA (Reactive Engine for Mocked Runtime Actions)** is a high-speed, extensible rules engine built in Go. Its primary purpose is to simulate and replicate complex behavioral patterns (defined as Minnow profiles) in response to specific incoming requests.  
The engine is designed to mock dynamic workflows and interactions in a controlled environment. It allows developers to define predictable yet flexible outcomes for testing, prototyping, or runtime simulation.  
The architecture is modular and consisting of two main layers:
* **Harbor**: The Engine access layer, responsible for request routing and workload coordination.
* **Tidal**: The computation layer, responsible for resource matching and behavior evaluation.
* **Benthos**: The data abstraction layer, responsible for persistence and retrieval of Minnow profiles.

The system aims to be deployed in both **monolithic** and **distributed modes**.  

In **monolithic mode**, the entire system is orchestrated in a main entry point, which initializes the configuration, attach to the persistence via Registry component, starts the core engine services by Tidal component and then launches the Harbor component to begin handling requests. This mode is particularly suitable for local testing and simple contexts where a high number of concurrent requests is not expected.    

In **distributed mode**, it is possible to deploy as separate contexts:
 - the proxying layer via the _Harbor_ component, which can provide balancing policies between instances of Tidal component according to correlation criteria.
 - the computation layer via the _Tidal_ component, which can be composed of multiple parallel instances, each of which uses a specific group of Minnows for correlation.
 - the registry layer via the _Benthos_ component, which is the abstraction layer that interconnects the Tidal component with the actual persistence of the Minnow profiles.

Currently, distributed management is purely theoretical, and the monolithic mode is the only one that can actually be implemented. In the future, distributed management will be developed to make REMORA Engine more flexible and suitable for different circumstances.