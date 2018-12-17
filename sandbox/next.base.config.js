module.exports = {
  distDir: "dist",
  useFileSystemPublicRoutes: true,
  publicRuntimeConfig: {
    API_URL: process.env.API_URL || "http://localhost:3000"
  }
};
