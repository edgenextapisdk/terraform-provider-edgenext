# Website Documentation Preview

This directory contains the generated Terraform provider documentation.

## Viewing Documentation Locally

### Method 1: Interactive Preview Page (Recommended)

1. **Generate document list (generates both JSON and JS files for fast loading):**
   ```bash
   cd website
   python3 generate_doc_list.py
   ```
   
   This will generate:
   - `doc_list.json` - JSON format (for fallback)
   - `doc_list.js` - JavaScript format (for fast loading, embedded in HTML)

2. **Start HTTP server:**
   ```bash
   python3 -m http.server 8080
   ```

3. **Open preview page in browser:**
   - Interactive preview: http://localhost:8080/preview.html
   - Features:
     - 📚 All documents organized by category (CDN, SSL, OSS, SCDN)
     - 🔍 Search functionality
     - 📖 Full markdown rendering with syntax highlighting
     - 📊 Document statistics
     - 🎨 GitHub-style markdown styling

### Method 2: Direct File Access

You can also view markdown files directly:

1. **View markdown files directly in browser:**
   - Main page: http://localhost:8080/docs/index.html.markdown
   - Data sources: http://localhost:8080/docs/d/
   - Resources: http://localhost:8080/docs/r/

2. **View specific SCDN documentation:**
   - Data source: http://localhost:8080/docs/d/scdn_domain.html.markdown
   - Resource: http://localhost:8080/docs/r/scdn_domain.html.markdown

### Method 2: Using Markdown Viewer

If you have a markdown viewer installed (like VS Code, Typora, or Marked):

```bash
# Open in VS Code
code website/docs/d/scdn_domain.html.markdown

# Or use any markdown viewer
open website/docs/d/scdn_domain.html.markdown
```

### Method 3: Using Marked (macOS)

If you have Marked app installed:

```bash
open -a Marked website/docs/d/scdn_domain.html.markdown
```

### Method 4: Using Pandoc (Convert to HTML)

```bash
# Install pandoc if needed
# brew install pandoc

# Convert to HTML
pandoc website/docs/d/scdn_domain.html.markdown -o scdn_domain.html

# Open in browser
open scdn_domain.html
```

## Documentation Structure

```
website/
├── docs/
│   ├── index.html.markdown          # Provider main page
│   ├── d/                            # Data sources
│   │   ├── scdn_domain.html.markdown
│   │   ├── scdn_domains.html.markdown
│   │   └── ...
│   └── r/                            # Resources
│       ├── scdn_domain.html.markdown
│       └── ...
└── edgenext.erb                      # Sidebar navigation template
```

## Regenerating Documentation

To regenerate documentation after making changes:

```bash
# From project root
make doc

# Or from gendoc directory
cd gendoc && go run ./...
```

## Notes

- The generated `.html.markdown` files are actually Markdown files with front matter
- They are designed to be processed by Terraform Registry's documentation system
- For local preview, you can view them as regular markdown files
- The front matter (YAML between `---` lines) contains metadata for rendering

## Checking Documentation

To verify all SCDN documentation files are generated:

```bash
# Count SCDN data sources
ls website/docs/d/scdn*.html.markdown | wc -l

# Count SCDN resources
ls website/docs/r/scdn*.html.markdown | wc -l

# List all SCDN docs
ls website/docs/d/scdn*.html.markdown website/docs/r/scdn*.html.markdown
```
