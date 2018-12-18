import React, { Component } from "react";
import Router from "next/router";
import { Container, Header, Divider, Card } from "semantic-ui-react";

export interface IImage {
  url: string;
}

export interface IMainPage {
  onImageSelect: (imageUrl: string) => void;
  images: Array<IImage>;
}

class MainPage extends Component<IMainPage> {
  static async getInitialProps() {
    // TODO API call to Unsplash here
    const images = [
      {
        url:
          "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      },
      {
        url:
          "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      },
      {
        url:
          "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      },
      {
        url:
          "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      },
      {
        url:
          "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      },
      {
        url:
          "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      },
      {
        url:
          "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      },
      {
        url:
          "https://images.unsplash.com/photo-1544617724-2d30b41f5d05?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
      }
    ];

    return { images };
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
            {images.map(({ url }) => (
              <Card
                image={url}
                onClick={() => {
                  onImageSelect(url);
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
