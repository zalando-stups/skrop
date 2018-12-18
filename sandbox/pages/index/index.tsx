import React, { Component } from "react";
import Router from "next/router";
// import getConfig from "next/config";
// import Unsplash from "unsplash-js";
import { Container, Header, Divider, Card } from "semantic-ui-react";

export interface IUnsplashImage {
  raw: string;
  full: string;
  small: string;
  thumb: string;
}

export interface IMainPage {
  onImageSelect: (imageUrl: string) => void;
  images: Array<IUnsplashImage>;
}

class MainPage extends Component<IMainPage> {
  static async getInitialProps() {
    // TODO: add Unsplash support
    // const {
    //   publicRuntimeConfig: { UNSPLASH_APP_ID, UNSPLASH_APP_SECRET }
    // } = getConfig();

    // const unsplash = new Unsplash({
    //   applicationId: UNSPLASH_APP_ID,
    //   secret: UNSPLASH_APP_SECRET
    // });

    // try {
    //   const imagesRequest = await unsplash.search.photos("fashion", 1, 72);
    //   const jsonImages = await imagesRequest.json();
    //   const images = jsonImages.results
    //     .filter(({ width, height }) => height > width)
    //     .map(({ urls }) => urls);
    //   return { images };
    // } catch (e) {
    //   console.log(e);
    //   return { images: [] };
    // }

    return {
      images: [
        {
          small:
            "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        },
        {
          small:
            "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        },
        {
          small:
            "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        },
        {
          small:
            "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        },
        {
          small:
            "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        },
        {
          small:
            "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        },
        {
          small:
            "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        },
        {
          small:
            "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
        }
      ]
    };
  }

  render() {
    const { onImageSelect, images } = this.props;
    return (
      <div>
        <Container>
          <Header as="h1" textAlign="center">
            Select Image
          </Header>
          <Divider hidden />
          <Card.Group itemsPerRow={4}>
            {images.map(({ small }) => (
              <Card
                image={small}
                onClick={() => {
                  onImageSelect(small);
                  Router.push("/artboard");
                }}
              />
            ))}
          </Card.Group>
        </Container>
      </div>
    );
  }
}
export default MainPage;
