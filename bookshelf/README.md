# 包复用生态书架

This workspace contains the Package Reuse Ecosystem book and its writing-structure reference:

- `books/08-package-reuse-ecosystem/` — the primary monograph, 《包复用生态：发现、管理与评价尺度》.
- `books/07-structured-writing-conventions/` — the writing-structure reference used to maintain book 08.

## Install

```bash
pnpm install
```

## Build

```bash
pnpm run build
```

Outputs:

- `build/html/index.html`
- `build/html/catalog.html`
- `build/html/books/08-package-reuse-ecosystem/book.html`
- `build/html/books/07-structured-writing-conventions/book.html`

Each book HTML includes continuous and paged reading views. The copy-source control copies an abundant-tree source bundle for the current book.

## Check

```bash
pnpm run check
```

This checks the catalog, existing books, book doctypes, cross-book xrefs, anchors, and local HTML resources.

## Clean

```bash
pnpm run clean
```

This removes `build/`.

## Navigation

Use `catalog.adoc` as the reader entry. It routes readers to:

- the main 08 book;
- reader-specific paths inside 08;
- the 07 writing conventions reference for maintenance rules.
