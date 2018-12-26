import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Board from './Board';
import './TV.css';

export default class TV extends Component {
  static propTypes = {
    game: PropTypes.object.isRequired
  };

  render() {
    const { game } = this.props;

    if (!game) return <div />;

    const expires = new Date(game.current_clue.expires_at);
    const now = new Date();
    const seconds_left = Math.round((+expires - +now) / 1000);

    return (
      <div className="App">
        <div className="tv-layout">
          <div className="tv-left">
            {game.current_players && (
              <div>
                <span className="score-header">Score</span>
                <ul className="scoreboard">
                  {game.current_players.map(p => (
                    <li key={p.name}>
                      {p.name}: {p.current_score} pts {!p.active && '(inactive)'}
                    </li>
                  ))}
                </ul>
              </div>
            )}
            {game.last_clue && game.last_clue.guesses && game.last_clue.guesses.length > 0 && (
              <div>
                <span className="last-clue-header">
                  Last Clue: <br />
                </span>
                <span className="last-clue">{game.last_clue.description}</span>
                <ul className="scoreboard-small">
                  {(game.last_clue.guesses || []).map(g => (
                    <li key={g.player.name}>
                      {g.player.name}: {g.guess} ({g.score} pts)
                    </li>
                  ))}
                </ul>
              </div>
            )}
            {game.current_clue && <p className="current-clue">{game.current_clue.description}</p>}
            <span className="countdown"> {seconds_left}s </span>
            <span className="tiny-header">
              {game.layout.title}
              <br />
              {game.layout.author}, Author - {game.layout.editor}, Editor
            </span>
            <p className="game-id">Game ID: {game.id}</p>
          </div>
          <div className="board-layout">
            <Board game={game} />
          </div>
        </div>
      </div>
    );
  }
}
