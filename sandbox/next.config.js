const path = require("path");
const withPlugins = require("next-compose-plugins");
const cssPlugin = require("@zeit/next-css");
const typescriptPlugin = require("@zeit/next-typescript");
const bundleAnalyzerPlugin = require("@zeit/next-bundle-analyzer");
const lessPlugin = require("@zeit/next-less");
const optimizeImagesPlugin = require("next-optimized-images");
const nextConfig = require("./next.base.config");

const webpackConfig = {
  webpack: (config, options) => {
    config.resolve.alias = {
      "../../theme.config$": path.join(__dirname, "skrop_theme/theme.config"),
      "ui-components": path.join(__dirname, "ui-components"),
      static: path.join(__dirname, "static")
    };

    config.module.rules.push({
      test: /\.(eot|otf|ttf|woff|woff2)$/,
      use: {
        loader: "url-loader",
        options: {
          limit: 8192,
          publicPath: "./",
          outputPath: "static/css/",
          name: "[name].[ext]"
        }
      }
    });

    return config;
  }
};

module.exports = withPlugins(
  [
    typescriptPlugin,
    cssPlugin,
    lessPlugin,
    optimizeImagesPlugin,
    [
      bundleAnalyzerPlugin,
      {
        analyzeServer: ["server", "both"].includes(process.env.BUNDLE_ANALYZE),
        analyzeBrowser: ["browser", "both"].includes(
          process.env.BUNDLE_ANALYZE
        ),
        bundleAnalyzerConfig: {
          server: {
            analyzerMode: "static",
            reportFilename: "../bundles/server.html"
          },
          browser: {
            analyzerMode: "static",
            reportFilename: "./bundles/client.html"
          }
        }
      }
    ]
  ],
  Object.assign({}, nextConfig, webpackConfig)
);
