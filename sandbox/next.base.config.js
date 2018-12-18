module.exports = {
  distDir: "dist",
  useFileSystemPublicRoutes: true,
  publicRuntimeConfig: {
    UNSPLASH_APP_ID: process.env.UNSPLASH_APP_ID || "",
    UNSPLASH_APP_SECRET: process.env.UNSPLASH_APP_SECRET || ""
  }
};
