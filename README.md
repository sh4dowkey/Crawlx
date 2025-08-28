# CrawlX 🕷️ – A Modern Go Web Crawler

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go" />
  <img src="https://img.shields.io/badge/Status-In%20Development-orange?style=flat-square" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=flat-square" />
  <img src="https://img.shields.io/github/stars/sh4dowkey/Crawlx?style=social" />
  <img src="https://img.shields.io/github/forks/sh4dowkey/Crawlx?style=social" />
</p>

<p align="center">
  <i>A fast, concurrent, and recursive web crawler built in Go. Designed to be a lightweight but powerful tool for exploring websites from the command line.</i>
</p>

---


## ✨ Features

- 🌐 **Recursive Crawling**: Explore web pages up to a user-defined depth.
- ⚡ **Concurrency (In Progress)**: Designed to crawl multiple links in parallel for maximum speed.
- 🎨 **Colored CLI Output**: User-friendly and readable terminal output.
- 🔧 **Customizable Flags**: Configure the starting URL, crawl depth, and verbosity.
- 🔗 **URL Resolution**: Correctly converts relative links into absolute, crawlable URLs.
- 🛡️ **Graceful Error Handling**: Skips broken links and handles HTTP errors without crashing.

---

## 🚀 Quick Start

Get up and running in under 60 seconds.

1.  **Download** the latest release for your OS from the **[Releases Page](https://github.com/sh4dowkey/Crawlx/releases/latest)**.
2.  **Unzip** the package.
3.  **Run the installer script** from your terminal (as Administrator on Windows).
    ```bash
    # On Linux/macOS
    chmod +x ./scripts/setup.sh
    sudo ./scripts/setup.sh
    
    # On Windows (in an Admin PowerShell)
    Set-ExecutionPolicy Bypass -Scope Process -Force
    .\scripts\setup.ps1
    ```
4.  **Crawl!** (Open a new terminal window first)
    ```bash
    crawlx -u [https://toscrape.com](https://toscrape.com) -d 2
    ```

---

## 📦 Installation

This is the recommended method for all users and does not require Go to be installed.

1.  **Download the Installer Package**
    * Go to the **[Latest Release](https://github.com/sh4dowkey/Crawlx/releases/latest)** page.
    * Under the **Assets** section, download the `.zip` (Windows) or `.tar.gz` (Linux/macOS) file for your operating system.

2.  **Unzip the File**
    * Extract the archive. It will contain the pre-compiled `crawlx` binary in the `dist/` folder and an installer script in the `scripts/` folder.

3.  **Run the Installer**
    * Open your terminal and navigate into the top-level unzipped folder (e.g., `cd crawlx-linux-amd64`).
    * Run the installation script. This will copy the binary to a system-wide location, allowing you to run `crawlx` from anywhere.

    * **On Linux or macOS:**
        ```bash
        # Make the script executable
        chmod +x ./scripts/setup.sh

        # Run the installer
        sudo ./scripts/setup.sh
        ```

    * **On Windows (run as Administrator):**
        ```powershell
        # In an Administrator PowerShell, run the installer
        Set-ExecutionPolicy Bypass -Scope Process -Force
        .\scripts\setup.ps1
        ```

That's it! Once the script finishes, open a **new** terminal window, and you will be able to run `crawlx` from any directory. 🎉

---

## 🛠️ Usage

### Options

| Flag | Shorthand | Description | Default |
|---|---|---|---|
| `--url` | `-u` | The starting URL to crawl **(required)** | |
| `--depth` | `-d` | The maximum depth for recursive crawling | `2` |
| `--verbose` | | Enable detailed, verbose output | `false` |

### Examples

**1. Basic crawl:**
```bash
crawlx -u [https://toscrape.com](https://toscrape.com) -d 3

```

**2. Shorthand flags:**
```bash
crawlx -u https://toscrape.com -d 3
```

**3. Verbose crawl:**
```bash
crawlx -u https://toscrape.com -d 2 --verbose
```

Example output:

> ***crawlx.exe --url https:&#8203;//toscrape.com -d 1 --verbose***

<details>
<summary><strong>➡️ Click to see Example Output</strong></summary>

```
 
[INFO] Starting crawl at: 11:43:13 AM IST
[+] Crawling: [https://toscrape.com](https://toscrape.com) (Depth 0)
  ↳ [200 OK] Found 10 links.
  Links found: 
    - [http://quotes.toscrape.com/](http://quotes.toscrape.com/)
    - [http://books.toscrape.com](http://books.toscrape.com)
    ...
  [+] Crawling: [http://quotes.toscrape.com/](http://quotes.toscrape.com/) (Depth 1)
    ↳ [200 OK] Found 49 links.
    ...
  [+] Crawling: [http://books.toscrape.com](http://books.toscrape.com) (Depth 1)
    ↳ [200 OK] Found 73 links.
    ...

============================================================

                   CRAWL SITEMAP REPORT


============================================================

[+] Crawled Pages (200 OK)
  - http://quotes.toscrape.com/ (726ms)  
  - http://quotes.toscrape.com/login (310ms)
  - http://quotes.toscrape.com/js (651ms)
  - http://quotes.toscrape.com/tableful (663ms)
  - http://quotes.toscrape.com/search.aspx (414ms)
  - http://quotes.toscrape.com (415ms)
  ...

------------------------------------------------------------

[~] Redirects (3xx)
  No redirects found.


------------------------------------------------------------

[✗] Client & Server Errors (4xx/5xx)
  No errors found.


------------------------------------------------------------

[!] External Links
  - https://www.goodreads.com/quotes
  - https://www.zyte.com


============================================================

Crawl Summary
  - Crawled 11 pages in 8.648s.
  - Success: 11
  - Warnings: 0
  - Errors: 0


============================================================

All Visited Links
  - [✓] http://quotes.toscrape.com/js
  - [✓] http://quotes.toscrape.com/js-delayed
  - [✓] http://quotes.toscrape.com/search.aspx
  - [✓] http://quotes.toscrape.com/scroll
  - [✓] http://quotes.toscrape.com/random
  - [✓] http://books.toscrape.com
  ...

============================================================.
```

</details>


 

---

## 👨‍💻 For Developers (Building from Source)

If you want to modify the code, you'll need to build the project from the source.

**Prerequisites:**
* [Git](https://git-scm.com/downloads)   [ ➡️ **[Official Go Installation Guide](https://go.dev/doc/install)** ]
* [Go](https://go.dev/dl/) version 1.22+
* [Make](https://www.gnu.org/software/make/) (usually pre-installed on Linux/macOS)

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/sh4dowkey/Crawlx.git
    cd Crawlx
    ```

2.  **Build the Binary**
    The compiled binaries will be placed in the `./dist/` folder.
    ```bash
    # Build for your current OS
    make build

    # Or, build for all platforms
    make build-all
    ```

3.  **(Optional) Install Your Local Build**
    Run the setup script to install your local version system-wide.

    * **On Linux/macOS:**
      ```bash
      chmod +x ./scripts/setup.sh
      sudo ./scripts/setup.sh
      ```
    * **On Windows Powershell (as Administrator):**
      ```powershell
      Set-ExecutionPolicy Bypass -Scope Process -Force
      .\scripts\setup.ps1
      ```

---


## 📂 Project Structure

The project follows the standard Go project layout for better organization and scalability.

```
crawlx/
├── Makefile          # Automates common tasks like building and testing
├── cmd/              # Main application entry point
│   └── crawlx/
│       └── main.go
├── dist/             # Contains compiled binaries after a build
├── internal/         # All private application logic (crawler, parser, etc.)
│   ├── crawl/
│   ├── parse/
│   └── util/
├── scripts/          # Installation and utility scripts (setup.sh, setup.ps1)
├── go.mod            # Go module definition
└── README.md         # This file
```

-----

## 🛣️ Roadmap

This project is actively being developed. Here are the next major features planned:

  - [x] **Recursive Crawling**
  - [ ] **Concurrent Crawling** (Worker Pool Model)
  - [ ] **robots.txt Support** (Polite Crawling)
  - [ ] **Rate Limiting**
  - [ ] **Export Results** (to JSON/CSV)
  - [ ] **Advanced Features** (AI/NLP Content Analysis)

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

---

## 📜 License

This project is licensed under the **MIT License**. See the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.

---

<p align="center">⭐ If you find this project useful, please give it a star on GitHub! ⭐</p>
<p align="center">🕸️ Crawl smarter. Crawl faster. CrawlX.</p>
