import React from "react";
import App, { Container } from "next/app";
import NextSeo from "next-seo";
import SEO from "../next-seo.config";
import "../skrop_theme/semantic.less";

export default class SkropApp extends App {
  public static async getInitialProps({ Component, router, ctx }) {
    let pageProps = {};

    if (Component.getInitialProps) {
      pageProps = await Component.getInitialProps(ctx);
    }

    return { pageProps };
  }

  public state = { error: null, errorInfo: null };

  public componentDidCatch(error, errorInfo) {
    this.setState({
      error,
      errorInfo
    });
  }

  public render() {
    const { Component, pageProps } = this.props;

    return (
      <Container>
        <div className="page">
          <NextSeo config={SEO} />
          <Component {...pageProps} />
        </div>
      </Container>
    );
  }
}
