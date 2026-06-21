# mooncakes-http-reader-probe

MoonBit native HTTP reader probe for Mooncakes API, asset, and source zip endpoints.

This spike verifies whether `oboard/mio` can support typed reader objects similar to a small `reqwest` or `httpx` wrapper:

- `MooncakesApiReader` reads JSON API endpoints.
- `MooncakesAssetReader` reads `module_index.json`, `package_data.json`, and optional `resource.json` assets.
- `SourceZipReader` reads published source zip bytes.

Run:

```bash
moon check --target native
moon build --target native
moon run cmd/main --target native
```

The output records status code, header count, body byte length, JSON parsing outcome, and a short preview.
