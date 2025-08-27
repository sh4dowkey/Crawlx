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

## Table of Contents

- [✨ Features](#-features)
- [🚀 Quick Start](#-quick-start)
- [📦 Getting Started](#-getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
- [👨‍💻 For Developers (Building from Source)](#-for-developers-building-from-source)
- [📂 Project Structure](#-project-structure)
- [🛣️ Roadmap](#️-roadmap)
- [🤝 Contributing](#-contributing)
- [📜 License](#-license)

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
    sudo ./setup.sh
    
    # On Windows (in an Admin PowerShell)
    .\setup.ps1
    ```
4.  **Crawl!** (Open a new terminal window first)
    ```bash
    crawlx -u [https://toscrape.com](https://toscrape.com) -d 2
    ```

---

## 📦 Getting Started

###  Installation

This is the recommended method for all users and does not require Go to be installed.

1.  **Download the Installer Package**
    * Go to the **[Latest Release](https://github.com/sh4dowkey/Crawlx/releases/latest)** page.
    * Under the **Assets** section, download the `.zip` or `.tar.gz` file for your operating system.

2.  **Unzip the File**
    * Extract the archive. It will contain the pre-compiled `crawlx` binary and an installer script located in the `scripts/` folder.

3.  **Run the Installer**
    * Open your terminal and navigate into the top-level unzipped folder (e.g., `cd crawlx-linux-amd64`).
    * Run the installation script from the `scripts/` directory. This will copy the `crawlx` binary to a system-wide location, allowing you to run it from anywhere.

    * **On Linux or macOS:**
        ```bash
        # Make the script executable
        chmod +x ./scripts/setup.sh

        # Run the installer
        sudo ./scripts/setup.sh
        ```

    * **On Windows (run as Administrator):**
        ```powershell
        # Run the installer from an Administrator PowerShell
        .\scripts\setup.ps1
        ```

That's it! Once the script is finished, you can open a **new** terminal and run `crawlx` from any directory. 🎉

---

## Usage


### Options

| Flag        | Shorthand | Description | Default |
| ---         |   ---     |   ---       |   ---   |
| `--url`     | `-u` | The starting URL to crawl **(required)** | |
| `--depth`   | `-d` | The maximum depth for recursive crawling | `2` |
| `--verbose` | | Enable detailed, verbose output | `false` |

### Examples

**1. Basic crawl with a depth of 3:**

```bash
crawlx -u [https://toscrape.com](https://toscrape.com) -d 3
```

**2. Verbose crawl to see all discovered links:**

```bash
crawlx --url [https://example.com](https://example.com) --depth 2 --verbose
```

-----

## 👨‍💻 For Developers (Building from Source)

If you want to modify the code, you'll need to build the project from the source.

Prerequisites:

   Git

   Go version 1.22+

   Make (usually pre-installed on Linux/macOS)

    
###  1. Clone the repository
```
git clone [https://github.com/sh4dowkey/Crawlx.git](https://github.com/sh4dowkey/Crawlx.git)
cd Crawlx
```

### 2. Build the binary
```
make build
```

### 3. (Optional) Run the installer to install your new version
```
sudo bash ./scripts/setup.sh
```

----

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

-----

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

-----

## 📜 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

-----

<p align="center">⭐ If you find this project useful, please give it a star on GitHub! ⭐</p>
<p align="center">🕸️ Crawl smarter. Crawl faster. CrawlX.</p>
