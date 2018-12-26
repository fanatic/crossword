import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Home from './Home';
import Mobile from './Mobile';
import TV from './TV';

export default class Debug extends Component {
  static propTypes = {
    game: PropTypes.object.isRequired,
    setGameID: PropTypes.func.isRequired
  };

  render() {
    const { game, player_name, setGameID } = this.props;

    return (
      <div className="App">
        {!game && <Home setGameID={setGameID} player_name={player_name} game={game} />}
        <hr />
        {game && (
          <React.Fragment>
            <Mobile setGameID={setGameID} player_name={player_name} game={game} />
            <hr />
            <TV game={game} />
          </React.Fragment>
        )}
      </div>
    );
  }
}
