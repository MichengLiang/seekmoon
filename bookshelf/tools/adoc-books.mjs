#!/usr/bin/env node
import { createRequire } from "node:module";
import { cp, mkdir, readdir, readFile, rm, stat, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import asciidoctorFactory from "@asciidoctor/core";
import { parseAbundantTree } from "asciidoc-abundant-tree";
const require = createRequire(import.meta.url);
const DIAGRAM_BLOCK_PATTERN = /^\[(?:actdiag|blockdiag|bpmn|bytefield|c4plantuml|d2|dbml|ditaa|erd|excalidraw|graphviz|mermaid|nomnoml|nwdiag|packetdiag|pikchr|plantuml|rackdiag|seqdiag|svgbob|symbolator|umlet|vega|vegalite|wavedrom|structurizr|diagramsnet|wireviz)(?:,|\])/m;
const CATALOG_BOOK_XREF_PATTERN = /xref:books\/([^/\]]+)\/book\.adoc(?:#[^\[]+)?\[/g;
const XREF_PATTERN = /xref:([^\[#]+)(?:#([A-Za-z0-9_-]+))?\[/g;
const ANCHOR_PATTERN = /^\[#([A-Za-z0-9_-]+)(?:[.,][^\]]*)?\]$/gm;
const LOCAL_TARGET_PATTERN = /\b(?:href|src)="([^"]+)"/g;
const SCHEME_PATTERN = /^[a-zA-Z][a-zA-Z0-9+.-]*:/;
const HOME_MARKER = "data-multi-book-home";
const CONTROLS_MARKER = "data-multi-book-controls";
const FETCH_DIAGRAMS_ENV = "ADOC_BOOKS_FETCH_DIAGRAMS";
const CONFIG_FILE = "adoc-books.config.mjs";
const DEFAULT_CONFIG = {
    rootIndex: {
        redirectTo: "catalog.html",
        title: "AsciiDoc Multi-Book Workspace"
    },
    homeLink: {
        label: "Back to catalog",
        subtitle: "AsciiDoc multi-book workspace"
    }
};
function objectValue(value) {
    return typeof value === "object" && value !== null ? value : {};
}
function stringValue(value, fallback) {
    return typeof value === "string" && value.trim() !== "" ? value : fallback;
}
async function loadRuntimeConfig(rootDir) {
    const configPath = path.join(rootDir, CONFIG_FILE);
    if (!await existsFile(configPath))
        return DEFAULT_CONFIG;
    const configUrl = pathToFileURL(configPath);
    configUrl.search = `mtime=${(await stat(configPath)).mtimeMs}`;
    const module = await import(configUrl.href);
    const rawConfig = objectValue(module.default ?? module);
    const rawRootIndex = objectValue(rawConfig.rootIndex);
    const rawHomeLink = objectValue(rawConfig.homeLink);
    return {
        rootIndex: {
            redirectTo: stringValue(rawRootIndex.redirectTo, DEFAULT_CONFIG.rootIndex.redirectTo),
            title: stringValue(rawRootIndex.title, DEFAULT_CONFIG.rootIndex.title)
        },
        homeLink: {
            label: stringValue(rawHomeLink.label, DEFAULT_CONFIG.homeLink.label),
            subtitle: stringValue(rawHomeLink.subtitle, DEFAULT_CONFIG.homeLink.subtitle)
        }
    };
}
async function existsFile(filePath) {
    try {
        const stats = await stat(filePath);
        return stats.isFile();
    }
    catch {
        return false;
    }
}
async function existsDir(dir) {
    try {
        const stats = await stat(dir);
        return stats.isDirectory();
    }
    catch {
        return false;
    }
}
async function collectFiles(dir, predicate) {
    const entries = await readdir(dir, { withFileTypes: true });
    const files = [];
    for (const entry of entries) {
        const entryPath = path.join(dir, entry.name);
        if (entry.isDirectory()) {
            files.push(...await collectFiles(entryPath, predicate));
        }
        else if (entry.isFile() && predicate(entryPath)) {
            files.push(entryPath);
        }
    }
    return files.sort();
}
async function discoverBooks(rootDir) {
    const booksDir = path.join(rootDir, "books");
    const entries = await readdir(booksDir, { withFileTypes: true });
    const books = [];
    for (const entry of entries) {
        if (!entry.isDirectory())
            continue;
        const bookDir = path.join(booksDir, entry.name);
        const input = path.join(bookDir, "book.adoc");
        if (await existsFile(input)) {
            books.push({
                bookId: entry.name,
                bookDir,
                input,
                htmlOutputDir: path.join(rootDir, "build", "html", "books", entry.name)
            });
        }
    }
    return books.sort((a, b) => a.bookId.localeCompare(b.bookId));
}
async function readIfExists(filePath) {
    return await existsFile(filePath) ? readFile(filePath, "utf8") : "";
}
async function combinedBookSource(bookDir) {
    const files = await collectFiles(bookDir, (filePath) => filePath.endsWith(".adoc") || filePath.endsWith(".mjs"));
    const sources = [];
    for (const file of files) {
        const relativePath = path.relative(path.dirname(bookDir), file).split(path.sep).join("/");
        const source = await readFile(file, "utf8");
        sources.push(`// file: ${relativePath}\n${source.endsWith("\n") ? source : `${source}\n`}`);
    }
    return sources.join("\n");
}
async function workspaceUsesDiagrams(rootDir, books) {
    const sources = [await readIfExists(path.join(rootDir, "catalog.adoc"))];
    for (const book of books)
        sources.push(await combinedBookSource(book.bookDir));
    return sources.some((source) => DIAGRAM_BLOCK_PATTERN.test(source));
}
function createAsciidoctor(loadKroki) {
    const asciidoctor = asciidoctorFactory();
    if (loadKroki) {
        try {
            require("asciidoctor-kroki").register(asciidoctor.Extensions);
        }
        catch (error) {
            const message = error instanceof Error ? error.message : String(error);
            throw new Error(`current workspace contains diagram blocks, but asciidoctor-kroki could not be loaded: ${message}`);
        }
    }
    return asciidoctor;
}
function shouldFetchDiagrams() {
    const value = process.env[FETCH_DIAGRAMS_ENV];
    return value === "1" || value === "true";
}
async function pruneStaleBookHtmlDirs(rootDir, books) {
    const htmlBooksDir = path.join(rootDir, "build", "html", "books");
    if (!await existsDir(htmlBooksDir))
        return;
    const current = new Set(books.map((book) => book.bookId));
    const entries = await readdir(htmlBooksDir, { withFileTypes: true });
    for (const entry of entries) {
        if (entry.isDirectory() && !current.has(entry.name)) {
            await rm(path.join(htmlBooksDir, entry.name), { force: true, recursive: true });
        }
    }
}
function convertFile(asciidoctor, input, outputFile, baseDir, attributes = {}) {
    asciidoctor.convertFile(input, {
        safe: "unsafe",
        base_dir: baseDir,
        to_file: outputFile,
        mkdirs: true,
        attributes
    });
}
async function buildHtml(rootDir, books, asciidoctor, useKroki, fetchDiagrams) {
    await mkdir(path.join(rootDir, "build", "html"), { recursive: true });
    convertFile(asciidoctor, path.join(rootDir, "catalog.adoc"), path.join(rootDir, "build", "html", "catalog.html"), rootDir);
    await pruneStaleBookHtmlDirs(rootDir, books);
    for (const book of books) {
        await mkdir(book.htmlOutputDir, { recursive: true });
        convertFile(asciidoctor, book.input, path.join(book.htmlOutputDir, "book.html"), book.bookDir, {
            ...(useKroki && fetchDiagrams ? { "kroki-fetch-diagram": "", "kroki-http-method": "post" } : {})
        });
    }
}
async function copyAssets(rootDir, books) {
    const sharedImages = path.join(rootDir, "shared", "images");
    if (await existsDir(sharedImages)) {
        await rm(path.join(rootDir, "build", "html", "shared", "images"), { force: true, recursive: true });
        await cp(sharedImages, path.join(rootDir, "build", "html", "shared", "images"), { recursive: true });
    }
    for (const book of books) {
        const assetsDir = path.join(book.bookDir, "assets");
        if (await existsDir(assetsDir)) {
            await rm(path.join(book.htmlOutputDir, "assets"), { force: true, recursive: true });
            await cp(assetsDir, path.join(book.htmlOutputDir, "assets"), { recursive: true });
        }
    }
}
function escapeHtmlAttribute(value) {
    return value
        .replaceAll("&", "&amp;")
        .replaceAll('"', "&quot;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;");
}
function escapeJsonScript(value) {
    return JSON.stringify(value)
        .replaceAll("<", "\\u003c")
        .replaceAll("\u2028", "\\u2028")
        .replaceAll("\u2029", "\\u2029");
}
function escapeJsonValueForScript(value) {
    return JSON.stringify(value)
        .replaceAll("<", "\\u003c")
        .replaceAll("\u2028", "\\u2028")
        .replaceAll("\u2029", "\\u2029");
}
function sectionStyle(section) {
    for (const metadata of section.metadata ?? []) {
        const style = metadata.attributes?.style;
        if (typeof style === "string")
            return style;
    }
    return undefined;
}
function firstSectionId(section) {
    const id = section.ids?.[0];
    if (!id)
        throw new Error(`section is missing an id: ${section.title ?? "(untitled)"}`);
    return id;
}
function childSections(section) {
    return (section.children ?? []).filter((child) => child.kind === "section");
}
function sourceSpan(section) {
    const relativePath = section.source?.relativePath;
    if (!relativePath)
        return undefined;
    return {
        relativePath,
        ...(typeof section.source?.span?.startLine === "number" ? { startLine: section.source.span.startLine } : {}),
        ...(typeof section.source?.span?.endLine === "number" ? { endLine: section.source.span.endLine } : {})
    };
}
function sectionAnchorTree(section) {
    const anchors = [firstSectionId(section)];
    for (const child of childSections(section))
        anchors.push(...sectionAnchorTree(child));
    return anchors;
}
function sectionToc(section, baseDepth = 0, topTitle) {
    const toc = [{ title: topTitle ?? section.title ?? firstSectionId(section), anchor: firstSectionId(section), depth: baseDepth }];
    for (const child of childSections(section)) {
        toc.push(...sectionToc(child, baseDepth + 1));
    }
    return toc;
}
function chapterPage(section, parentPageId) {
    return {
        id: firstSectionId(section),
        kind: "chapter",
        title: section.title ?? firstSectionId(section),
        anchors: sectionAnchorTree(section),
        toc: sectionToc(section, 0, "概述"),
        source: sourceSpan(section),
        ...(parentPageId ? { parentPageId } : {})
    };
}
function backmatterPage(section, kind) {
    return {
        id: firstSectionId(section),
        kind,
        title: section.title ?? firstSectionId(section),
        anchors: sectionAnchorTree(section),
        toc: sectionToc(section, 0, "概述"),
        source: sourceSpan(section)
    };
}
function assertNoToolErrors(bookId, diagnostics) {
    for (const diagnostic of diagnostics) {
        const entry = objectValue(diagnostic);
        if (entry.level === "error") {
            throw new Error(`${bookId}: ${stringValue(entry.code, "error")}: ${stringValue(entry.message, "")}`);
        }
    }
}
function sourceBundleFromDocument(document, fallbackPath, fallbackSource) {
    const sourceFiles = Array.isArray(document.sourceFiles) ? document.sourceFiles : [];
    const files = sourceFiles.filter((file) => typeof file.relativePath === "string" && typeof file.raw === "string");
    if (files.length === 0)
        return fallbackSource;
    return files.map((file) => {
        const body = file.raw?.endsWith("\n") ? file.raw : `${file.raw ?? ""}\n`;
        return `// file: ${file.relativePath ?? fallbackPath}\n${body}`;
    }).join("\n");
}
function fallbackReaderBookData(rootDir, book, fallbackSource, reason) {
    console.warn(`${book.bookId}: reader page map fell back to cover-only mode: ${reason}`);
    const entry = path.relative(rootDir, book.input).split(path.sep).join("/");
    return {
        pageMap: {
            version: 1,
            book: {
                id: book.bookId,
                title: book.bookId,
                entry
            },
            pages: [{
                    id: "cover",
                    kind: "cover",
                    title: book.bookId,
                    anchors: [],
                    toc: []
                }]
        },
        sourceBundle: fallbackSource
    };
}
function buildReaderBookData(rootDir, book, fallbackSource) {
    let document;
    try {
        document = parseAbundantTree({
            sourcePath: book.input,
            mode: "book-entry",
            documentRoot: rootDir
        });
    }
    catch (error) {
        const message = error instanceof Error ? error.message : String(error);
        return fallbackReaderBookData(rootDir, book, fallbackSource, message);
    }
    const diagnostics = Array.isArray(document.toolDiagnostics) ? document.toolDiagnostics : [];
    assertNoToolErrors(book.bookId, diagnostics);
    const title = objectValue(document.title).text;
    const pages = [{
            id: "cover",
            kind: "cover",
            title: typeof title === "string" ? title : book.bookId,
            anchors: [],
            toc: []
        }];
    const sections = (Array.isArray(document.children) ? document.children : []).filter((child) => {
        return typeof child === "object" && child !== null && child.kind === "section";
    });
    const frontmatterStyles = new Set(["abstract", "colophon", "dedication", "preface", "acknowledgments"]);
    const backmatterStyles = new Set(["appendix", "glossary", "bibliography", "index"]);
    const frontmatter = sections.filter((section) => section.level === 1 && frontmatterStyles.has(sectionStyle(section) ?? ""));
    if (frontmatter.length > 0) {
        pages.push({
            id: "frontmatter",
            kind: "frontmatter",
            title: "前置",
            anchors: frontmatter.map(firstSectionId),
            toc: frontmatter.map((section) => ({
                title: section.title ?? firstSectionId(section),
                anchor: firstSectionId(section),
                depth: 0
            }))
        });
    }
    for (const section of sections) {
        const style = sectionStyle(section);
        if (frontmatterStyles.has(style ?? ""))
            continue;
        if (backmatterStyles.has(style)) {
            pages.push(backmatterPage(section, style));
            continue;
        }
        if (section.level === 0) {
            const partId = `part-${firstSectionId(section)}`;
            const children = childSections(section);
            const childPages = children.map((child) => chapterPage(child, partId));
            pages.push({
                id: partId,
                kind: "part",
                title: section.title ?? firstSectionId(section),
                anchors: [firstSectionId(section)],
                toc: [{ title: "概述", anchor: firstSectionId(section), depth: 0 }],
                source: sourceSpan(section),
                childPageIds: childPages.map((page) => page.id)
            });
            pages.push(...childPages);
            continue;
        }
        pages.push(chapterPage(section));
    }
    const pageMap = {
        version: 1,
        book: {
            id: book.bookId,
            title: pages[0].title,
            entry: path.relative(rootDir, book.input).split(path.sep).join("/")
        },
        pages
    };
    return {
        pageMap,
        sourceBundle: sourceBundleFromDocument(document, pageMap.book.entry, fallbackSource)
    };
}
function readerStyles() {
    return `
.multi-book-reader-toggle {
  display: grid;
  gap: .35rem;
  grid-template-columns: 1fr 1fr;
  margin-top: .75rem;
}
.multi-book-reader-toggle button {
  appearance: none;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  color: #1f2937;
  cursor: pointer;
  font-size: .82rem;
  font-weight: 600;
  padding: .42rem .5rem;
}
.multi-book-reader-toggle button[aria-pressed="true"] {
  border-color: #0f766e;
  color: #0f766e;
}
[data-multi-book-page-nav] {
  border-top: 1px solid #e5e7eb;
  margin-top: .75rem;
  padding-top: .75rem;
}
[data-multi-book-page-nav][hidden] {
  display: none;
}
[data-multi-book-page-nav] a {
  border-left: 3px solid transparent;
  color: #334155;
  display: block;
  font-size: .85rem;
  line-height: 1.35;
  padding: .3rem .35rem .3rem .55rem;
  text-decoration: none;
}
[data-multi-book-page-nav] a[data-reader-depth="1"] {
  padding-left: 1.45rem;
}
[data-multi-book-page-nav] a[aria-current="page"] {
  border-left-color: #0f766e;
  color: #0f766e;
  font-weight: 700;
}
.multi-book-reader-cover {
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
}
.multi-book-reader-cover h2 {
  margin-top: 0;
}
[data-multi-book-page-toc] {
  border-left: 1px solid #e5e7eb;
  color: #334155;
  font-size: .85rem;
  padding-left: 1rem;
}
[data-multi-book-page-toc][hidden] {
  display: none;
}
[data-multi-book-page-toc] strong {
  color: #1f2937;
  display: block;
  font-size: .78rem;
  margin-bottom: .35rem;
}
[data-multi-book-page-toc] a {
  color: #334155;
  display: block;
  margin: .35rem 0;
  text-decoration: none;
}
[data-multi-book-pagination] {
  border-top: 1px solid #e5e7eb;
  display: flex;
  gap: 1rem;
  justify-content: space-between;
  margin-bottom: 2rem;
  margin-top: 2.5rem;
  padding-top: 1rem;
}
[data-multi-book-pagination][hidden] {
  display: none;
}
[data-multi-book-pagination] a {
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  color: #1f2937;
  flex: 1 1 0;
  padding: .65rem;
  text-decoration: none;
}
[data-multi-book-pagination] span {
  color: #64748b;
  display: block;
  font-size: .78rem;
}
@media (min-width: 1024px) {
  body.toc2.toc-left.multi-book-reader-paged #content {
    --multi-book-reader-content-gutter: max(0px, calc((100vw - 20em - 1000px) / 2));
    box-sizing: border-box;
    margin-left: var(--multi-book-reader-content-gutter);
    margin-right: 16rem;
    max-width: min(800px, calc(100vw - 20em - 20rem - var(--multi-book-reader-content-gutter)));
    padding-left: 15px;
  }
  body.toc2.toc-left.multi-book-reader-paged [data-multi-book-pagination] {
    --multi-book-reader-content-gutter: max(0px, calc((100vw - 20em - 1000px) / 2));
    margin-left: var(--multi-book-reader-content-gutter);
    max-width: min(800px, calc(100vw - 20em - 20rem - var(--multi-book-reader-content-gutter)));
    width: calc(100vw - 20em - 20rem - var(--multi-book-reader-content-gutter));
  }
  [data-multi-book-page-toc] {
    position: fixed;
    right: 2rem;
    top: 7rem;
    width: 12rem;
  }
}
@media (min-width: 1024px) and (max-width: 1359px) {
  [data-multi-book-page-toc] {
    right: 1rem;
    width: 8rem;
  }
}
body.multi-book-reader-paged #footer {
  display: none;
}
@media (max-width: 767px) {
  #toc.toc2 {
    position: static;
    width: auto;
  }
  [data-multi-book-pagination] {
    flex-direction: column;
  }
  pre, table {
    max-width: 100%;
    overflow-x: auto;
  }
}
`;
}
function readerScript() {
    return `
(function () {
function initReader() {
  var mapElement = document.getElementById("multi-book-page-map");
  if (!mapElement) return;
  var pageMap;
  try {
    pageMap = JSON.parse(mapElement.textContent || "{}");
  } catch (error) {
    return;
  }
  if (!pageMap || !Array.isArray(pageMap.pages)) return;

  var content = document.getElementById("content");
  var toc = document.getElementById("toc");
  if (!content || !toc) return;

  var pages = pageMap.pages;
  var pagesById = new Map();
  var anchorToPageId = new Map();
  pages.forEach(function (page) {
    pagesById.set(page.id, page);
    (page.anchors || []).forEach(function (anchor) { anchorToPageId.set(anchor, page.id); });
  });

  var controls = document.querySelector("[data-multi-book-controls]");
  var toggle = document.createElement("div");
  toggle.className = "multi-book-reader-toggle";
  toggle.setAttribute("data-multi-book-view-toggle", "");
  toggle.innerHTML = '<button type="button" data-reader-view="continuous">连续</button><button type="button" data-reader-view="paged">页面</button>';
  if (controls) controls.appendChild(toggle);

  var nav = document.createElement("nav");
  nav.setAttribute("data-multi-book-page-nav", "");
  nav.setAttribute("aria-label", "页面");
  if (controls && controls.nextSibling) toc.insertBefore(nav, controls.nextSibling);
  else toc.appendChild(nav);

  var pageToc = document.createElement("aside");
  pageToc.setAttribute("data-multi-book-page-toc", "");
  content.parentNode.insertBefore(pageToc, content.nextSibling);

  var pagination = document.createElement("nav");
  pagination.setAttribute("data-multi-book-pagination", "");
  pagination.setAttribute("aria-label", "分页");
  content.parentNode.insertBefore(pagination, pageToc.nextSibling);

  var cover = document.createElement("section");
  cover.id = "multi-book-reader-cover";
  cover.className = "multi-book-reader-cover";
  var sourceTitle = document.querySelector("#header > h1");
  var details = document.querySelector("#header > .details");
  cover.innerHTML = "<h2></h2>";
  cover.querySelector("h2").textContent = sourceTitle ? sourceTitle.textContent : pageMap.book.title;
  if (details) cover.appendChild(details.cloneNode(true));
  content.insertBefore(cover, content.firstChild);

  var originalTocNodes = Array.from(toc.children).filter(function (node) {
    return node !== controls && node !== nav;
  });
  var originalNodes = Array.from(content.children).filter(function (node) { return node !== cover; });
  var currentView = "continuous";
  var currentPageId = "cover";

  function query() {
    return new URLSearchParams(window.location.search);
  }

  function writeQuery(values) {
    var url = new URL(window.location.href);
    Object.keys(values).forEach(function (key) {
      if (values[key] == null) url.searchParams.delete(key);
      else url.searchParams.set(key, values[key]);
    });
    history.pushState(null, "", url);
  }

  function pageNodes(page) {
    if (page.kind === "cover") return [cover];
    if (page.kind === "frontmatter") {
      return (page.anchors || []).map(function (anchor) {
        var heading = document.getElementById(anchor);
        return heading ? heading.closest(".sect1") : null;
      }).filter(Boolean);
    }
    if (page.kind === "part") {
      var partHeading = document.getElementById(page.anchors[0]);
      if (!partHeading) return [];
      var nodes = [partHeading];
      var next = partHeading.nextElementSibling;
      while (next && !next.matches(".sect1") && !next.matches("h1.sect0")) {
        nodes.push(next);
        next = next.nextElementSibling;
      }
      return nodes;
    }
    var heading = document.getElementById(page.anchors[0]);
    var section = heading ? heading.closest(".sect1") : null;
    return section ? [section] : [];
  }

  function renderNav() {
    nav.hidden = currentView !== "paged";
    nav.textContent = "";
    if (currentView !== "paged") return;
    pages.forEach(function (page) {
      var link = document.createElement("a");
      link.href = "?view=paged&page=" + encodeURIComponent(page.id);
      link.textContent = page.title;
      link.dataset.pageId = page.id;
      link.dataset.readerDepth = page.parentPageId ? "1" : "0";
      if (page.id === currentPageId) link.setAttribute("aria-current", "page");
      link.addEventListener("click", function (event) {
        event.preventDefault();
        setView("paged", page.id, true);
      });
      nav.appendChild(link);
    });
  }

  function renderPageToc(page) {
    pageToc.hidden = currentView !== "paged";
    pageToc.textContent = "";
    if (currentView !== "paged") return;
    var title = document.createElement("strong");
    title.textContent = "本页内容";
    pageToc.appendChild(title);
    (page.toc || []).forEach(function (item) {
      var link = document.createElement("a");
      link.href = "#" + item.anchor;
      link.textContent = item.title;
      link.style.paddingLeft = (item.depth * 0.8) + "rem";
      link.addEventListener("click", function (event) {
        event.preventDefault();
        var target = document.getElementById(item.anchor);
        if (target) target.scrollIntoView();
      });
      pageToc.appendChild(link);
    });
  }

  function renderPagination(page) {
    pagination.hidden = currentView !== "paged";
    pagination.textContent = "";
    if (currentView !== "paged") return;
    var index = pages.findIndex(function (candidate) { return candidate.id === page.id; });
    function addLink(label, target) {
      if (!target) {
        var spacer = document.createElement("span");
        pagination.appendChild(spacer);
        return;
      }
      var link = document.createElement("a");
      link.href = "?view=paged&page=" + encodeURIComponent(target.id);
      link.innerHTML = "<span>" + label + "</span>" + target.title;
      link.addEventListener("click", function (event) {
        event.preventDefault();
        setView("paged", target.id, true);
        window.scrollTo(0, 0);
      });
      pagination.appendChild(link);
    }
    addLink("上一页", pages[index - 1]);
    addLink("下一页", pages[index + 1]);
  }

  function setView(view, pageId, persist) {
    currentView = view;
    currentPageId = pagesById.has(pageId) ? pageId : "cover";
    var page = pagesById.get(currentPageId) || pages[0];
    document.body.classList.toggle("multi-book-reader-paged", view === "paged");
    document.body.classList.toggle("multi-book-reader-continuous", view === "continuous");
    toggle.querySelectorAll("button").forEach(function (button) {
      button.setAttribute("aria-pressed", button.dataset.readerView === view ? "true" : "false");
    });
    if (persist) {
      localStorage.setItem("multi-book-reader-view", view);
      writeQuery({ view: view, page: view === "paged" ? currentPageId : null });
    }
    if (view === "continuous") {
      cover.hidden = true;
      originalNodes.forEach(function (node) { node.hidden = false; });
      originalTocNodes.forEach(function (node) { node.hidden = false; });
      pageToc.hidden = true;
      pagination.hidden = true;
      renderNav();
      return;
    }
    originalTocNodes.forEach(function (node) { node.hidden = true; });
    var visible = new Set(pageNodes(page));
    if (visible.size === 0 && page.kind !== "cover") {
      cover.hidden = false;
      cover.querySelector("h2").textContent = "当前页面无法在生成 HTML 中定位。";
      visible.add(cover);
    }
    cover.hidden = !visible.has(cover);
    originalNodes.forEach(function (node) { node.hidden = !visible.has(node); });
    renderNav();
    renderPageToc(page);
    renderPagination(page);
  }

  toggle.addEventListener("click", function (event) {
    var button = event.target.closest("button[data-reader-view]");
    if (!button) return;
    setView(button.dataset.readerView, currentPageId, true);
  });

  document.addEventListener("click", function (event) {
    if (currentView !== "paged") return;
    var link = event.target.closest('a[href^="#"]');
    if (!link) return;
    var anchor = decodeURIComponent(link.getAttribute("href").slice(1));
    var targetPageId = anchorToPageId.get(anchor);
    if (!targetPageId) return;
    if (targetPageId !== currentPageId) {
      event.preventDefault();
      setView("paged", targetPageId, true);
      var target = document.getElementById(anchor);
      if (target) target.scrollIntoView();
    }
  });

  window.addEventListener("popstate", function () {
    var params = query();
    setView(params.get("view") === "paged" ? "paged" : "continuous", params.get("page") || "cover", false);
  });

  var params = query();
  var initialView = params.get("view") || localStorage.getItem("multi-book-reader-view") || "continuous";
  setView(initialView === "paged" ? "paged" : "continuous", params.get("page") || "cover", false);
}
if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", initReader);
} else {
  initReader();
}
}());
`;
}
function addBookControlsToBookHtml(html, href, homeLink, bookSource, pageMap, bookId) {
    if (html.includes(CONTROLS_MARKER))
        return html;
    const marker = '<div id="toc" class="toc2">';
    const index = html.indexOf(marker);
    if (index === -1)
        throw new Error(`${bookId}: missing left TOC container`);
    const insertAt = index + marker.length;
    const controlsBlock = `
<style>
.multi-book-controls {
  margin: 0 0 1rem;
  padding-bottom: .75rem;
  border-bottom: 1px solid #e5e7eb;
}
.multi-book-home {
  color: #1f2937;
  display: block;
  font-weight: 600;
  line-height: 1.35;
  text-decoration: none;
}
.multi-book-home:hover {
  color: #0f766e;
  text-decoration: underline;
}
.multi-book-home span {
  color: #64748b;
  display: block;
  font-size: .78rem;
  font-weight: 400;
  margin-top: .15rem;
}
.multi-book-copy-source {
  appearance: none;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  color: #1f2937;
  cursor: pointer;
  display: block;
  font-size: .82rem;
  font-weight: 600;
  line-height: 1.2;
  margin-top: .75rem;
  padding: .42rem .5rem;
  text-align: center;
  width: 100%;
}
.multi-book-copy-source:hover,
.multi-book-copy-source:focus {
  border-color: #0f766e;
  color: #0f766e;
}
.multi-book-copy-status {
  color: #64748b;
  display: block;
  font-size: .75rem;
  line-height: 1.35;
  margin-top: .4rem;
  min-height: 1em;
}
${readerStyles()}
</style>
<div class="multi-book-controls" ${CONTROLS_MARKER}>
  <a class="multi-book-home" ${HOME_MARKER} href="${escapeHtmlAttribute(href)}">${escapeHtmlAttribute(homeLink.label)}<span>${escapeHtmlAttribute(homeLink.subtitle)}</span></a>
  <button type="button" class="multi-book-copy-source" data-multi-book-source-copy>复制本书为纯文本</button>
  <span class="multi-book-copy-status" data-multi-book-source-status aria-live="polite"></span>
</div>
<script type="application/json" id="multi-book-source-data">${escapeJsonScript(bookSource)}</script>
<script type="application/json" id="multi-book-page-map">${escapeJsonValueForScript(pageMap)}</script>
<script>
(function () {
  var sourceElement = document.getElementById("multi-book-source-data");
  var button = document.querySelector("[data-multi-book-source-copy]");
  var status = document.querySelector("[data-multi-book-source-status]");
  if (!sourceElement || !button) return;

  var source = "";
  var fallbackMode = false;
  try {
    source = JSON.parse(sourceElement.textContent || '""');
  } catch (error) {
    if (status) status.textContent = "纯文本数据读取失败";
    button.disabled = true;
    return;
  }

  function setStatus(message) {
    if (status) status.textContent = message;
  }

  function openSourcePage() {
    var blob = new Blob([source], { type: "text/plain;charset=utf-8" });
    var url = URL.createObjectURL(blob);
    window.open(url, "_blank", "noopener");
    setTimeout(function () { URL.revokeObjectURL(url); }, 60000);
    setStatus("已尝试打开纯文本页；如果没有出现，请允许弹出窗口后再试");
  }

  button.addEventListener("click", async function () {
    if (fallbackMode) {
      openSourcePage();
      return;
    }

    try {
      if (!navigator.clipboard || !navigator.clipboard.writeText) throw new Error("clipboard unavailable");
      await navigator.clipboard.writeText(source);
      button.textContent = "已复制";
      setStatus("");
      setTimeout(function () { button.textContent = "复制本书为纯文本"; }, 1800);
    } catch (error) {
      fallbackMode = true;
      button.textContent = "打开纯文本页";
      openSourcePage();
    }
  });
}());
</script>`;
    const withReaderScript = `${controlsBlock}
<script>${readerScript()}</script>`;
    return `${html.slice(0, insertAt)}${withReaderScript}${html.slice(insertAt)}`;
}
async function addHomeLinks(rootDir, books, homeLink) {
    const catalog = path.join(rootDir, "build", "html", "catalog.html");
    for (const book of books) {
        const htmlFile = path.join(book.htmlOutputDir, "book.html");
        const html = await readFile(htmlFile, "utf8");
        const href = path.relative(path.dirname(htmlFile), catalog);
        const fallbackSource = await combinedBookSource(book.bookDir);
        const readerData = buildReaderBookData(rootDir, book, fallbackSource);
        await writeFile(htmlFile, addBookControlsToBookHtml(html, href, homeLink, readerData.sourceBundle, readerData.pageMap, book.bookId), "utf8");
    }
}
async function writeRootIndex(rootDir, rootIndex) {
    const outputDir = path.join(rootDir, "build", "html");
    await mkdir(outputDir, { recursive: true });
    const redirectTo = rootIndex.redirectTo;
    await writeFile(path.join(outputDir, "index.html"), `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta http-equiv="refresh" content="0; url=${escapeHtmlAttribute(redirectTo)}">
  <title>${escapeHtmlAttribute(rootIndex.title)}</title>
</head>
<body>
  <p><a href="${escapeHtmlAttribute(redirectTo)}">${escapeHtmlAttribute(redirectTo)}</a></p>
</body>
</html>
`, "utf8");
}
function extractCatalogBookIds(catalogSource) {
    const ids = new Set();
    for (const match of catalogSource.matchAll(CATALOG_BOOK_XREF_PATTERN))
        ids.add(match[1]);
    return [...ids].sort((a, b) => a.localeCompare(b));
}
function explicitAnchors(source) {
    const anchors = new Set();
    for (const match of source.matchAll(ANCHOR_PATTERN))
        anchors.add(match[1]);
    return anchors;
}
function attributesFromBookSource(bookSource) {
    const attributes = new Map();
    for (const line of bookSource.split(/\r?\n/)) {
        const match = line.match(/^:([A-Za-z0-9_-]+):\s+(.+)$/);
        if (match)
            attributes.set(match[1], match[2]);
    }
    return attributes;
}
function expandXrefTarget(target, attributes) {
    return target.replaceAll(/\{([A-Za-z0-9_-]+)\}/g, (_, name) => attributes.get(name) ?? `{${name}}`);
}
function issue(code, detail) {
    return { code, detail };
}
async function workspaceContractIssues(rootDir, books) {
    const issues = [];
    const catalogSource = await readIfExists(path.join(rootDir, "catalog.adoc"));
    const catalogBookIds = extractCatalogBookIds(catalogSource);
    const bookSet = new Set(books.map((book) => book.bookId));
    const catalogBookSet = new Set(catalogBookIds);
    for (const bookId of catalogBookIds) {
        if (!bookSet.has(bookId))
            issues.push(issue("CATALOG_TARGET_MISSING", bookId));
    }
    for (const book of books) {
        if (!catalogBookSet.has(book.bookId))
            issues.push(issue("BOOK_MISSING_FROM_CATALOG", book.bookId));
    }
    for (const book of books) {
        const bookSource = await readFile(book.input, "utf8");
        const allSource = await combinedBookSource(book.bookDir);
        const attributes = attributesFromBookSource(bookSource);
        if (!/^:doctype:\s+book$/m.test(bookSource))
            issues.push(issue("MISSING_DOCTYPE", book.bookId));
        for (const match of allSource.matchAll(XREF_PATTERN)) {
            const [, target, anchor] = match;
            if (/^https?:/.test(target))
                continue;
            const expandedTarget = expandXrefTarget(target, attributes);
            if (!expandedTarget.endsWith(".adoc"))
                continue;
            const resolved = path.resolve(book.bookDir, expandedTarget);
            if (!await existsFile(resolved)) {
                issues.push(issue("XREF_TARGET_MISSING", `${book.bookId} -> ${expandedTarget}`));
                if (anchor)
                    issues.push(issue("MISSING_ANCHOR", `${book.bookId} -> ${expandedTarget}#${anchor}`));
                continue;
            }
            if (anchor) {
                const targetSource = await combinedBookSource(path.dirname(resolved));
                if (!explicitAnchors(targetSource).has(anchor)) {
                    issues.push(issue("MISSING_ANCHOR", `${book.bookId} -> ${expandedTarget}#${anchor}`));
                }
            }
        }
    }
    return issues;
}
async function collectHtmlFiles(dir) {
    return collectFiles(dir, (filePath) => filePath.endsWith(".html"));
}
function extractLocalTargets(html) {
    const targets = [];
    for (const match of html.matchAll(LOCAL_TARGET_PATTERN)) {
        const rawTarget = match[1];
        if (rawTarget === "" || rawTarget.startsWith("#") || rawTarget.startsWith("//") || SCHEME_PATTERN.test(rawTarget)) {
            continue;
        }
        const targetWithoutFragment = rawTarget.split("#", 1)[0].split("?", 1)[0];
        if (targetWithoutFragment !== "")
            targets.push(rawTarget);
    }
    return targets;
}
async function missingLocalResources(rootDir) {
    const htmlDir = path.join(rootDir, "build", "html");
    const htmlFiles = await collectHtmlFiles(htmlDir);
    const issues = [];
    for (const htmlFile of htmlFiles) {
        const html = await readFile(htmlFile, "utf8");
        for (const target of extractLocalTargets(html)) {
            const targetPath = target.split("#", 1)[0].split("?", 1)[0];
            const resolved = path.resolve(path.dirname(htmlFile), targetPath);
            if (!await existsFile(resolved)) {
                issues.push(issue("HTML_RESOURCE_MISSING", `${path.relative(htmlDir, htmlFile)} -> ${target}`));
            }
        }
    }
    return issues;
}
async function assertNoIssues(label, issues) {
    if (issues.length === 0)
        return;
    for (const entry of issues)
        console.error(`${entry.code}: ${entry.detail}`);
    throw new Error(`${label} failed with ${issues.length} issue(s)`);
}
export async function buildWorkspace(rootDir = process.cwd()) {
    const catalog = path.join(rootDir, "catalog.adoc");
    if (!await existsFile(catalog))
        throw new Error(`missing catalog.adoc in ${rootDir}`);
    const books = await discoverBooks(rootDir);
    if (books.length === 0)
        throw new Error(`missing book.adoc entries in ${path.join(rootDir, "books")}`);
    const config = await loadRuntimeConfig(rootDir);
    const useKroki = await workspaceUsesDiagrams(rootDir, books);
    const fetchDiagrams = shouldFetchDiagrams();
    const asciidoctor = createAsciidoctor(useKroki);
    await buildHtml(rootDir, books, asciidoctor, useKroki, fetchDiagrams);
    await copyAssets(rootDir, books);
    await addHomeLinks(rootDir, books, config.homeLink);
    await writeRootIndex(rootDir, config.rootIndex);
    await assertNoIssues("HTML local resource check", await missingLocalResources(rootDir));
    await assertNoIssues("workspace contract check", await workspaceContractIssues(rootDir, books));
}
export async function checkWorkspace(rootDir = process.cwd()) {
    await buildWorkspace(rootDir);
}
export async function cleanWorkspace(rootDir = process.cwd()) {
    await rm(path.join(rootDir, "build"), { force: true, recursive: true });
}
async function main() {
    const command = process.argv[2];
    if (command === "build") {
        await buildWorkspace();
        return;
    }
    if (command === "check") {
        await checkWorkspace();
        return;
    }
    if (command === "clean") {
        await cleanWorkspace();
        return;
    }
    console.error("Usage: node tools/adoc-books.mjs <build|check|clean>");
    process.exitCode = 1;
}
const executedPath = process.argv[1] ? path.resolve(process.argv[1]) : "";
const modulePath = fileURLToPath(import.meta.url);
if (executedPath === modulePath) {
    main().catch((error) => {
        console.error(error instanceof Error ? error.message : error);
        process.exitCode = 1;
    });
}
