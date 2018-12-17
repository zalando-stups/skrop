import React from "react";
import { Container, Grid, Header, Button } from "semantic-ui-react";

const MainPage: React.SFC = () => (
  <Container>
    <Grid container>
      <Grid.Row columns="1">
        <Grid.Column>
          <Header as="h1">Skrop</Header>
          <Button primary>Hello</Button>
        </Grid.Column>
      </Grid.Row>
    </Grid>
  </Container>
);

export default MainPage;
