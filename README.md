# CrawlX 🕷️ – A Lightweight Go Web Crawler

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

## ✨ Overview
CrawlX is a lightweight CLI-based web crawler built in Go. It can recursively crawl web pages, follow links, and provide insights into visited URLs.  
The project is structured for **cross-platform usage**, with installation scripts for both Linux/macOS and Windows.  

---

## 🚀 Features
- 🌐 **Recursive Crawling** — explore web pages up to a user-defined depth.  
- ⚡ **Concurrency with Goroutines** — crawl multiple links in parallel for faster performance.  
- 📝 **Customizable Flags** — configure URL, depth, verbosity.  
- 📦 **Cross-Platform Installation** — works on Linux/macOS and Windows.  
- 🔗 **URL resolution** – converts relative links into absolute URLs  
- 📋 **Simple CLI** – easy to run with flags  
- 🛡️ **Error handling** – skips broken links gracefully

---

## 📦 Installation

### Linux / macOS
```bash
git clone https://github.com/sh4dowkey/Crawlx.git
cd Crawlx
chmod +x setup.sh
sudo ./setup.sh
```

### Windows (PowerShell)
```powershell
git clone https://github.com/sh4dowkey/Crawlx.git
cd Crawlx
Set-ExecutionPolicy Bypass -Scope Process -Force
.\setup.ps1
```

That’s it! 🎉  
Once installed, you can run `crawlx` from anywhere in your terminal.  

---

## 🛠️ Usage

Basic usage:
```bash
crawlx --url https://example.com --depth 4
```

Shorthand flags:
```bash
crawlx -u https://example.com -d 3
```

Verbose mode:
```bash
crawlx -u https://example.com -d 2 --verbose
```

Example output:
```
[+] Crawling: https://toscrape.com (Depth 0)
    [200 OK] Found 10 links.
    [+] Crawling: http://quotes.toscrape.com/ (Depth 1)
        [200 OK] Found 49 links.
    [+] Crawling: http://quotes.toscrape.com/login (Depth 1)
        [200 OK] Found 4 links.
    [+] Crawling: http://quotes.toscrape.com/js (Depth 1)
        [200 OK] Found 5 links.
    [+] Crawling: http://quotes.toscrape.com/js-delayed (Depth 1)
        [200 OK] Found 5 links.
    [+] Crawling: http://books.toscrape.com (Depth 1)
        [200 OK] Found 73 links.
    [+] Crawling: http://quotes.toscrape.com (Depth 1)
        [200 OK] Found 49 links.
    [+] Crawling: http://quotes.toscrape.com/random (Depth 1)
        [200 OK] Found 6 links.
    [+] Crawling: http://quotes.toscrape.com/scroll (Depth 1)
        [200 OK] Found 4 links.
    [+] Crawling: http://quotes.toscrape.com/search.aspx (Depth 1)
        [200 OK] Found 4 links.
    [+] Crawling: http://quotes.toscrape.com/tableful (Depth 1)
        [200 OK] Found 35 links.

--- Crawl Complete ---

--- Crawl Summary ---
- Crawled 11 pages in 7.444s.
- Success: 11 pages
- Warnings: 0 pages
- Errors: 0 pages

--- Visited Links ---
- Visited: http://quotes.toscrape.com/random
- Visited: http://quotes.toscrape.com/search.aspx
- Visited: https://toscrape.com
- Visited: http://quotes.toscrape.com/
- Visited: http://quotes.toscrape.com/login
- Visited: http://quotes.toscrape.com/js
- Visited: http://quotes.toscrape.com/js-delayed
- Visited: http://quotes.toscrape.com
- Visited: http://quotes.toscrape.com/scroll
- Visited: http://quotes.toscrape.com/tableful
- Visited: http://books.toscrape.com
```

---


## ⚙️ Usage

```bash
Usage:
  crawlx [options]

Options:
  -u, --url string     The starting URL to crawl (required)
  -d, --depth int      Depth level for recursive crawling (default: 2)
  --verbose bool       Enable verbose output with more details
```

---

## 📂 Project Structure
```
Crawlx/
│── cmd/            # Main source code (main.go, crawl.go)
│── dist/           # Built binaries
│── setup.sh        # Linux/macOS installer
│── setup.ps1       # Windows installer
│── go.mod          # Go module file
│── go.sum          # Go dependencies
```

---

## 🛣 Roadmap
- [ ] Add support for robots.txt parsing  
- [ ] Export results (JSON/CSV)  
- [ ] Ignore filetypes (e.g., images, PDFs)  
- [ ] Crawl domain restrictions  

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


<p align="center">⭐ If you like this project, give it a star on GitHub! ⭐</p>

<p align="center">🕸️ Crawl smarter. Crawl faster. CrawlX.</p>
