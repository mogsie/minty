# InsureQuote Example

A comprehensive example demonstrating all mintydyn patterns with SVG icons.

## Running within minty

This example is designed to run within the minty repository:

```bash
cd minty/examples/insurance-quote
go run ./cmd/main.go
# Open http://localhost:8080
```

## Standalone Version

For a fully self-contained project with Makefile and local minty copy, download the standalone `insurance-quote.zip` from the releases.

## Patterns Demonstrated

| Pattern | Location | Description |
|---------|----------|-------------|
| **States** | Quote wizard, Settings tabs, Compare plans | Tab-like navigation |
| **Rules** | Quote form | Conditional field visibility |
| **ClientFilterable** | Claims page | JSON data with client-side filtering |

## Pages

- `/` - Dashboard with stats and quick actions
- `/quote` - Quote wizard (States + Rules patterns)
- `/claims` - Claims list (ClientFilterable pattern)
- `/compare` - Plan comparison (States with icons)
- `/settings` - User settings (States pattern)

## Icon System

Uses Heroicons (MIT licensed) via `Icon()` and `IconHTML()` functions in `icons.go`.

## See Also

- [AssetTrack](../../../) - Demonstrates ServerRenderedFilter pattern
- [minty documentation](../../docs/)
