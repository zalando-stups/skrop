import React from "react";
import Router from "next/router";
import { Container, Header, Button } from "semantic-ui-react";

export interface IMainPage {
  onImageSelect: (imageUrl: string) => void;
}

const MainPage: React.SFC<IMainPage> = ({ onImageSelect }) => (
  <div>
    <Container>
      <Header as="h1" textAlign="center">
        Select Image
      </Header>
      <Button
        onClick={() => {
          onImageSelect(
            "https://images.unsplash.com/photo-1544868501-b2f493d76e35?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=500&q=60"
          );
          Router.push("/artboard");
        }}
        primary
      >
        Simulate Image Select
      </Button>
    </Container>
  </div>
);
export default MainPage;
