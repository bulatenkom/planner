# Planner

*⚠️ EXPERIMENTAL For private usage only. Author doesn't guarantee any future maintenance, releases or breaking changes.*

Planner - a basic single-user SPA-like webapp task tracker with built-in UI.

[Video (task create/update/delete)](doc/v0.0.5.webm)\
[Screen 1 (Task form)](doc/v0.0.5.webp)\
[Screen 2 (Task statuses)](doc/v0.0.6.webp)

## Features

- ✅ find/create/delete/update tasks
- filter tasks
  - ✅ by status
  - ❌ by age
  - ❌ by priority
  - ❌ by tag
- ❌ schedule task / create an event linked to tasks
- ❌ calendar
- ❌ task priority
- ❌ expired-task monitoring
- ❌ tags
- ❌ criteria-search (by tag / name-pattern / ...)

## Usage

You shouldn't.

## Idiology of the repo/project

- Go mod doesn't contain deps
- Go files sit in one 'main' package
- Bare minimum CSS library (Just to make it more visually appealing)
- No direct/manual code in JS/TS/Hyperscript/whatever (one exception is invocation of markdown renderer on page).
- Author has no prior Go/HTMX knowledge
- Zero-motivation to finish this project

## Dependencies / Credits

- [Missing.css](https://missing.style/)
- [HTMX](https://htmx.org/)
- [Marked.js](https://marked.js.org/)

---
<div style="text-align: right">MIT</div>