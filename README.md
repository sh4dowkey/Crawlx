# 🌐 CrawlX – A Lightweight Go Web Crawler

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go" />
  <img src="https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=flat-square" />
  <img src="https://img.shields.io/github/stars/yourusername/crawlx?style=social" />
  <img src="https://img.shields.io/github/forks/yourusername/crawlx?style=social" />
</p>

CrawlX is a **fast, recursive web crawler built in Go**.  
It can fetch links up to a specified depth, resolve relative URLs into absolute ones, and supports **concurrent crawling with goroutines & channels** 🚀  

---

## ✨ Features
- 🌍 **Recursive crawling** – follows links up to a configurable depth  
- ⚡ **Concurrent fetching** – crawl multiple links in parallel using goroutines  
- 🔗 **URL resolution** – converts relative links into absolute URLs  
- 📋 **Simple CLI** – easy to run with flags  
- 🛡️ **Error handling** – skips broken links gracefully  

---

## 📦 Installation

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/crawlx.git
cd crawlx
```

### 2. Build the binary
```bash
go build -o crawlx ./cmd
```

### 3. Run the crawler
```bash
./crawlx --url https://example.com --depth 2
```

Or using short flags:
```bash
./crawlx -u https://example.com -d 2
```

---

## ⚙️ Usage

```bash
Usage:
  crawlx [options]

Options:
  -u, --url string     The starting URL to crawl (required)
  -d, --depth int      Depth level for recursive crawling (default: 2)
```

---

## 📂 Project Structure
```
crawlx/
│── cmd/
│   ├── main.go      # Entry point (parses flags, starts crawl)
│   ├── crawl.go     # Core crawler logic
│── go.mod           # Go module file
│── README.md        # Documentation
│── dist/            # Compiled binaries (ignored in git)
```

---

## 🖼️ Example Output

```bash
$ ./crawlx -u https://golang.org -d 1
[+] Crawled: https://golang.org [200]
[+] Crawled: https://golang.org/doc/ [200]
[+] Crawled: https://golang.org/pkg/ [200]
```

---

## 🚀 Concurrency Model

CrawlX leverages Go’s **goroutines + sync.WaitGroup** to fetch links in parallel:

```go
wg.Add(1)
go Crawl(link, depth-1, &wg)
```

This ensures faster crawling without blocking on a single request.

---

## 🛠️ Roadmap

- [x] Basic recursive crawling  
- [x] Depth control  
- [x] Absolute URL resolution  
- [x] Concurrency with goroutines  
- [ ] Add `--verbose` flag for detailed logs  
- [ ] Robots.txt handling  
- [ ] Export crawled URLs to JSON/CSV  
- [ ] Add colored terminal output  
- [ ] Configurable concurrency limit  

---

## 🤝 Contributing

Contributions are welcome!  
1. Fork the repo  
2. Create your feature branch (`git checkout -b feature-name`)  
3. Commit changes (`git commit -m "Add feature"`)  
4. Push to branch (`git push origin feature-name`)  
5. Open a Pull Request 🚀  

---

## 📜 License

This project is licensed under the **MIT License** – see the [LICENSE](LICENSE) file for details.

---

## 💡 Inspiration
This project was created as part of a **learning roadmap** to master Go, concurrency, and system-level programming concepts in the context of a simple but powerful tool: a web crawler.

---

<p align="center">⭐ If you like this project, give it a star on GitHub! ⭐</p>
