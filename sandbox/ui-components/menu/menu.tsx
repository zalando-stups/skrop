import React from "react";
import { Container, Image, Menu, Icon } from "semantic-ui-react";

interface IMenu {
  logoSrc: string;
}

const PageMenu: React.SFC<IMenu> = ({ logoSrc }) => (
  <Menu fixed="top" inverted>
    <Container>
      <Menu.Item as="a" href="/" header>
        <Image size="mini" src={logoSrc} alt="Skrop Logo" />
      </Menu.Item>
      <Menu.Item as="a" href="/">
        Home
      </Menu.Item>
      <Menu.Menu position="right">
        <Menu.Item as="a" href="#" target="_blank">
          <Icon name="file alternate outline" size="large" link />
        </Menu.Item>
        <Menu.Item
          as="a"
          href="https://github.com/zalando-stups/skrop"
          target="_blank"
        >
          <Icon name="github" size="large" link />
        </Menu.Item>
      </Menu.Menu>
    </Container>
  </Menu>
);

export default PageMenu;
