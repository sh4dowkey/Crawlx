# CrawlX 🕷️

<div align="center">

**A fast, concurrent web crawler built in Go**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/sh4dowkey/crawlx?style=flat-square)](https://github.com/sh4dowkey/crawlx/releases/latest)
[![Build Status](https://img.shields.io/github/actions/workflow/status/sh4dowkey/crawlx/build.yml?style=flat-square)](https://github.com/sh4dowkey/crawlx/actions)

*Professional web crawling with concurrent processing and intelligent error handling*

[🚀 Quick Start](#installation) • [📖 Documentation](#usage) • [🤝 Contributing](#contributing) • [📋 Releases](https://github.com/sh4dowkey/crawlx/releases)

</div>

---

## 🌟 Overview

CrawlX is a professional-grade web crawler designed for speed, reliability, and ease of use. Built with Go's powerful concurrency features, it can crawl websites **10x faster** than traditional sequential crawlers while maintaining respectful server interaction through built-in delays and proper error handling.

Whether you're doing SEO analysis, testing website deployments, or conducting security research, CrawlX provides clean, actionable output with comprehensive error categorization.

## ✨ Features

- **⚡ Concurrent Architecture**: 10 workers processing URLs simultaneously with intelligent task distribution
- **🔍 Smart URL Validation**: Pre-flight checks handle malformed URLs, IP addresses, and domain validation  
- **📊 Professional Reports**: Real-time progress with categorized error reporting (4xx, 5xx, network failures)
- **🎯 Depth Control**: Configurable recursive crawling with same-domain restriction
- **🛡️ Respectful Crawling**: Built-in delays, retry logic, and proper User-Agent identification
- **🎨 Dual Output Modes**: Clean standard output or detailed verbose mode with link discovery
- **🌐 Cross-Platform**: Native binaries for Windows, Linux, and macOS

## 📦 Installation

### Quick Install (Recommended)

1. **Download** the latest release for your operating system:
   - [Linux (amd64)](https://github.com/sh4dowkey/crawlx/releases/latest) • [Windows](https://github.com/sh4dowkey/crawlx/releases/latest) • [macOS](https://github.com/sh4dowkey/crawlx/releases/latest)

2. **Extract and install**:
   ```bash
   # Linux/macOS
   tar -xzf crawlx-*.tar.gz
   cd crawlx-*
   sudo ./scripts/setup.sh
   
   # Windows (Run as Administrator)
   # Extract ZIP file
   .\scripts\setup.ps1
   ```

3. **Verify installation**:
   ```bash
   crawlx -u https://example.com -d 1
   ```

### Build from Source

For developers who want to modify the code or contribute:

```bash
# Prerequisites: Go 1.22+, Git, Make
git clone https://github.com/sh4dowkey/crawlx.git
cd crawlx

# Build for current OS
make build

# Build for all platforms  
make build-all

# Run directly (development)
make run
```

## 🚀 Usage

### Basic Examples

```bash
# Quick website scan
crawlx -u https://example.com

# Deep crawl with detailed output
crawlx -u https://example.com -d 3 --verbose

# Local development server
crawlx -u http://localhost:8080 --allow-ip

# Quick link validation
crawlx -u https://mywebsite.com -d 1
```

### Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--url` | `-u` | Target URL to crawl **(required)** | - |
| `--depth` | `-d` | Maximum crawling depth | `2` |
| `--verbose` | `-v` | Enable detailed progress output | `false` |
| `--allow-ip` | `-i` | Allow crawling IP addresses | `false` |

### Output Examples

**Standard Mode:**
```
Crawling https://example.com (depth: 2)

 [200] https://example.com (245ms)
 [200] https://example.com/about (189ms)
 [404] https://example.com/missing (89ms)
 [FAIL] https://timeout.com (Network: connection timeout)

============================================================
                    CRAWL REPORT
============================================================

Scan completed in 2.3s using 10 workers
Total: 25 | Success: 20 | Errors: 5

Successful Pages (20):
  https://example.com (245ms)
  https://example.com/about (189ms)
  ...

Client Errors - 4xx (3):
  https://example.com/missing (404 Not Found)
  https://example.com/private (403 Forbidden)

Network Errors (2):
  https://timeout.com (Connection failed: timeout)
  
External Links (5):
  https://github.com
  https://stackoverflow.com
  ...

============================================================
```

**Verbose Mode:**
```
  [+] Crawling: https://example.com (Depth 0)
    ↳ [200 OK] Found 15 links (245ms)
    Links found:
      - https://example.com/about
      - https://example.com/contact
      - https://example.com/products
    [+] Crawling: https://example.com/about (Depth 1)
      ↳ [200 OK] Found 8 links (189ms)
      Links found:
        - https://example.com/team
        - https://example.com/history
        ... and 6 more
```

## 🏗️ How It Works

CrawlX employs a **worker pool architecture** for maximum efficiency:

1. **URL Validation**: Comprehensive pre-flight checks validate protocols, domains, and handle IP addresses securely
2. **Concurrent Processing**: 10 workers process URLs simultaneously while a dispatcher manages the crawling queue
3. **Intelligent Discovery**: Each page's HTML is parsed to extract and resolve relative links to absolute URLs
4. **Smart Categorization**: Results are automatically categorized into success (2xx), client errors (4xx), server errors (5xx), and network failures
5. **Respectful Behavior**: Built-in 100ms delays between requests prevent server overload, with automatic retry logic for reliability

The crawler maintains a visited URL map to prevent infinite loops and duplicate processing, while external links are discovered and reported separately.

## 🎯 Use Cases

**SEO & Website Analysis**
- Site structure mapping and navigation flow analysis
- Broken link detection for improved user experience  
- Page performance monitoring and response time analysis
- External link discovery for competitor research

**Development & Testing**  
- Local development server validation (`localhost` support)
- Post-deployment link verification and health checks
- Large website content inventory and auditing
- Performance regression testing

**Security Research**
- Website reconnaissance and scope mapping (with proper authorization)
- Link validation and discovery for security assessments
- Infrastructure analysis through response pattern examination

**Educational & Learning**
- Understanding web crawling concepts and HTTP protocols
- Studying concurrent programming patterns in Go
- Learning about professional tool development and architecture

## 🔧 Advanced Features

### Error Handling
CrawlX provides comprehensive error categorization to help identify issues quickly:
- **Network Errors**: DNS failures, connection timeouts, connection refused
- **4xx Client Errors**: Not found (404), forbidden (403), unauthorized (401)  
- **5xx Server Errors**: Internal server error (500), bad gateway (502), service unavailable (503)
- **Parsing Errors**: Malformed HTML, invalid link formats

### Performance Characteristics
- **Speed**: 10-20x faster than sequential crawlers for most websites
- **Memory Efficient**: Buffered channels prevent excessive memory usage
- **Respectful**: 100ms delays and 3-attempt retry limits prevent server overload
- **Scalable**: Handles small personal sites to large enterprise websites

### URL Management
- **Same-domain restriction**: Prevents scope creep during crawls
- **Relative link resolution**: Converts relative URLs to absolute for proper processing
- **Duplicate prevention**: Efficient tracking prevents redundant crawling
- **External link discovery**: Identifies and reports off-site references

## 🗺️ Project Structure

```
crawlx/
├── cmd/crawlx/          # Main application entry point
│   └── main.go          # CLI interface and validation
├── internal/            # Private application code  
│   ├── crawl/          # Core crawling engine
│   ├── parse/          # HTML parsing and link extraction
│   └── util/           # Utility functions and helpers
├── scripts/            # Installation scripts
├── Makefile           # Build automation
└── go.mod             # Go module definition
```

## 🔮 Roadmap

**Version 2.0** - Advanced Crawling
- robots.txt support for respectful crawling
- JSON/CSV export capabilities  
- YAML configuration file support
- Sitemap.xml integration

**Version 2.1** - Intelligence & Analysis
- Content-type filtering (skip images, PDFs)
- Page title and meta description extraction
- Real-time progress indicators
- Advanced performance metrics

**Version 2.2** - Professional Features  
- Resume interrupted crawls
- Custom headers and authentication
- URL include/exclude patterns
- Per-domain rate limiting

## 🤝 Contributing

We welcome contributions from developers of all skill levels! Here's how to get started:

### Quick Start
1. Fork the repository and clone your fork
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes and add tests if applicable
4. Commit with a clear message: `git commit -m "Add feature: description"`
5. Push and create a Pull Request

### Development Setup
```bash
git clone https://github.com/yourusername/crawlx.git
cd crawlx
go mod download
make run  # Test your changes
```

### Ways to Contribute
- **🐛 Bug Reports**: Help identify issues with detailed reproduction steps
- **💡 Feature Requests**: Suggest improvements with clear use cases
- **📝 Documentation**: Improve README, code comments, or add examples
- **🧪 Testing**: Add test cases or improve existing test coverage
- **💻 Code**: Implement new features or fix existing bugs

Check our [Issues](https://github.com/sh4dowkey/crawlx/issues) page for tasks labeled `good first issue` to get started.

## 📜 License

This project is licensed under the MIT License, allowing free use in both personal and commercial projects. See the [LICENSE](LICENSE) file for complete details.

## ⭐ Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/sh4dowkey/crawlx/issues)
- **GitHub Discussions**: [Ask questions or share ideas](https://github.com/sh4dowkey/crawlx/discussions)  
- **Documentation**: Comprehensive guides and examples in this README
- **Source Code**: Well-commented code for understanding implementation

## 🙏 Acknowledgments

Built with appreciation for:
- **Go Programming Language** for excellent concurrency support and HTTP libraries
- **golang.org/x/net/html** for robust HTML parsing capabilities
- **Open Source Community** for inspiration, feedback, and contributions

---

<div align="center">

**If you find CrawlX useful, please consider giving it a ⭐ on GitHub!**

Made with ❤️ by [sh4dowkey](https://github.com/sh4dowkey) and [contributors](https://github.com/sh4dowkey/crawlx/graphs/contributors)

</div>