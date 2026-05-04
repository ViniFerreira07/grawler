# 🚀 Grawler: High-Performance Concurrent Web Scraper

**Grawler** is a robust CLI tool built in Go designed to demonstrate the power of **CSP (Communicating Sequential Processes)**. It processes batches of URLs concurrently, extracting metadata and validating security standards with high efficiency and low memory footprint.

---

## 🎯 Technical Highlights

This project was developed to explore the core strengths of the Go language, focusing on:

*   **Advanced Concurrency:** Utilizing `errgroup` and `sync.WaitGroup` for managed goroutine lifecycles.
*   **The Go Way:** Idiomatic error handling and explicit memory management using pointers.
*   **Decoupled Architecture:** Heavy use of **Interfaces** to allow easy extension of "Checker" logic (Open/Closed Principle).
*   **Resource Safety:** Strict use of `defer` for closing network connections and file descriptors to prevent memory leaks.
*   **I/O Optimization:** Efficient JSON marshalling and buffered CSV reading.

## 🏗️ Architecture

The system follows a **Producer-Consumer pipeline**:
1.  **Scanner:** Reads URLs from a CSV source.
2.  **Worker Pool:** A set of concurrent workers (Checkers) processes the URLs.
3.  **Reporter:** A centralized consumer gathers results via **Channels** and exports them to a JSON report.



## ⚡ Performance

| Metric | Result |
| :--- | :--- |
| **Concurrency Model** | Goroutines + Channels |
| **Throughput** | ~100 URLs / 2.5 seconds* |
| **Memory Usage** | < 20MB RAM |
| **Timeout Handling** | Strict 10s per request |

*\*Results may vary based on network latency and bandwidth.*

## 🛠️ Installation & Usage

### Prerequisites
*   Go 1.22+
*   Make (optional)

### Setup
```bash
# Clone the repository
git clone [https://github.com/your-user/grawler.git](https://github.com/your-user/grawler.git)
cd grawler

# Build the project using the Makefile
make build