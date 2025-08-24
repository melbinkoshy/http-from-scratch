# HTTP from Scratch (Go)

A work-in-progress project to implement a minimal HTTP server **starting from raw TCP sockets in Go** — without using Go's built-in `net/http` package.  
The goal is to understand how HTTP works under the hood by building it step by step.

---

## 🚀 Motivation
We often use Go’s `net/http` for building servers, but HTTP is really just **text over TCP**.  
This project explores:
- How to accept TCP connections  
- How to read and parse HTTP requests manually  
- How to construct valid HTTP responses  
- Eventually, how to support basic routing

---

