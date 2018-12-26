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

    return (
      <div>
        <h2>TV Play</h2>
        <div className="tv-layout">
          <div className="tv-left">
            {game.current_players && (
              <p>
                <span className="score-header">Score</span>
                <ul className="scoreboard">
                  {game.current_players.map(p => (
                    <li key={p.name}>
                      {p.name}: {p.current_score} pts
                    </li>
                  ))}
                </ul>
              </p>
            )}
            {game.last_clue && (
              <p>
                <span className="last-clue-header">
                  Last Clue: <br />
                </span>
                <span className="last-clue">{game.last_clue.description}</span>
                <ul className="last-clue-score">
                  {(game.last_clue.guesses || []).map(g => (
                    <li key={g.player.name}>
                      {g.player.name}: {g.guess} ({g.score} pts)
                    </li>
                  ))}
                </ul>
              </p>
            )}
            {game.current_clue && <p className="current-clue">{game.current_clue.description}</p>}
            {/* <span className="countdown"> :30 </span> */}
            <p className="game-id">
              {/* <button is="google-cast-button" /> */}
              Game ID: {game.id}
            </p>
          </div>
          <div className="board-layout">
            <Board game={game} />
          </div>
        </div>
      </div>
    );
  }
}
