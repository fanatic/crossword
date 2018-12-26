import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Home from './Home';
import Mobile from './Mobile';

export default class MobileWrapper extends Component {
  static propTypes = {
    game: PropTypes.object,
    setGameID: PropTypes.func.isRequired
  };

  render() {
    const { game, player_name, setGameID } = this.props;

    return (
      <div className="App">
        {!game && <Home setGameID={setGameID} player_name={player_name} game={game} />}
        {game && <Mobile setGameID={setGameID} player_name={player_name} game={game} />}
      </div>
    );
  }
}
