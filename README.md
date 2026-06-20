# seekmoon

`seekmoon` is a public AsciiDoc bookshelf for the book project
《包复用生态：发现、管理与评价尺度》.

The repository keeps the source text, research notes, and the local build
workspace used to publish the HTML bookshelf. The default published entry is
the bookshelf catalog.

## Read

The source catalog is [bookshelf/catalog.adoc](bookshelf/catalog.adoc).

After the GitHub Pages workflow is enabled for the repository, the published
site is expected at:

<https://michengliang.github.io/seekmoon/>

The site root redirects to `catalog.html`.

## Contents

```text
bookshelf/
  catalog.adoc                         # Reader entry for the published site
  books/08-package-reuse-ecosystem/    # Main book
  books/07-structured-writing-conventions/
                                       # Writing-structure reference
  shared/                              # Shared AsciiDoc attributes and images
  tools/adoc-books.mjs                 # Local bookshelf build tool
docs/                                  # Research notes and source material
```

## Build

The bookshelf is a private Node workspace because it is a documentation build
environment, not a package intended for registry publication.

```bash
cd bookshelf
pnpm install
pnpm run check
pnpm run build
```

Generated HTML is written to `bookshelf/build/html/`.

## GitHub Pages

GitHub Pages is built from GitHub Actions. The Pages workflow installs the
bookshelf dependencies, runs the structural checks, builds the HTML output, and
deploys `bookshelf/build/html`.

In the GitHub repository settings, configure Pages to use **GitHub Actions** as
the publishing source.

## License

Apache-2.0. See [LICENSE](LICENSE).
