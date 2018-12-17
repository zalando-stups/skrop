import React from "react";
import { Grid, Header, Tab } from "semantic-ui-react";

export interface IArtboard {
  selectedImageUrl: string;
}

const panes = [
  {
    menuItem: "Properties",
    render: () => <Tab.Pane>Filter Props</Tab.Pane>
  },
  { menuItem: "Export", render: () => <Tab.Pane>Export Content</Tab.Pane> }
];

const Artboard: React.SFC<IArtboard> = ({ selectedImageUrl }) => (
  <div>
    <Grid columns={3}>
      <Grid.Row>
        <Grid.Column width="4">
          <Header as="h3">Filters</Header>
        </Grid.Column>
        <Grid.Column>
          <img src={selectedImageUrl} />{" "}
        </Grid.Column>
        <Grid.Column width="4">
          <Header as="h3">Filters</Header>
          <Tab panes={panes} />
        </Grid.Column>
      </Grid.Row>
    </Grid>
  </div>
);
export default Artboard;
