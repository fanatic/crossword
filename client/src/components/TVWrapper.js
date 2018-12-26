import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { ConnectedGameWrapper } from './GameWrapper';
import TV from './TV';

export default class TVWrapper extends Component {
  static propTypes = {
    match: PropTypes.object.isRequired
  };

  render() {
    const {
      match: {
        params: { id }
      }
    } = this.props;

    return (
      <ConnectedGameWrapper game_id={id} player_name={'TV'}>
        <TV />
      </ConnectedGameWrapper>
    );
  }
}
