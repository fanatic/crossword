import React, { Component } from 'react';
import { connect, PromiseState } from 'react-refetch';
import PropTypes from 'prop-types';
import MobileWrapper from './MobileWrapper';

export default class GameState extends Component {
  constructor(props) {
    super(props);
    this.state = {
      game_id: localStorage.getItem('game_id'),
      player_name: localStorage.getItem('player_name')
    };
  }

  setGameID = (game_id, player_name) => {
    this.setState({ game_id, player_name });
    localStorage.setItem('game_id', game_id);
    localStorage.setItem('player_name', player_name);
  };

  render() {
    return (
      <ConnectedGameWrapper game_id={this.state.game_id} player_name={this.state.player_name}>
        <MobileWrapper setGameID={this.setGameID} />
      </ConnectedGameWrapper>
    );
  }
}

class GameWrapper extends Component {
  static propTypes = {
    getGame: PropTypes.instanceOf(PromiseState).isRequired
  };

  render() {
    const { player_name, getGame, children } = this.props;

    let game;
    if (getGame.fulfilled) {
      game = getGame.value;
    }

    const childrenWithProps = React.Children.map(children, child => React.cloneElement(child, { player_name, game }));

    return <div>{childrenWithProps}</div>;
  }
}

export const ConnectedGameWrapper = connect(props => ({
  getGame: {
    url: `http://192.168.3.38:8080/games/${props.game_id}`,
    refreshInterval: 1000
  }
}))(GameWrapper);
