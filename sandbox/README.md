## Requirements

- Node.js 10
- NPM 6
- Docker

## Languages & Frameworks

Skrop Sandbox Web Application is build using the following:

- TypeScript & LESS
- [Next.js](https://github.com/zeit/next.js/) as a React framework with a bunch of plugins to simplify SSR, code splitting, images optimization, and SEO
- [Semantic UI React](https://react.semantic-ui.com/) with custom theme as Design System and building blocks for UI
- [Express.js](https://expressjs.com/) as a web framework to run server side code

## Setup

Please follow these steps to initialize the project on your machine:

1. Clone this repo, grab some ☕️ and run:
1. In order to make sure that the project is compilable, build it for the production environment locally:

```bash
# build the project
npm run build

# start server
npm start
```

## Run local development environment

```bash
npm run dev
```

## Run production environment

Since we dockerize Web Application to run it in production
it's possible to check production build localy:

1. Build the project

```bash
npm run build
```

2. Build new Docker image

```bash
./scripts/docker-build.sh
```

3. Run Docker image

```bash
./scripts/docker-run.sh
```

Now you should see application running at `http://localhost:3000`.
