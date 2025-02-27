# Go Portfolio

An extremely performant, accessible, and SEO-optimized portfolio website built with Go.

## Project Goals

- **Performance**: Optimize for Core Web Vitals and fast loading times
- **Accessibility**: Follow WCAG guidelines for maximum accessibility
- **SEO**: Implement best practices including structured data and semantic HTML
- **Simplicity**: Create a maintainable codebase with clear architecture

## Planned Tech Stack

- Go (Golang) for backend and build processes
- Templ for type-safe HTML templating
- Docker for containerized development
- Make for build automation
- Modern CSS techniques
- Optimized image formats (WebP/AVIF)

## Project Structure (Planned)

```
portfolio/
├── cmd/              # Application entrypoints
│   └── portfolio/    # Main portfolio application
│       └── main.go
├── templates/        # UI components with templ
│   ├── components/
│   │   └── layout.templ
│   └── pages/
│       ├── home.templ
│       └── about.templ
├── static/
│   ├── styles.css
│   └── images/
├── dist/            # Generated output (gitignored)
├── go.mod
└── Makefile
```

## Development (Coming Soon)

Review **Makefile**
Instructions for setting up the development environment will be added as the project progresses.

## Deployment (Coming Soon)

Instructions for deploying to GitHub Pages or Vercel will be added.

## License

MIT
