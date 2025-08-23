# CrawlX 🕷️  
*A fast, recursive, and concurrent web crawler written in Go.*

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
- 📊 **Summary Output** — keep track of visited links.  

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
crawlx --url https://example.com --depth 2
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
Crawling: https://example.com
Found: https://example.com/about
Found: https://example.com/contact
Visited: 3 links
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

## 🛣️ Roadmap
- [ ] Add support for robots.txt parsing  
- [ ] Export results (JSON/CSV)  
- [ ] Ignore filetypes (e.g., images, PDFs)  
- [ ] Crawl domain restrictions  

---

## 🤝 Contributing
Contributions are welcome!  
1. Fork the repo  
2. Create a new branch (`feature-xyz`)  
3. Commit your changes  
4. Submit a pull request  

---

## 📜 License
This project is licensed under the **MIT License**.  

---

### 🕸️ Crawl smarter. Crawl faster. CrawlX.
