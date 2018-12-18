import React from "react";
import App, { Container } from "next/app";
import NextSeo from "next-seo";
import SEO from "../next-seo.config";
import PageMenu from "ui-components/menu/menu";
import "../skrop_theme/semantic.less";
import { string } from "prop-types";

export default class SkropApp extends App {
  public static async getInitialProps({ Component, router, ctx }) {
    let pageProps = {};

    if (Component.getInitialProps) {
      pageProps = await Component.getInitialProps(ctx);
    }

    return { pageProps };
  }

  public state = {
    error: null,
    errorInfo: null,
    filters: Array,
    selectedImageUrl: string
  };

  public componentDidCatch(error, errorInfo) {
    this.setState({
      error,
      errorInfo
    });
  }

  public render() {
    const { Component, pageProps } = this.props;
    const propagatedProps = {
      onImageSelect: this.selectImage
    };

    return (
      <Container>
        <PageMenu logoSrc={require("static/favicon/apple-touch-icon.png")} />
        <div className="page" style={{ marginTop: "7rem" }}>
          <NextSeo config={SEO} />
          <Component {...pageProps} {...this.state} {...propagatedProps} />
        </div>
      </Container>
    );
  }

  private selectImage = selectedImageUrl => {
    this.setState({ selectedImageUrl });
  };
}
