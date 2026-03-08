# Page-First Component Strategy

## Default Approach

Start with a complete page template and one page handler.

## When to Split

Split a section into a standalone component + fragment route + handler when at least one is true:

- It updates independently via HTMX.
- It has its own form lifecycle.
- It is reused in multiple page contexts.
- It needs independent loading or refresh triggers.

## Typical Evolution

1. `FeaturePage` renders everything.
2. Extract interactive section into `FeatureSection` component.
3. Add `GET` fragment endpoint for refresh.
4. Add `POST` endpoint(s) for actions.
5. Wire with HTMX attributes (`hx-target`, `hx-swap`, `hx-trigger`).

This keeps complexity localized and preserves server-rendered clarity.
