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
- [📦 Installation](#-installation)
- [🛠️ Usage](#️-usage)
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

***crawlx.exe --url https://toscrape.com -d 1 --verbose***
<details>
<summary><strong>Click to see Example Output</strong></summary>

```
 
[INFO] Starting crawl at: 11:43:13 AM IST
[+] Crawling: https://toscrape.com (Depth 0)
  ↳ [200 OK] Found 10 links.
  Links found: 
    - http://quotes.toscrape.com/
    - http://quotes.toscrape.com/login
    - http://quotes.toscrape.com/js
    - http://quotes.toscrape.com/tableful
    - http://quotes.toscrape.com/js-delayed
    - http://quotes.toscrape.com/search.aspx
    - http://quotes.toscrape.com/scroll
    - http://quotes.toscrape.com
    - http://quotes.toscrape.com/random
    - http://books.toscrape.com
  [+] Crawling: http://quotes.toscrape.com/ (Depth 1)
    ↳ [200 OK] Found 49 links.
    Links found:
      - http://quotes.toscrape.com/tag/simile/
      - http://quotes.toscrape.com/author/Jane-Austen
      - http://quotes.toscrape.com/tag/change/page/1/
      - http://quotes.toscrape.com/tag/deep-thoughts/page/1/
      - http://quotes.toscrape.com/tag/misattributed-eleanor-roosevelt/page/1/
      - http://quotes.toscrape.com/tag/live/page/1/
      - http://quotes.toscrape.com/login
      - https://www.goodreads.com/quotes
      - http://quotes.toscrape.com/tag/books/page/1/
      - http://quotes.toscrape.com/tag/obvious/page/1/
      - http://quotes.toscrape.com/tag/simile/page/1/
      - http://quotes.toscrape.com/tag/miracles/page/1/
      - http://quotes.toscrape.com/tag/be-yourself/page/1/
      - http://quotes.toscrape.com/tag/humor/
      - http://quotes.toscrape.com/tag/paraphrased/page/1/
      - http://quotes.toscrape.com/tag/humor/page/1/
      - http://quotes.toscrape.com/tag/thinking/page/1/
      - http://quotes.toscrape.com/tag/books/
      - http://quotes.toscrape.com/tag/adulthood/page/1/
      - http://quotes.toscrape.com/tag/life/
      - http://quotes.toscrape.com/tag/reading/
      - http://quotes.toscrape.com/tag/edison/page/1/
      - http://quotes.toscrape.com/tag/classic/page/1/
      - http://quotes.toscrape.com/author/Marilyn-Monroe
      - http://quotes.toscrape.com/tag/friends/
      - http://quotes.toscrape.com/author/J-K-Rowling
      - http://quotes.toscrape.com/tag/friendship/
      - http://quotes.toscrape.com/author/Thomas-A-Edison
      - http://quotes.toscrape.com/tag/truth/
      - http://quotes.toscrape.com/tag/miracle/page/1/
      - http://quotes.toscrape.com/author/Andre-Gide
      - http://quotes.toscrape.com/tag/failure/page/1/
      - http://quotes.toscrape.com/tag/world/page/1/
      - http://quotes.toscrape.com/tag/inspirational/
      - http://quotes.toscrape.com/tag/choices/page/1/
      - http://quotes.toscrape.com/tag/love/page/1/
      - http://quotes.toscrape.com/author/Albert-Einstein
      - http://quotes.toscrape.com/tag/aliteracy/page/1/
      - http://quotes.toscrape.com/page/2/
      - https://www.zyte.com
      - http://quotes.toscrape.com/author/Eleanor-Roosevelt
      - http://quotes.toscrape.com/tag/love/
      - http://quotes.toscrape.com/author/Steve-Martin
      - http://quotes.toscrape.com/tag/abilities/page/1/
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/tag/success/page/1/
      - http://quotes.toscrape.com/tag/value/page/1/
      - http://quotes.toscrape.com/tag/life/page/1/
      - http://quotes.toscrape.com/tag/inspirational/page/1/
  [+] Crawling: http://quotes.toscrape.com/login (Depth 1)
    ↳ [200 OK] Found 4 links.
    Links found:
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/login
      - https://www.goodreads.com/quotes
      - https://www.zyte.com
  [+] Crawling: http://quotes.toscrape.com/js (Depth 1)
    ↳ [200 OK] Found 5 links.
    Links found:
      - https://www.zyte.com
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/login
      - http://quotes.toscrape.com/js/page/2/
      - https://www.goodreads.com/quotes
  [+] Crawling: http://quotes.toscrape.com/tableful (Depth 1)
    ↳ [200 OK] Found 35 links.
    Links found:
      - http://quotes.toscrape.com/tableful/page/2/
      - http://quotes.toscrape.com/tableful/tag/live/page/1/
      - http://quotes.toscrape.com/tableful/tag/aliteracy/page/1/
      - http://quotes.toscrape.com/tableful/tag/classic/page/1/
      - http://quotes.toscrape.com/tableful/tag/books/page/1/
      - http://quotes.toscrape.com/tableful/tag/humor/page/1/
      - http://quotes.toscrape.com/tableful/tag/miracle/page/1/
      - http://quotes.toscrape.com/login
      - http://quotes.toscrape.com/tableful/tag/be-yourself/page/1/
      - http://quotes.toscrape.com/tableful/tag/choices/page/1/
      - http://quotes.toscrape.com/tableful/tag/reading/page/1/
      - https://www.zyte.com
      - http://quotes.toscrape.com/tableful/tag/truth/page/1/
      - http://quotes.toscrape.com/tableful/tag/change/page/1/
      - http://quotes.toscrape.com/tableful/tag/love/page/1/
      - http://quotes.toscrape.com/tableful/tag/friends/page/1/
      - http://quotes.toscrape.com/tableful/tag/thinking/page/1/
      - http://quotes.toscrape.com/tableful/tag/life/page/1/
      - http://quotes.toscrape.com/tableful/tag/failure/page/1/
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/tableful/tag/world/page/1/
      - https://www.goodreads.com/quotes
      - http://quotes.toscrape.com/tableful/tag/miracles/page/1/
      - http://quotes.toscrape.com/tableful/tag/abilities/page/1/
      - http://quotes.toscrape.com/tableful/tag/obvious/page/1/
      - http://quotes.toscrape.com/tableful/tag/simile/page/1/
      - http://quotes.toscrape.com/tableful/tag/inspirational/page/1/
      - http://quotes.toscrape.com/tableful/tag/misattributed-eleanor-roosevelt/page/1/
      - http://quotes.toscrape.com/tableful/tag/edison/page/1/
      - http://quotes.toscrape.com/tableful/tag/success/page/1/
      - http://quotes.toscrape.com/tableful/tag/paraphrased/page/1/
      - http://quotes.toscrape.com/tableful/tag/deep-thoughts/page/1/
      - http://quotes.toscrape.com/tableful/tag/adulthood/page/1/
      - http://quotes.toscrape.com/tableful/tag/friendship/page/1/
      - http://quotes.toscrape.com/tableful/tag/value/page/1/
  [+] Crawling: http://quotes.toscrape.com/js-delayed (Depth 1)
    ↳ [200 OK] Found 5 links.
    Links found:
      - https://www.goodreads.com/quotes
      - https://www.zyte.com
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/login
      - http://quotes.toscrape.com/js-delayed/page/2/
  [+] Crawling: http://quotes.toscrape.com/search.aspx (Depth 1)
    ↳ [200 OK] Found 4 links.
    Links found:
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/login
      - https://www.goodreads.com/quotes
      - https://www.zyte.com
  [+] Crawling: http://quotes.toscrape.com/scroll (Depth 1)
    ↳ [200 OK] Found 4 links.
    Links found:
      - https://www.zyte.com
      - https://www.goodreads.com/quotes
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/login
  [+] Crawling: http://quotes.toscrape.com (Depth 1)
    ↳ [200 OK] Found 49 links.
    Links found:
      - http://quotes.toscrape.com/page/2/
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/author/Albert-Einstein
      - http://quotes.toscrape.com/author/J-K-Rowling
      - http://quotes.toscrape.com/tag/misattributed-eleanor-roosevelt/page/1/
      - http://quotes.toscrape.com/tag/world/page/1/
      - http://quotes.toscrape.com/tag/miracle/page/1/
      - http://quotes.toscrape.com/tag/deep-thoughts/page/1/
      - https://www.goodreads.com/quotes
      - http://quotes.toscrape.com/tag/adulthood/page/1/
      - http://quotes.toscrape.com/tag/thinking/page/1/
      - http://quotes.toscrape.com/tag/inspirational/
      - http://quotes.toscrape.com/tag/life/page/1/
      - http://quotes.toscrape.com/tag/simile/page/1/
      - http://quotes.toscrape.com/tag/love/page/1/
      - http://quotes.toscrape.com/tag/humor/
      - http://quotes.toscrape.com/tag/failure/page/1/
      - http://quotes.toscrape.com/tag/love/
      - http://quotes.toscrape.com/tag/be-yourself/page/1/
      - http://quotes.toscrape.com/tag/paraphrased/page/1/
      - http://quotes.toscrape.com/tag/truth/
      - http://quotes.toscrape.com/tag/aliteracy/page/1/
      - http://quotes.toscrape.com/tag/classic/page/1/
      - http://quotes.toscrape.com/tag/inspirational/page/1/
      - http://quotes.toscrape.com/tag/value/page/1/
      - http://quotes.toscrape.com/tag/live/page/1/
      - http://quotes.toscrape.com/tag/simile/
      - http://quotes.toscrape.com/tag/books/page/1/
      - http://quotes.toscrape.com/tag/reading/
      - http://quotes.toscrape.com/tag/abilities/page/1/
      - http://quotes.toscrape.com/author/Steve-Martin
      - http://quotes.toscrape.com/tag/success/page/1/
      - http://quotes.toscrape.com/author/Jane-Austen
      - http://quotes.toscrape.com/tag/humor/page/1/
      - http://quotes.toscrape.com/tag/choices/page/1/
      - http://quotes.toscrape.com/login
      - http://quotes.toscrape.com/tag/friends/
      - http://quotes.toscrape.com/author/Thomas-A-Edison
      - http://quotes.toscrape.com/tag/miracles/page/1/
      - http://quotes.toscrape.com/tag/edison/page/1/
      - http://quotes.toscrape.com/tag/obvious/page/1/
      - http://quotes.toscrape.com/tag/life/
      - http://quotes.toscrape.com/tag/books/
      - http://quotes.toscrape.com/tag/friendship/
      - http://quotes.toscrape.com/author/Eleanor-Roosevelt
      - http://quotes.toscrape.com/author/Marilyn-Monroe
      - https://www.zyte.com
      - http://quotes.toscrape.com/tag/change/page/1/
      - http://quotes.toscrape.com/author/Andre-Gide
  [+] Crawling: http://quotes.toscrape.com/random (Depth 1)
    ↳ [200 OK] Found 6 links.
    Links found:
      - http://quotes.toscrape.com/
      - http://quotes.toscrape.com/login
      - http://quotes.toscrape.com/tag/truth/page/1/
      - https://www.goodreads.com/quotes
      - https://www.zyte.com
      - http://quotes.toscrape.com/author/J-K-Rowling
  [+] Crawling: http://books.toscrape.com (Depth 1)
    ↳ [200 OK] Found 73 links.
    Links found:
      - http://books.toscrape.com/catalogue/page-2.html
      - http://books.toscrape.com/catalogue/the-requiem-red_995/index.html
      - http://books.toscrape.com/catalogue/category/books/academic_40/index.html
      - http://books.toscrape.com/catalogue/category/books/sports-and-games_17/index.html
      - http://books.toscrape.com/catalogue/soumission_998/index.html
      - http://books.toscrape.com/catalogue/category/books/short-stories_45/index.html
      - http://books.toscrape.com/catalogue/a-light-in-the-attic_1000/index.html
      - http://books.toscrape.com/catalogue/category/books/suspense_44/index.html
      - http://books.toscrape.com/catalogue/category/books/childrens_11/index.html
      - http://books.toscrape.com/catalogue/category/books/thriller_37/index.html
      - http://books.toscrape.com/catalogue/category/books/contemporary_38/index.html
      - http://books.toscrape.com/catalogue/category/books/sequential-art_5/index.html
      - http://books.toscrape.com/catalogue/category/books/add-a-comment_18/index.html
      - http://books.toscrape.com/catalogue/category/books/fiction_10/index.html
      - http://books.toscrape.com/catalogue/scott-pilgrims-precious-little-life-scott-pilgrim-1_987/index.html
      - http://books.toscrape.com/catalogue/category/books/fantasy_19/index.html
      - http://books.toscrape.com/catalogue/category/books/nonfiction_13/index.html
      - http://books.toscrape.com/catalogue/category/books/adult-fiction_29/index.html
      - http://books.toscrape.com/catalogue/category/books/biography_36/index.html
      - http://books.toscrape.com/catalogue/category/books/travel_2/index.html
      - http://books.toscrape.com/catalogue/category/books/religion_12/index.html
      - http://books.toscrape.com/catalogue/category/books/parenting_28/index.html
      - http://books.toscrape.com/catalogue/category/books/womens-fiction_9/index.html
      - http://books.toscrape.com/catalogue/category/books/historical_42/index.html
      - http://books.toscrape.com/catalogue/mesaerion-the-best-science-fiction-stories-1800-1849_983/index.html
      - http://books.toscrape.com/catalogue/category/books/psychology_26/index.html
      - http://books.toscrape.com/catalogue/category/books/self-help_41/index.html
      - http://books.toscrape.com/catalogue/the-boys-in-the-boat-nine-americans-and-their-epic-quest-for-gold-at-the-1936-berlin-olympics_992/index.html
      - http://books.toscrape.com/catalogue/category/books/christian_43/index.html
      - http://books.toscrape.com/catalogue/category/books/business_35/index.html
      - http://books.toscrape.com/catalogue/category/books/poetry_23/index.html
      - http://books.toscrape.com/catalogue/rip-it-up-and-start-again_986/index.html
      - http://books.toscrape.com/catalogue/olio_984/index.html
      - http://books.toscrape.com/catalogue/category/books/mystery_3/index.html
      - http://books.toscrape.com/index.html
      - http://books.toscrape.com/catalogue/category/books_1/index.html
      - http://books.toscrape.com/catalogue/category/books/science_22/index.html
      - http://books.toscrape.com/catalogue/sapiens-a-brief-history-of-humankind_996/index.html
      - http://books.toscrape.com/catalogue/category/books/autobiography_27/index.html
      - http://books.toscrape.com/catalogue/category/books/art_25/index.html
      - http://books.toscrape.com/catalogue/category/books/novels_46/index.html
      - http://books.toscrape.com/catalogue/category/books/humor_30/index.html
      - http://books.toscrape.com/catalogue/category/books/erotica_50/index.html
      - http://books.toscrape.com/catalogue/category/books/philosophy_7/index.html
      - http://books.toscrape.com/catalogue/category/books/food-and-drink_33/index.html
      - http://books.toscrape.com/catalogue/sharp-objects_997/index.html
      - http://books.toscrape.com/catalogue/category/books/horror_31/index.html
      - http://books.toscrape.com/catalogue/category/books/music_14/index.html
      - http://books.toscrape.com/catalogue/the-black-maria_991/index.html
      - http://books.toscrape.com/catalogue/category/books/historical-fiction_4/index.html
      - http://books.toscrape.com/catalogue/category/books/politics_48/index.html
      - http://books.toscrape.com/catalogue/category/books/paranormal_24/index.html
      - http://books.toscrape.com/catalogue/category/books/crime_51/index.html
      - http://books.toscrape.com/catalogue/category/books/romance_8/index.html
      - http://books.toscrape.com/catalogue/set-me-free_988/index.html
      - http://books.toscrape.com/catalogue/its-only-the-himalayas_981/index.html
      - http://books.toscrape.com/catalogue/the-dirty-little-secrets-of-getting-your-dream-job_994/index.html
      - http://books.toscrape.com/catalogue/category/books/christian-fiction_34/index.html
      - http://books.toscrape.com/catalogue/tipping-the-velvet_999/index.html
      - http://books.toscrape.com/catalogue/the-coming-woman-a-novel-based-on-the-life-of-the-infamous-feminist-victoria-woodhull_993/index.html
      - http://books.toscrape.com/catalogue/category/books/health_47/index.html
      - http://books.toscrape.com/catalogue/category/books/default_15/index.html
      - http://books.toscrape.com/catalogue/category/books/cultural_49/index.html
      - http://books.toscrape.com/catalogue/category/books/science-fiction_16/index.html
      - http://books.toscrape.com/catalogue/our-band-could-be-your-life-scenes-from-the-american-indie-underground-1981-1991_985/index.html
      - http://books.toscrape.com/catalogue/category/books/new-adult_20/index.html
      - http://books.toscrape.com/catalogue/category/books/young-adult_21/index.html
      - http://books.toscrape.com/catalogue/category/books/spirituality_39/index.html
      - http://books.toscrape.com/catalogue/category/books/history_32/index.html
      - http://books.toscrape.com/catalogue/category/books/classics_6/index.html
      - http://books.toscrape.com/catalogue/shakespeares-sonnets_989/index.html
      - http://books.toscrape.com/catalogue/starving-hearts-triangular-trade-trilogy-1_990/index.html
      - http://books.toscrape.com/catalogue/libertarianism-for-beginners_982/index.html

```

</details>

```
 
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
  - https://toscrape.com (1553ms)
  - http://quotes.toscrape.com/js-delayed (895ms)
  - http://quotes.toscrape.com/scroll (502ms)
  - http://quotes.toscrape.com/random (369ms)
  - http://books.toscrape.com (1041ms)


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
  - [✓] https://toscrape.com
  - [✓] http://quotes.toscrape.com/
  - [✓] http://quotes.toscrape.com/login
  - [✓] http://quotes.toscrape.com/tableful
  - [✓] http://quotes.toscrape.com

============================================================.
```

---

## 👨‍💻 For Developers (Building from Source)

If you want to modify the code, you'll need to build the project from the source.

**Prerequisites:**
* [Git](https://git-scm.com/downloads)
* [Go](https://go.dev/dl/) version 1.22+
* [Make](https://www.gnu.org/software/make/) (usually pre-installed on Linux/macOS)

1.  **Clone the Repository**
    ```bash
    git clone [https://github.com/sh4dowkey/Crawlx.git](https://github.com/sh4dowkey/Crawlx.git)
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
